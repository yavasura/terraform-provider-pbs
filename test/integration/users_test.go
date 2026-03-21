package integration

import (
	"context"
	"fmt"
	"os"
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

func testUserID() string {
	realm := os.Getenv("PBS_TEST_USER_REALM")
	if realm == "" {
		realm = "pbs"
	}

	return fmt.Sprintf("%s@%s", GenerateTestName("tfuser"), realm)
}
