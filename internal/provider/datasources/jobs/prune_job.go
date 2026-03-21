/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package jobs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfstate"
	"github.com/yavasura/terraform-provider-pbs/pbs"
)

var (
	_ datasource.DataSource              = &pruneJobDataSource{}
	_ datasource.DataSourceWithConfigure = &pruneJobDataSource{}
)

// NewPruneJobDataSource is a helper function to simplify the provider implementation.
func NewPruneJobDataSource() datasource.DataSource {
	return &pruneJobDataSource{}
}

// pruneJobDataSource is the data source implementation.
type pruneJobDataSource struct {
	client *pbs.Client
}

// pruneJobDataSourceModel maps the data source schema data.
type pruneJobDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Store       types.String `tfsdk:"store"`
	Schedule    types.String `tfsdk:"schedule"`
	KeepLast    types.Int64  `tfsdk:"keep_last"`
	KeepHourly  types.Int64  `tfsdk:"keep_hourly"`
	KeepDaily   types.Int64  `tfsdk:"keep_daily"`
	KeepWeekly  types.Int64  `tfsdk:"keep_weekly"`
	KeepMonthly types.Int64  `tfsdk:"keep_monthly"`
	KeepYearly  types.Int64  `tfsdk:"keep_yearly"`
	MaxDepth    types.Int64  `tfsdk:"max_depth"`
	Namespace   types.String `tfsdk:"namespace"`
	Comment     types.String `tfsdk:"comment"`
	Disable     types.Bool   `tfsdk:"disable"`
	Digest      types.String `tfsdk:"digest"`
}

// Metadata returns the data source type name.
func (d *pruneJobDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prune_job"
}

// Schema defines the schema for the data source.
func (d *pruneJobDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads information about an existing PBS prune job.",
		MarkdownDescription: "Reads information about an existing PBS prune job configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier for the prune job.",
				MarkdownDescription: "The unique identifier for the prune job.",
				Required:            true,
			},
			"store": schema.StringAttribute{
				Description:         "The datastore name where pruning is performed.",
				MarkdownDescription: "The datastore name where pruning is performed.",
				Computed:            true,
			},
			"schedule": schema.StringAttribute{
				Description:         "When the prune job runs (systemd calendar event format).",
				MarkdownDescription: "When the prune job runs (systemd calendar event format).",
				Computed:            true,
			},
			"keep_last": schema.Int64Attribute{
				Description:         "Keep the last N backup snapshots.",
				MarkdownDescription: "Keep the last N backup snapshots.",
				Computed:            true,
			},
			"keep_hourly": schema.Int64Attribute{
				Description:         "Keep hourly backups for the last N hours.",
				MarkdownDescription: "Keep hourly backups for the last N hours.",
				Computed:            true,
			},
			"keep_daily": schema.Int64Attribute{
				Description:         "Keep daily backups for the last N days.",
				MarkdownDescription: "Keep daily backups for the last N days.",
				Computed:            true,
			},
			"keep_weekly": schema.Int64Attribute{
				Description:         "Keep weekly backups for the last N weeks.",
				MarkdownDescription: "Keep weekly backups for the last N weeks.",
				Computed:            true,
			},
			"keep_monthly": schema.Int64Attribute{
				Description:         "Keep monthly backups for the last N months.",
				MarkdownDescription: "Keep monthly backups for the last N months.",
				Computed:            true,
			},
			"keep_yearly": schema.Int64Attribute{
				Description:         "Keep yearly backups for the last N years.",
				MarkdownDescription: "Keep yearly backups for the last N years.",
				Computed:            true,
			},
			"max_depth": schema.Int64Attribute{
				Description:         "Maximum depth for namespace traversal.",
				MarkdownDescription: "Maximum depth for namespace traversal.",
				Computed:            true,
			},
			"namespace": schema.StringAttribute{
				Description:         "Namespace filter (regex).",
				MarkdownDescription: "Namespace filter as a regular expression.",
				Computed:            true,
			},
			"comment": schema.StringAttribute{
				Description:         "A comment describing this prune job.",
				MarkdownDescription: "A comment describing this prune job.",
				Computed:            true,
			},
			"disable": schema.BoolAttribute{
				Description:         "Whether this prune job is disabled.",
				MarkdownDescription: "Whether this prune job is disabled.",
				Computed:            true,
			},
			"digest": schema.StringAttribute{
				Description:         "Opaque digest returned by PBS.",
				MarkdownDescription: "Opaque digest returned by PBS.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *pruneJobDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config.ConfigureDataSourceClient(&d.client, req.ProviderData, &resp.Diagnostics)
}

// Read refreshes the Terraform state with the latest data.
func (d *pruneJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pruneJobDataSourceModel
	if !tfstate.Decode(ctx, req.Config, &state, &resp.Diagnostics) {
		return
	}

	// Get prune job from API
	job, err := d.client.Jobs.GetPruneJob(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Prune Job",
			fmt.Sprintf("Could not read prune job %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	setPruneJobDataSourceState(job, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}
