package access

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

func TestTokenNameRegex(t *testing.T) {
	valid := []string{
		"terraform",
		"ci-bot",
		"ops.token",
		"backup_2026",
	}

	for _, tokenName := range valid {
		assert.True(t, tokenNameRegex.MatchString(tokenName), "expected %q to be valid", tokenName)
	}

	invalid := []string{
		"",
		"has space",
		"bad/token",
		"bad!token",
	}

	for _, tokenName := range invalid {
		assert.False(t, tokenNameRegex.MatchString(tokenName), "expected %q to be invalid", tokenName)
	}
}

func TestAPITokenImportIDSplit(t *testing.T) {
	userID, tokenName, ok := pbsaccess.SplitAPITokenID("admin@pbs!terraform")
	assert.True(t, ok)
	assert.Equal(t, "admin@pbs", userID)
	assert.Equal(t, "terraform", tokenName)

	_, _, ok = pbsaccess.SplitAPITokenID("admin@pbs")
	assert.False(t, ok)
}

func TestBuildAPITokenFromModel(t *testing.T) {
	model := &apiTokenResourceModel{
		UserID:    types.StringValue("admin@pbs"),
		TokenName: types.StringValue("terraform"),
		Comment:   types.StringValue("Managed by Terraform"),
		Enable:    types.BoolValue(true),
		Expire:    types.Int64Value(1900000000),
	}

	token := buildAPITokenFromModel(model)

	assert.Equal(t, "admin@pbs", token.UserID)
	assert.Equal(t, "terraform", token.TokenName)
	assert.Equal(t, "admin@pbs!terraform", token.TokenID)
	assert.Equal(t, "Managed by Terraform", token.Comment)
	if assert.NotNil(t, token.Enable) {
		assert.True(t, *token.Enable)
	}
	if assert.NotNil(t, token.Expire) {
		assert.EqualValues(t, 1900000000, *token.Expire)
	}
}
