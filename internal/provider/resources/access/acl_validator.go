package access

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

type ugidValidator struct{}

func (v ugidValidator) Description(_ context.Context) string {
	return "must be a valid PBS user ID, token ID, or group ID"
}

func (v ugidValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ugidValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if userIDRegex.MatchString(value) {
		return
	}
	if _, _, ok := pbsaccess.SplitAPITokenID(value); ok {
		return
	}
	if groupIDRegex.MatchString(value) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid ACL Subject",
		fmt.Sprintf("ACL subject %q must be a PBS user ID, token ID, or group ID.", value),
	)
}
