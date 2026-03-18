package access

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestUserIDRegex(t *testing.T) {
	valid := []string{
		"admin@pam",
		"john@ldap",
		"svc-backup@ad",
		"ops.user@openid",
		"terraform@pbs",
	}

	for _, userID := range valid {
		assert.True(t, userIDRegex.MatchString(userID), "expected %q to be valid", userID)
	}

	invalid := []string{
		"admin",
		"admin@",
		"@pam",
		"admin @pam",
		"admin/pam",
	}

	for _, userID := range invalid {
		assert.False(t, userIDRegex.MatchString(userID), "expected %q to be invalid", userID)
	}
}

func TestComputeUserDeletes(t *testing.T) {
	state := &userResourceModel{
		Comment:   types.StringValue("existing"),
		Enable:    types.BoolValue(true),
		Expire:    types.Int64Value(1700000000),
		FirstName: types.StringValue("Jane"),
		LastName:  types.StringValue("Doe"),
		Email:     types.StringValue("jane@example.com"),
	}
	plan := &userResourceModel{
		Comment:   types.StringNull(),
		Enable:    types.BoolNull(),
		Expire:    types.Int64Null(),
		FirstName: types.StringNull(),
		LastName:  types.StringNull(),
		Email:     types.StringNull(),
	}

	assert.ElementsMatch(t, []string{
		"comment",
		"enable",
		"expire",
		"firstname",
		"lastname",
		"email",
	}, computeUserDeletes(plan, state))
}
