/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package notifications

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
)

func stringValueOrNull(value string) types.String {
	return tfvalue.StringOrNull(value)
}

func boolValueOrNull(value *bool) types.Bool {
	return tfvalue.BoolPtrOrNull(value)
}

func intValueOrNull(value *int) types.Int64 {
	return tfvalue.IntPtrOrNull(value)
}
