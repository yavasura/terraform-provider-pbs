/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

// Package metrics provides Terraform data sources for PBS metrics server configurations
package metrics

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfstate"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	"github.com/yavasura/terraform-provider-pbs/pbs/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &metricsServerDataSource{}
	_ datasource.DataSourceWithConfigure = &metricsServerDataSource{}
)

// NewMetricsServerDataSource is a helper function to simplify the provider implementation.
func NewMetricsServerDataSource() datasource.DataSource {
	return &metricsServerDataSource{}
}

// metricsServerDataSource is the data source implementation.
type metricsServerDataSource struct {
	client *pbs.Client
}

// metricsServerDataSourceModel maps the data source schema data.
type metricsServerDataSourceModel struct {
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	URL          types.String `tfsdk:"url"`
	Server       types.String `tfsdk:"server"`
	Port         types.Int64  `tfsdk:"port"`
	Enable       types.Bool   `tfsdk:"enable"`
	MTU          types.Int64  `tfsdk:"mtu"`
	Organization types.String `tfsdk:"organization"`
	Bucket       types.String `tfsdk:"bucket"`
	MaxBodySize  types.Int64  `tfsdk:"max_body_size"`
	VerifyTLS    types.Bool   `tfsdk:"verify_tls"`
	Comment      types.String `tfsdk:"comment"`
}

// Metadata returns the data source type name.
func (d *metricsServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metrics_server"
}

// Schema defines the schema for the data source.
func (d *metricsServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads a specific metrics server configuration from Proxmox Backup Server.",
		MarkdownDescription: "Reads a specific metrics server configuration from Proxmox Backup Server.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The unique name identifier for the metrics server.",
				MarkdownDescription: "The unique name identifier for the metrics server.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				Description:         "The type of metrics server.",
				MarkdownDescription: "The type of metrics server. Valid values: `influxdb-udp`, `influxdb-http`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("influxdb-udp", "influxdb-http"),
				},
			},
			"url": schema.StringAttribute{
				Description:         "Full URL for InfluxDB HTTP.",
				MarkdownDescription: "Full URL for InfluxDB HTTP (e.g., `http://influxdb:8086`). Only applicable for `influxdb-http` type.",
				Computed:            true,
			},
			"server": schema.StringAttribute{
				Description:         "The server address (hostname or IP).",
				MarkdownDescription: "The server address (hostname or IP) extracted from URL or Host field.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				Description:         "The server port.",
				MarkdownDescription: "The server port extracted from URL or Host field.",
				Computed:            true,
			},
			"enable": schema.BoolAttribute{
				Description:         "Whether this metrics server is enabled.",
				MarkdownDescription: "Whether metrics export to this server is enabled.",
				Computed:            true,
			},
			"mtu": schema.Int64Attribute{
				Description:         "MTU for the metrics connection.",
				MarkdownDescription: "Maximum transmission unit for the metrics connection.",
				Computed:            true,
			},
			"organization": schema.StringAttribute{
				Description:         "InfluxDB organization (InfluxDB HTTP only).",
				MarkdownDescription: "InfluxDB organization name. Only applicable for `influxdb-http` type.",
				Computed:            true,
			},
			"bucket": schema.StringAttribute{
				Description:         "InfluxDB bucket name (InfluxDB HTTP only).",
				MarkdownDescription: "InfluxDB bucket name. Only applicable for `influxdb-http` type.",
				Computed:            true,
			},
			"max_body_size": schema.Int64Attribute{
				Description:         "Maximum body size for HTTP requests in bytes (InfluxDB HTTP only).",
				MarkdownDescription: "Maximum body size for HTTP requests in bytes. Only applicable for `influxdb-http` type.",
				Computed:            true,
			},
			"verify_tls": schema.BoolAttribute{
				Description:         "Verify TLS certificate for HTTPS connections (InfluxDB HTTP only).",
				MarkdownDescription: "Whether to verify TLS certificate for HTTPS connections. Only applicable for `influxdb-http` type.",
				Computed:            true,
			},
			"comment": schema.StringAttribute{
				Description:         "A comment describing this metrics server.",
				MarkdownDescription: "A comment describing this metrics server configuration.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *metricsServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config.ConfigureDataSourceClient(&d.client, req.ProviderData, &resp.Diagnostics)
}

// Read refreshes the Terraform state with the latest data.
func (d *metricsServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state metricsServerDataSourceModel
	if !tfstate.Decode(ctx, req.Config, &state, &resp.Diagnostics) {
		return
	}

	// Get metrics server from API
	serverType := metrics.MetricsServerType(state.Type.ValueString())
	server, err := d.client.Metrics.GetMetricsServer(ctx, serverType, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Metrics Server",
			fmt.Sprintf("Could not read metrics server %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	// Map API response to state
	if err := metricsServerToState(server, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Metrics Server",
			fmt.Sprintf("Could not convert metrics server to state: %s", err.Error()),
		)
		return
	}

	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

// metricsServerToState converts a metrics server struct to Terraform state
func metricsServerToState(server *metrics.MetricsServer, state *metricsServerDataSourceModel) error {
	state.Name = types.StringValue(server.Name)
	state.Type = types.StringValue(string(server.Type))
	state.Comment = tfvalue.StringOrNull(server.Comment)
	state.Enable = tfvalue.BoolPtrOrDefault(server.Enable, true)

	// URL and parsed server/port fields
	state.URL = tfvalue.StringOrNull(server.URL)
	state.Server = tfvalue.StringOrNull(server.Server)
	if server.Port > 0 {
		state.Port = types.Int64Value(int64(server.Port))
	} else {
		state.Port = types.Int64Null()
	}

	// Type-specific fields
	if server.Type == metrics.MetricsServerTypeInfluxDBUDP {
		state.MTU = tfvalue.IntPtrOrNull(server.MTU)
	}

	if server.Type == metrics.MetricsServerTypeInfluxDBHTTP {
		state.Organization = tfvalue.StringOrNull(server.Organization)
		state.Bucket = tfvalue.StringOrNull(server.Bucket)
		state.MaxBodySize = tfvalue.IntPtrOrNull(server.MaxBodySize)
		state.VerifyTLS = tfvalue.BoolPtrOrNull(server.VerifyTLS)
	}

	return nil
}
