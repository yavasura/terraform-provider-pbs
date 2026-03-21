package access

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestACLPathRegex(t *testing.T) {
	valid := []string{
		"/",
		"/datastore",
		"/datastore/backups",
		"/remote/secondary",
		"/tape",
		"/system",
		"/access",
	}

	for _, aclPath := range valid {
		assert.True(t, aclPathRegex.MatchString(aclPath), "expected %q to be valid", aclPath)
	}

	invalid := []string{
		"",
		"datastore/backups",
		"/datastore/",
		"/bad path",
	}

	for _, aclPath := range invalid {
		assert.False(t, aclPathRegex.MatchString(aclPath), "expected %q to be invalid", aclPath)
	}
}

func TestBuildACLFromModel(t *testing.T) {
	model := &aclResourceModel{
		Path:      types.StringValue("/datastore/backups"),
		UGID:      types.StringValue("admin@pbs"),
		RoleID:    types.StringValue("DatastoreAdmin"),
		Propagate: types.BoolValue(true),
	}

	acl := buildACLFromModel(model)
	assert.Equal(t, "/datastore/backups", acl.Path)
	assert.Equal(t, "admin@pbs", acl.UGID)
	assert.Equal(t, "DatastoreAdmin", acl.RoleID)
	if assert.NotNil(t, acl.Propagate) {
		assert.True(t, *acl.Propagate)
	}
}
