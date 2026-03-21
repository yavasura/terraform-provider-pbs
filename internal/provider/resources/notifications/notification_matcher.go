/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package notifications

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	"github.com/yavasura/terraform-provider-pbs/pbs/notifications"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &notificationMatcherResource{}
	_ resource.ResourceWithConfigure   = &notificationMatcherResource{}
	_ resource.ResourceWithImportState = &notificationMatcherResource{}
)

// NewNotificationMatcherResource is a helper function to simplify the provider implementation.
func NewNotificationMatcherResource() resource.Resource {
	return &notificationMatcherResource{}
}

// notificationMatcherResource is the resource implementation.
type notificationMatcherResource struct {
	client *pbs.Client
}

// notificationMatcherResourceModel maps the resource schema data.
type notificationMatcherResourceModel struct {
	Name          types.String `tfsdk:"name"`
	Targets       types.List   `tfsdk:"targets"`
	MatchSeverity types.List   `tfsdk:"match_severity"`
	MatchField    types.List   `tfsdk:"match_field"`
	MatchCalendar types.List   `tfsdk:"match_calendar"`
	Mode          types.String `tfsdk:"mode"`
	InvertMatch   types.Bool   `tfsdk:"invert_match"`
	Comment       types.String `tfsdk:"comment"`
	Disable       types.Bool   `tfsdk:"disable"`
	Origin        types.String `tfsdk:"origin"`
}

// Metadata returns the resource type name.
func (r *notificationMatcherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_matcher"
}

// Schema defines the schema for the resource.
func (r *notificationMatcherResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a notification matcher (routing rule).",
		MarkdownDescription: `Manages a notification matcher.

Notification matchers define rules for routing notification events to specific targets or endpoints. 
They can filter notifications based on severity, custom fields, and calendar schedules. This allows 
for sophisticated notification routing, such as sending critical errors to on-call staff or filtering 
informational messages to specific channels.`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The unique name identifier for the notification matcher.",
				MarkdownDescription: "The unique name identifier for the notification matcher.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"targets": schema.ListAttribute{
				Description:         "List of notification target or endpoint names to route matching notifications to.",
				MarkdownDescription: "List of notification target or endpoint names to route matching notifications to.",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"match_severity": schema.ListAttribute{
				Description:         "List of severity levels to match (info, notice, warning, error).",
				MarkdownDescription: "List of severity levels to match. Valid values: `info`, `notice`, `warning`, `error`.",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"match_field": schema.ListAttribute{
				Description:         "List of field=value pairs to match against notification metadata.",
				MarkdownDescription: "List of `field=value` pairs to match against notification metadata (e.g., `type=prune`, `hostname=server01`).",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"match_calendar": schema.ListAttribute{
				Description:         "List of calendar IDs to match for time-based routing.",
				MarkdownDescription: "List of calendar IDs to match for time-based routing (requires calendar configuration in PBS).",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"mode": schema.StringAttribute{
				Description:         "Match mode: all (all conditions must match) or any (at least one condition must match).",
				MarkdownDescription: "Match mode: `all` (all conditions must match) or `any` (at least one condition must match). Defaults to `all`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("all"),
				Validators: []validator.String{
					stringvalidator.OneOf("all", "any"),
				},
			},
			"invert_match": schema.BoolAttribute{
				Description:         "Invert the matching logic (route notifications that DON'T match the criteria).",
				MarkdownDescription: "Invert the matching logic (route notifications that DON'T match the criteria). Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"comment": schema.StringAttribute{
				Description:         "A comment describing this notification matcher.",
				MarkdownDescription: "A comment describing this notification matcher.",
				Optional:            true,
			},
			"disable": schema.BoolAttribute{
				Description:         "Disable this notification matcher.",
				MarkdownDescription: "Disable this notification matcher. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"origin": schema.StringAttribute{
				Description:         "Origin of this configuration as reported by PBS.",
				MarkdownDescription: "Origin of this configuration as reported by PBS (e.g., `user`, `builtin`).",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *notificationMatcherResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *notificationMatcherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationMatcherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	matcher := &notifications.NotificationMatcher{
		Name: plan.Name.ValueString(),
	}
	matcher = buildNotificationMatcher(ctx, &plan, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Notifications.CreateNotificationMatcher(ctx, matcher)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification matcher",
			fmt.Sprintf("Could not create notification matcher %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	created, err := r.client.Notifications.GetNotificationMatcher(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created notification matcher",
			fmt.Sprintf("Created notification matcher %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	setNotificationMatcherState(ctx, created, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *notificationMatcherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationMatcherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	matcher, err := r.client.Notifications.GetNotificationMatcher(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading notification matcher",
			fmt.Sprintf("Could not read notification matcher %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	setNotificationMatcherState(ctx, matcher, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationMatcherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationMatcherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state notificationMatcherResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	matcher := buildNotificationMatcher(ctx, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Notifications.UpdateNotificationMatcher(ctx, plan.Name.ValueString(), matcher)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating notification matcher",
			fmt.Sprintf("Could not update notification matcher %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	updated, err := r.client.Notifications.GetNotificationMatcher(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated notification matcher",
			fmt.Sprintf("Updated notification matcher %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	setNotificationMatcherState(ctx, updated, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *notificationMatcherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationMatcherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Notifications.DeleteNotificationMatcher(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting notification matcher",
			fmt.Sprintf("Could not delete notification matcher %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *notificationMatcherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func computeMatcherDeletes(plan, state *notificationMatcherResourceModel) []string {
	if state == nil {
		return nil
	}

	var deletes []string

	if shouldDeleteListAttr(plan.MatchSeverity, state.MatchSeverity) {
		deletes = append(deletes, "match-severity")
	}
	if shouldDeleteListAttr(plan.MatchField, state.MatchField) {
		deletes = append(deletes, "match-field")
	}
	if shouldDeleteListAttr(plan.MatchCalendar, state.MatchCalendar) {
		deletes = append(deletes, "match-calendar")
	}
	if shouldDeleteStringAttr(plan.Comment, state.Comment) {
		deletes = append(deletes, "comment")
	}

	return deletes
}
