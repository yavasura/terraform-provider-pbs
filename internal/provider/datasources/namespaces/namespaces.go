/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package namespaces

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/pbs"
)

var (
	_ datasource.DataSource              = &namespacesDataSource{}
	_ datasource.DataSourceWithConfigure = &namespacesDataSource{}
)

// NewNamespacesDataSource returns a datastore namespace listing data source.
func NewNamespacesDataSource() datasource.DataSource {
	return &namespacesDataSource{}
}

type namespacesDataSource struct {
	client *pbs.Client
}

type namespacesDataSourceModel struct {
	Store      types.String             `tfsdk:"store"`
	Prefix     types.String             `tfsdk:"prefix"`
	MaxDepth   types.Int64              `tfsdk:"max_depth"`
	Namespaces []namespaceListItemModel `tfsdk:"namespaces"`
}

type namespaceListItemModel struct {
	Namespace types.String `tfsdk:"namespace"`
	Comment   types.String `tfsdk:"comment"`
}

func (d *namespacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespaces"
}

func (d *namespacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Lists namespaces in a PBS datastore.",
		MarkdownDescription: "Lists namespaces in a PBS datastore, with optional prefix and max-depth filtering.",
		Attributes: map[string]schema.Attribute{
			"store": schema.StringAttribute{
				Description:         "Datastore to inspect.",
				MarkdownDescription: "Datastore to inspect.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"prefix": schema.StringAttribute{
				Description:         "Optional namespace prefix filter.",
				MarkdownDescription: "Optional namespace prefix filter. Only namespace paths starting with this prefix are returned.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"max_depth": schema.Int64Attribute{
				Description:         "Optional maximum namespace depth to request from PBS.",
				MarkdownDescription: "Optional maximum namespace depth to request from PBS. Must be between 0 and 7.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 7),
				},
			},
			"namespaces": schema.ListNestedAttribute{
				Description:         "Matching namespaces.",
				MarkdownDescription: "Matching namespaces.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"namespace": schema.StringAttribute{
							Description:         "Namespace path.",
							MarkdownDescription: "Namespace path.",
							Computed:            true,
						},
						"comment": schema.StringAttribute{
							Description:         "Namespace comment.",
							MarkdownDescription: "Namespace comment.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *namespacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config.ConfigureDataSourceClient(&d.client, req.ProviderData, &resp.Diagnostics)
}

func (d *namespacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state namespacesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var maxDepth *int
	if !state.MaxDepth.IsNull() && !state.MaxDepth.IsUnknown() {
		value := int(state.MaxDepth.ValueInt64())
		maxDepth = &value
	}

	namespaces, err := d.client.Datastores.ListNamespaces(ctx, state.Store.ValueString(), maxDepth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing namespaces",
			fmt.Sprintf("Could not list namespaces for datastore %s: %s", state.Store.ValueString(), err.Error()),
		)
		return
	}

	prefix := ""
	if !state.Prefix.IsNull() && !state.Prefix.IsUnknown() {
		prefix = state.Prefix.ValueString()
	}

	items := make([]namespaceListItemModel, 0, len(namespaces))
	for _, ns := range namespaces {
		if prefix != "" && !strings.HasPrefix(ns.Path, prefix) {
			continue
		}

		items = append(items, namespaceListItemModel{
			Namespace: types.StringValue(ns.Path),
			Comment:   types.StringValue(ns.Comment),
		})
	}

	state.Namespaces = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
