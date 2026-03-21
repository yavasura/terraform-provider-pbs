package namespaces

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pbsdatastores "github.com/yavasura/terraform-provider-pbs/pbs/datastores"
)

type namespacePathValidator struct{}

func (v namespacePathValidator) Description(_ context.Context) string {
	return "must be a valid PBS namespace path with at most 7 segments"
}

func (v namespacePathValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v namespacePathValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if !pbsdatastores.IsValidNamespacePath(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Namespace Path",
			fmt.Sprintf("Namespace %q must use PBS-safe path segments and contain at most 7 levels.", value),
		)
	}
}

func boolDefaultFalse() types.Bool {
	return types.BoolValue(false)
}
