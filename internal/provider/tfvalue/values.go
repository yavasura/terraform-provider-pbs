/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tfvalue

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func BoolPtrOrNull(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

func BoolPtrOrDefault(value *bool, defaultValue bool) types.Bool {
	if value == nil {
		return types.BoolValue(defaultValue)
	}
	return types.BoolValue(*value)
}

func IntPtrOrNull(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

func Int64PtrOrNull(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*value)
}

func StringListOrNull(ctx context.Context, value []string) (types.List, diag.Diagnostics) {
	if len(value) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, value)
}

func StringMapOrNull(ctx context.Context, value map[string]string) (types.Map, diag.Diagnostics) {
	if len(value) == 0 {
		return types.MapNull(types.StringType), nil
	}
	return types.MapValueFrom(ctx, types.StringType, value)
}

func StringSetOrNull(ctx context.Context, value []string) (types.Set, diag.Diagnostics) {
	if len(value) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, value)
}
