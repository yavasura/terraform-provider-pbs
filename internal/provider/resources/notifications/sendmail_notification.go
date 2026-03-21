/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package notifications

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	"github.com/yavasura/terraform-provider-pbs/pbs/notifications"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sendmailNotificationResource{}
	_ resource.ResourceWithConfigure   = &sendmailNotificationResource{}
	_ resource.ResourceWithImportState = &sendmailNotificationResource{}
)

// NewSendmailNotificationResource is a helper function to simplify the provider implementation.
func NewSendmailNotificationResource() resource.Resource {
	return &sendmailNotificationResource{}
}

// sendmailNotificationResource is the resource implementation.
type sendmailNotificationResource struct {
	client *pbs.Client
}

// sendmailNotificationResourceModel maps the resource schema data.
type sendmailNotificationResourceModel struct {
	Name       types.String `tfsdk:"name"`
	From       types.String `tfsdk:"from_address"`
	Mailto     types.List   `tfsdk:"mailto"`
	MailtoUser types.List   `tfsdk:"mailto_user"`
	Author     types.String `tfsdk:"author"`
	Comment    types.String `tfsdk:"comment"`
	Disable    types.Bool   `tfsdk:"disable"`
	Origin     types.String `tfsdk:"origin"`
}

// Metadata returns the resource type name.
func (r *sendmailNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sendmail_notification"
}

// Schema defines the schema for the resource.
func (r *sendmailNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Sendmail notification target for PBS alerts and notifications.",
		MarkdownDescription: `Manages a Sendmail notification target.

Configure local sendmail to receive notifications from PBS about backup jobs,
verification tasks, and system events.`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The unique name identifier for the Sendmail target.",
				MarkdownDescription: "The unique name identifier for the Sendmail notification target.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"from_address": schema.StringAttribute{
				Description:         "Sender email address.",
				MarkdownDescription: "Sender email address. This will appear as the 'From' address in notification emails.",
				Required:            true,
			},
			"mailto": schema.ListAttribute{
				Description:         "Recipient email address(es).",
				MarkdownDescription: "Recipient email address(es). Specify as a list of email strings.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"mailto_user": schema.ListAttribute{
				Description:         "User(s) from PBS user database to receive notifications.",
				MarkdownDescription: "User(s) from PBS user database to receive notifications. Specify as PBS user IDs.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"author": schema.StringAttribute{
				Description:         "Author name for notification emails.",
				MarkdownDescription: "Author name that will appear in the email headers.",
				Optional:            true,
			},
			"comment": schema.StringAttribute{
				Description:         "A comment describing this notification target.",
				MarkdownDescription: "A comment describing this notification target.",
				Optional:            true,
			},
			"disable": schema.BoolAttribute{
				Description:         "Disable this notification target.",
				MarkdownDescription: "Disable this notification target. Defaults to `false`.",
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
func (r *sendmailNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *sendmailNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sendmailNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	target := &notifications.SendmailTarget{
		Name: plan.Name.ValueString(),
		From: plan.From.ValueString(),
	}

	if !plan.Mailto.IsNull() && !plan.Mailto.IsUnknown() {
		target.Mailto = decodeStringList(ctx, plan.Mailto, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.MailtoUser.IsNull() && !plan.MailtoUser.IsUnknown() {
		target.MailtoUser = decodeStringList(ctx, plan.MailtoUser, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.Author.IsNull() {
		target.Author = plan.Author.ValueString()
	}
	if !plan.Comment.IsNull() {
		target.Comment = plan.Comment.ValueString()
	}
	if !plan.Disable.IsNull() {
		disable := plan.Disable.ValueBool()
		target.Disable = &disable
	}

	err := r.client.Notifications.CreateSendmailTarget(ctx, target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Sendmail notification target",
			fmt.Sprintf("Could not create Sendmail notification target %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	// Read back from API to get actual values
	created, err := r.client.Notifications.GetSendmailTarget(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created Sendmail notification target",
			fmt.Sprintf("Created Sendmail notification target %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	applySendmailNotificationState(ctx, created, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sendmailNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sendmailNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	target, err := r.client.Notifications.GetSendmailTarget(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Sendmail notification target",
			fmt.Sprintf("Could not read Sendmail notification target %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	applySendmailNotificationState(ctx, target, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sendmailNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sendmailNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	target := &notifications.SendmailTarget{
		Name: plan.Name.ValueString(),
		From: plan.From.ValueString(),
	}

	if plan.Mailto.IsNull() {
		target.Mailto = []string{}
	} else if !plan.Mailto.IsUnknown() {
		target.Mailto = decodeStringList(ctx, plan.Mailto, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if plan.MailtoUser.IsNull() {
		target.MailtoUser = []string{}
	} else if !plan.MailtoUser.IsUnknown() {
		target.MailtoUser = decodeStringList(ctx, plan.MailtoUser, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.Author.IsNull() {
		target.Author = plan.Author.ValueString()
	}
	if !plan.Comment.IsNull() {
		target.Comment = plan.Comment.ValueString()
	}
	if !plan.Disable.IsNull() {
		disable := plan.Disable.ValueBool()
		target.Disable = &disable
	}

	err := r.client.Notifications.UpdateSendmailTarget(ctx, plan.Name.ValueString(), target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Sendmail notification target",
			fmt.Sprintf("Could not update Sendmail notification target %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	updated, err := r.client.Notifications.GetSendmailTarget(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated Sendmail notification target",
			fmt.Sprintf("Updated Sendmail notification target %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	applySendmailNotificationState(ctx, updated, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sendmailNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sendmailNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Notifications.DeleteSendmailTarget(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Sendmail notification target",
			fmt.Sprintf("Could not delete Sendmail notification target %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *sendmailNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func applySendmailNotificationState(ctx context.Context, target *notifications.SendmailTarget, state *sendmailNotificationResourceModel, diags *diag.Diagnostics) {
	state.From = types.StringValue(target.From)

	mailto, listDiags := stringListState(ctx, target.Mailto)
	diags.Append(listDiags...)
	state.Mailto = mailto

	mailtoUser, listDiags := stringListState(ctx, target.MailtoUser)
	diags.Append(listDiags...)
	state.MailtoUser = mailtoUser

	state.Author = tfvalue.StringOrNull(target.Author)
	setNotificationCommonState(target.Comment, target.Disable, target.Origin, &state.Comment, &state.Disable, &state.Origin)
}
