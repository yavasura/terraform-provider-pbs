package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

func TestUsersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	userID := testUserID()
	usersClient := pbsaccess.NewClient(tc.APIClient)

	defer func() {
		_ = usersClient.DeleteUser(context.Background(), userID, "")
	}()

	config := fmt.Sprintf(`
resource "pbs_user" "test_user" {
  userid    = "%s"
  comment   = "Managed by Terraform"
  enable    = true
  firstname = "Terraform"
  lastname  = "User"
  email     = "terraform.user@example.com"
}
`, userID)

	tc.WriteMainTF(t, config)
	tc.ApplyTerraform(t)

	resource := tc.GetResourceFromState(t, "pbs_user.test_user")
	assert.Equal(t, userID, resource.AttributeValues["userid"])
	assert.Equal(t, "Managed by Terraform", resource.AttributeValues["comment"])
	assert.Equal(t, "Terraform", resource.AttributeValues["firstname"])
	assert.Equal(t, "User", resource.AttributeValues["lastname"])
	assert.Equal(t, "terraform.user@example.com", resource.AttributeValues["email"])

	user, err := usersClient.GetUser(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, "Managed by Terraform", user.Comment)
	if assert.NotNil(t, user.Enable) {
		assert.True(t, *user.Enable)
	}

	updatedConfig := fmt.Sprintf(`
resource "pbs_user" "test_user" {
  userid    = "%s"
  comment   = "Updated by Terraform"
  enable    = false
  expire    = 1893456000
  firstname = "Updated"
  lastname  = "Person"
  email     = "updated.user@example.com"
}
`, userID)

	tc.WriteMainTF(t, updatedConfig)
	tc.ApplyTerraform(t)

	user, err = usersClient.GetUser(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, "Updated by Terraform", user.Comment)
	if assert.NotNil(t, user.Enable) {
		assert.False(t, *user.Enable)
	}
	if assert.NotNil(t, user.Expire) {
		assert.EqualValues(t, 1893456000, *user.Expire)
	}
	assert.Equal(t, "Updated", user.FirstName)
	assert.Equal(t, "Person", user.LastName)
	assert.Equal(t, "updated.user@example.com", user.Email)

	clearedConfig := fmt.Sprintf(`
resource "pbs_user" "test_user" {
  userid = "%s"
  enable = false
}
`, userID)

	tc.WriteMainTF(t, clearedConfig)
	tc.ApplyTerraform(t)

	user, err = usersClient.GetUser(context.Background(), userID)
	require.NoError(t, err)
	if assert.NotNil(t, user.Enable) {
		assert.False(t, *user.Enable)
	}
	assert.Empty(t, user.Comment)
	assert.Nil(t, user.Expire)
	assert.Empty(t, user.FirstName)
	assert.Empty(t, user.LastName)
	assert.Empty(t, user.Email)
}

func TestUserImportIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	userID := testUserID()
	usersClient := pbsaccess.NewClient(tc.APIClient)

	enable := false
	importUser := &pbsaccess.User{
		UserID:    userID,
		Comment:   "Imported by Terraform",
		Enable:    &enable,
		FirstName: "Import",
		LastName:  "User",
		Email:     "import.user@example.com",
	}

	err := usersClient.CreateUser(context.Background(), importUser)
	require.NoError(t, err)

	defer func() {
		_ = usersClient.DeleteUser(context.Background(), userID, "")
	}()

	importConfig := fmt.Sprintf(`
resource "pbs_user" "imported" {
  userid    = "%s"
  comment   = "Imported by Terraform"
  enable    = false
  firstname = "Import"
  lastname  = "User"
  email     = "import.user@example.com"
}
`, userID)

	tc.WriteMainTF(t, importConfig)
	tc.ImportResource(t, "pbs_user.imported", userID)

	resource := tc.GetResourceFromState(t, "pbs_user.imported")
	assert.Equal(t, userID, resource.AttributeValues["userid"])
	assert.Equal(t, "Imported by Terraform", resource.AttributeValues["comment"])

	tc.ApplyTerraform(t)
}

func TestUsersDockerSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	usersClient := pbsaccess.NewClient(tc.APIClient)
	userID := testUserID()
	destroyedByTerraform := false

	t.Cleanup(func() {
		if destroyedByTerraform {
			return
		}

		err := usersClient.DeleteUser(context.Background(), userID, "")
		if err != nil && !isPBSUserMissingError(err, userID) {
			t.Logf("cleanup warning: failed to delete %s directly: %v", userID, err)
		}
	})

	initialConfig := fmt.Sprintf(`
resource "pbs_user" "docker_smoke" {
  userid  = "%s"
  comment = "Docker smoke test user"
  enable  = true
}
`, userID)

	t.Logf("terraform apply: create resource pbs_user.docker_smoke for %s", userID)
	tc.WriteMainTF(t, initialConfig)
	tc.ApplyTerraform(t)

	resource := tc.GetResourceFromState(t, "pbs_user.docker_smoke")
	t.Logf("verified terraform state: pbs_user.docker_smoke.userid=%v comment=%v", resource.AttributeValues["userid"], resource.AttributeValues["comment"])
	assert.Equal(t, userID, resource.AttributeValues["userid"])
	assert.Equal(t, "Docker smoke test user", resource.AttributeValues["comment"])

	user, err := usersClient.GetUser(context.Background(), userID)
	require.NoError(t, err)
	t.Logf("verified PBS API after create: user %s exists with comment=%q enabled=%v", user.UserID, user.Comment, user.Enable != nil && *user.Enable)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, "Docker smoke test user", user.Comment)
	if assert.NotNil(t, user.Enable) {
		assert.True(t, *user.Enable)
	}

	updatedConfig := fmt.Sprintf(`
resource "pbs_user" "docker_smoke" {
  userid  = "%s"
  comment = "Docker smoke test user updated"
  enable  = false
}
`, userID)

	t.Logf("terraform apply: update resource pbs_user.docker_smoke for %s", userID)
	tc.WriteMainTF(t, updatedConfig)
	tc.ApplyTerraform(t)

	user, err = usersClient.GetUser(context.Background(), userID)
	require.NoError(t, err)
	t.Logf("verified PBS API after update: user %s has comment=%q enabled=%v", user.UserID, user.Comment, user.Enable != nil && *user.Enable)
	assert.Equal(t, "Docker smoke test user updated", user.Comment)
	if assert.NotNil(t, user.Enable) {
		assert.False(t, *user.Enable)
	}

	t.Logf("terraform destroy: remove resource pbs_user.docker_smoke for %s", userID)
	tc.DestroyTerraform(t)
	destroyedByTerraform = true

	_, err = usersClient.GetUser(context.Background(), userID)
	require.Error(t, err)
	t.Logf("verified PBS API after destroy: user %s no longer exists", userID)
}

func testUserID() string {
	realm := os.Getenv("PBS_TEST_USER_REALM")
	if realm == "" {
		realm = "pbs"
	}

	return fmt.Sprintf("%s@%s", GenerateTestName("tfuser"), realm)
}

func isPBSUserMissingError(err error, userID string) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	if !strings.Contains(msg, userID) {
		return false
	}

	return strings.Contains(msg, "no such user") || strings.Contains(msg, "does not exist")
}
