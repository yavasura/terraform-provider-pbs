/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tfstate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type getter interface {
	Get(context.Context, any) diag.Diagnostics
}

type setter interface {
	Set(context.Context, any) diag.Diagnostics
}

func Decode[T any](ctx context.Context, src getter, dst *T, diags *diag.Diagnostics) bool {
	diags.Append(src.Get(ctx, dst)...)
	return !diags.HasError()
}

func Encode[T any](ctx context.Context, dst setter, src *T, diags *diag.Diagnostics) bool {
	diags.Append(dst.Set(ctx, src)...)
	return !diags.HasError()
}
