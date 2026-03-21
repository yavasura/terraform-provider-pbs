/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs/datastores"
)

func optionalInt64Pointer(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := int(value.ValueInt64())
	return &v
}

func optionalBoolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueBool()
	return &v
}

func hasBoolValue(value types.Bool) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func hasInt64Value(value types.Int64) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func isEmptyTuning(t *datastores.DatastoreTuning) bool {
	if t == nil {
		return true
	}
	return t.ChunkOrder == "" && t.GCAtimeCutoff == nil && t.GCAtimeSafetyCheck == nil && t.GCCacheCapacity == nil && t.SyncLevel == ""
}

func boolValueIsTrue(value types.Bool) bool {
	return !value.IsNull() && !value.IsUnknown() && value.ValueBool()
}

func tuneLevelToSyncLevel(level int) (string, error) {
	switch level {
	case 0:
		return "none", nil
	case 1:
		return "filesystem", nil
	case 2:
		return "file", nil
	default:
		return "", fmt.Errorf("unsupported tune_level %d; valid values are 0-2", level)
	}
}

func syncLevelToTuneLevel(syncLevel string) (int, bool) {
	switch strings.ToLower(syncLevel) {
	case "none":
		return 0, true
	case "filesystem":
		return 1, true
	case "file":
		return 2, true
	default:
		return 0, false
	}
}

func stringValueOrNull(value string) types.String {
	return tfvalue.StringOrNull(value)
}

func intValueOrNull(ptr *int) types.Int64 {
	return tfvalue.IntPtrOrNull(ptr)
}

func boolValueOrNull(ptr *bool) types.Bool {
	return tfvalue.BoolPtrOrNull(ptr)
}
