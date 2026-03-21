package jobs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
)

func stringValueOrNull(value string) types.String {
	return tfvalue.StringOrNull(value)
}

func boolValueOrNull(value *bool) types.Bool {
	return tfvalue.BoolPtrOrNull(value)
}

func int64ValueOrNull(value *int) types.Int64 {
	return tfvalue.IntPtrOrNull(value)
}

func intPointerFromAttr(attr types.Int64) *int {
	if attr.IsNull() || attr.IsUnknown() {
		return nil
	}
	v := int(attr.ValueInt64())
	return &v
}

func boolPointerFromAttr(attr types.Bool) *bool {
	if attr.IsNull() || attr.IsUnknown() {
		return nil
	}
	v := attr.ValueBool()
	return &v
}

func stringPointerFromAttr(attr types.String) *string {
	if attr.IsNull() || attr.IsUnknown() {
		return nil
	}
	v := attr.ValueString()
	return &v
}

func shouldDeleteStringAttr(plan, state types.String) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func shouldDeleteIntAttr(plan, state types.Int64) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func shouldDeleteBoolAttr(plan, state types.Bool) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func shouldDeleteListAttr(plan, state types.List) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func stringListFromAttribute(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var items []string
	diags := list.ElementsAs(ctx, &items, false)
	return items, diags
}
