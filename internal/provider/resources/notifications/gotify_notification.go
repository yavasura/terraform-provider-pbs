/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package notifications

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfschema"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	"github.com/yavasura/terraform-provider-pbs/pbs/notifications"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &gotifyNotificationResource{}
	_ resource.ResourceWithConfigure   = &gotifyNotificationResource{}
	_ resource.ResourceWithImportState = &gotifyNotificationResource{}
)

// NewGotifyNotificationResource is a helper function to simplify the provider implementation.
func NewGotifyNotificationResource() resource.Resource {
	return &gotifyNotificationResource{}
}

// gotifyNotificationResource is the resource implementation.
type gotifyNotificationResource struct {
	client *pbs.Client
}

// gotifyNotificationResourceModel maps the resource schema data.
type gotifyNotificationResourceModel struct {
	Name    types.String `tfsdk:"name"`
	Server  types.String `tfsdk:"server"`
	Token   types.String `tfsdk:"token"`
	Comment types.String `tfsdk:"comment"`
	Disable types.Bool   `tfsdk:"disable"`
	Origin  types.String `tfsdk:"origin"`
}

// Metadata returns the resource type name.
func (r *gotifyNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gotify_notification"
}

// Schema defines the schema for the resource.
func (r *gotifyNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Gotify notification target for PBS alerts and notifications.",
		MarkdownDescription: `Manages a Gotify notification target.

Configure a Gotify server to receive push notifications from PBS about backup jobs,
verification tasks, and system events.`,
		Attributes: map[string]schema.Attribute{
			"name": tfschema.RequiredReplaceStringAttribute(
				"The unique name identifier for the Gotify target.",
				"The unique name identifier for the Gotify notification target.",
			),
			"server": schema.StringAttribute{
				Description:         "Gotify server URL.",
				MarkdownDescription: "Gotify server URL (e.g., `https://gotify.example.com`).",
				Required:            true,
			},
			"token": schema.StringAttribute{
				Description:         "Gotify application token.",
				MarkdownDescription: "Gotify application token for authentication.",
				Required:            true,
				Sensitive:           true,
			},
			"comment": tfschema.OptionalCommentAttribute(
				"A comment describing this notification target.",
				"A comment describing this notification target.",
			),
			"disable": tfschema.OptionalComputedBoolDefaultFalseAttribute(
				"Disable this notification target.",
				"Disable this notification target. Defaults to `false`.",
			),
			"origin": schema.StringAttribute{
				Description:         "Origin of this configuration as reported by PBS.",
				MarkdownDescription: "Origin of this configuration as reported by PBS (e.g., `user`, `builtin`).",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *gotifyNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *gotifyNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan gotifyNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	target := &notifications.GotifyTarget{
		Name:   plan.Name.ValueString(),
		Server: plan.Server.ValueString(),
		Token:  plan.Token.ValueString(),
	}

	if !plan.Comment.IsNull() {
		target.Comment = plan.Comment.ValueString()
	}
	if !plan.Disable.IsNull() {
		disable := plan.Disable.ValueBool()
		target.Disable = &disable
	}

	err := r.client.Notifications.CreateGotifyTarget(ctx, target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Gotify notification target",
			fmt.Sprintf("Could not create Gotify notification target %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	created, err := r.client.Notifications.GetGotifyTarget(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created Gotify notification target",
			fmt.Sprintf("Created Gotify notification target %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	applyGotifyNotificationState(created, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *gotifyNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gotifyNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	target, err := r.client.Notifications.GetGotifyTarget(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Gotify notification target",
			fmt.Sprintf("Could not read Gotify notification target %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	applyGotifyNotificationState(target, &state)
	// Don't update token from API (sensitive field)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *gotifyNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan gotifyNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	target := &notifications.GotifyTarget{
		Name:   plan.Name.ValueString(),
		Server: plan.Server.ValueString(),
		Token:  plan.Token.ValueString(),
	}

	if !plan.Comment.IsNull() {
		target.Comment = plan.Comment.ValueString()
	}
	if !plan.Disable.IsNull() {
		disable := plan.Disable.ValueBool()
		target.Disable = &disable
	}

	err := r.client.Notifications.UpdateGotifyTarget(ctx, plan.Name.ValueString(), target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Gotify notification target",
			fmt.Sprintf("Could not update Gotify notification target %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	updated, err := r.client.Notifications.GetGotifyTarget(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated Gotify notification target",
			fmt.Sprintf("Updated Gotify notification target %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	applyGotifyNotificationState(updated, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *gotifyNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state gotifyNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Notifications.DeleteGotifyTarget(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Gotify notification target",
			fmt.Sprintf("Could not delete Gotify notification target %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *gotifyNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func applyGotifyNotificationState(target *notifications.GotifyTarget, state *gotifyNotificationResourceModel) {
	state.Server = types.StringValue(target.Server)
	setNotificationCommonState(target.Comment, target.Disable, target.Origin, &state.Comment, &state.Disable, &state.Origin)
}
