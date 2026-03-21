package access

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var aclRootSegments = []string{
	"access",
	"datastore",
	"remote",
	"system",
	"tape",
}

type aclPathValidator struct{}

func (v aclPathValidator) Description(_ context.Context) string {
	return "must be a valid PBS ACL path such as /, /datastore/<name>, /remote/<name>, /tape, /access, or /system"
}

func (v aclPathValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v aclPathValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "/" {
		return
	}

	parts := strings.Split(strings.TrimPrefix(value, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid ACL Path",
			fmt.Sprintf("ACL path %q must be a valid PBS access path.", value),
		)
		return
	}

	if !slices.Contains(aclRootSegments, parts[0]) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid ACL Path",
			fmt.Sprintf("ACL path %q must start with one of /access, /datastore, /remote, /system, or /tape, or be /.", value),
		)
	}
}
