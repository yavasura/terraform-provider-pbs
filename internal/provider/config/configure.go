/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/yavasura/terraform-provider-pbs/pbs"
)

// ConfigureResourceClient loads the configured provider client for a resource.
func ConfigureResourceClient(target **pbs.Client, providerData any, diags *diag.Diagnostics) {
	if providerData == nil {
		return
	}

	cfg, ok := providerData.(*Resource)
	if !ok {
		diags.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.Resource, got: %T. Please report this issue to the provider developers.", providerData),
		)
		return
	}

	*target = cfg.Client
}

// ConfigureDataSourceClient loads the configured provider client for a data source.
func ConfigureDataSourceClient(target **pbs.Client, providerData any, diags *diag.Diagnostics) {
	if providerData == nil {
		return
	}

	cfg, ok := providerData.(*DataSource)
	if !ok {
		diags.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *config.DataSource, got: %T. Please report this issue to the provider developers.", providerData),
		)
		return
	}

	*target = cfg.Client
}
