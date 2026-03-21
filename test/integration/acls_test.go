package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

func TestACLUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	userID := testUserID()
	aclClient := pbsaccess.NewClient(tc.APIClient)
	usersClient := pbsaccess.NewClient(tc.APIClient)

	defer func() {
		_ = aclClient.DeleteACL(context.Background(), "/", userID, "Audit")
		_ = usersClient.DeleteUser(context.Background(), userID, "")
	}()

	config := fmt.Sprintf(`
resource "pbs_user" "acl_user" {
  userid  = "%s"
  comment = "ACL integration user"
  enable  = true
}

resource "pbs_acl" "acl_user_root" {
  path      = "/"
  ugid      = pbs_user.acl_user.userid
  role_id   = "Audit"
  propagate = true
}
`, userID)

	tc.WriteMainTF(t, config)
	tc.ApplyTerraform(t)

	acl, err := aclClient.GetACL(context.Background(), "/", userID)
	require.NoError(t, err)
	assert.Equal(t, "Audit", acl.RoleID)
	if assert.NotNil(t, acl.Propagate) {
		assert.True(t, *acl.Propagate)
	}
}

func TestACLTokenIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	userID := testUserID()
	tokenID := userID + "!terraform"
	aclClient := pbsaccess.NewClient(tc.APIClient)
	usersClient := pbsaccess.NewClient(tc.APIClient)

	defer func() {
		_ = aclClient.DeleteACL(context.Background(), "/remote", tokenID, "RemoteAdmin")
		_ = usersClient.DeleteUserToken(context.Background(), userID, "terraform", "")
		_ = usersClient.DeleteUser(context.Background(), userID, "")
	}()

	config := fmt.Sprintf(`
resource "pbs_user" "acl_token_user" {
  userid  = "%s"
  comment = "ACL token integration user"
  enable  = true
}

resource "pbs_api_token" "terraform" {
  userid     = pbs_user.acl_token_user.userid
  token_name = "terraform"
}

resource "pbs_acl" "token_remote_acl" {
  path      = "/remote"
  ugid      = pbs_api_token.terraform.tokenid
  role_id   = "RemoteAdmin"
  propagate = false
}
`, userID)

	tc.WriteMainTF(t, config)
	tc.ApplyTerraform(t)

	acl, err := aclClient.GetACL(context.Background(), "/remote", tokenID)
	require.NoError(t, err)
	assert.Equal(t, "RemoteAdmin", acl.RoleID)
	if assert.NotNil(t, acl.Propagate) {
		assert.False(t, *acl.Propagate)
	}
}
