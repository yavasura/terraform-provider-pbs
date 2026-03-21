/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tfplan

import "github.com/hashicorp/terraform-plugin-framework/types"

// DeleteSet tracks unique PBS delete keys.
type DeleteSet struct {
	keys []string
}

func (d *DeleteSet) Add(key string) {
	for _, existing := range d.keys {
		if existing == key {
			return
		}
	}
	d.keys = append(d.keys, key)
}

func (d *DeleteSet) AddIf(key string, cond bool) {
	if cond {
		d.Add(key)
	}
}

func (d *DeleteSet) Values() []string {
	return d.keys
}

func ShouldDeleteString(plan, state types.String) bool {
	if state.IsNull() || state.IsUnknown() {
		return false
	}
	return plan.IsNull() || plan.IsUnknown()
}

func ShouldDeleteBool(plan, state types.Bool) bool {
	if state.IsNull() || state.IsUnknown() {
		return false
	}
	return plan.IsNull() || plan.IsUnknown()
}

func ShouldDeleteInt64(plan, state types.Int64) bool {
	if state.IsNull() || state.IsUnknown() {
		return false
	}
	return plan.IsNull() || plan.IsUnknown()
}
