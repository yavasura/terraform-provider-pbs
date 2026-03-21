/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package notifications

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
)

func shouldDeleteListAttr(plan, state types.List) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func shouldDeleteStringAttr(plan, state types.String) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func setNotificationCommonState(comment string, disable *bool, origin string, commentState *types.String, disableState *types.Bool, originState *types.String) {
	*commentState = tfvalue.StringOrNull(comment)
	*disableState = tfvalue.BoolPtrOrNull(disable)
	*originState = tfvalue.StringOrNull(origin)
}

func stringListState(ctx context.Context, value []string) (types.List, diag.Diagnostics) {
	if len(value) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, value)
}

func stringMapState(ctx context.Context, value map[string]string) (types.Map, diag.Diagnostics) {
	if len(value) == 0 {
		return types.MapNull(types.StringType), nil
	}
	return types.MapValueFrom(ctx, types.StringType, value)
}

func decodeStringList(ctx context.Context, value types.List, diags *diag.Diagnostics) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var items []string
	diags.Append(value.ElementsAs(ctx, &items, false)...)
	if diags.HasError() {
		return nil
	}
	return items
}
