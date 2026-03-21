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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	"github.com/yavasura/terraform-provider-pbs/pbs/notifications"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &smtpNotificationResource{}
	_ resource.ResourceWithConfigure   = &smtpNotificationResource{}
	_ resource.ResourceWithImportState = &smtpNotificationResource{}
)

// NewSMTPNotificationResource is a helper function to simplify the provider implementation.
func NewSMTPNotificationResource() resource.Resource {
	return &smtpNotificationResource{}
}

// smtpNotificationResource is the resource implementation.
type smtpNotificationResource struct {
	client *pbs.Client
}

// smtpNotificationResourceModel maps the resource schema data.
type smtpNotificationResourceModel struct {
	Name       types.String `tfsdk:"name"`
	Server     types.String `tfsdk:"server"`
	Port       types.Int64  `tfsdk:"port"`
	Mode       types.String `tfsdk:"mode"`
	From       types.String `tfsdk:"from_address"`
	Mailto     types.List   `tfsdk:"mailto"`
	MailtoUser types.List   `tfsdk:"mailto_user"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Author     types.String `tfsdk:"author"`
	Comment    types.String `tfsdk:"comment"`
	Disable    types.Bool   `tfsdk:"disable"`
	Origin     types.String `tfsdk:"origin"`
}

// Metadata returns the resource type name.
func (r *smtpNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smtp_notification"
}

// Schema defines the schema for the resource.
func (r *smtpNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an SMTP notification target for PBS alerts and notifications.",
		MarkdownDescription: `Manages an SMTP notification target.

Configure an SMTP server to receive notifications from PBS about backup jobs, 
verification tasks, and system events.`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The unique name identifier for the SMTP target.",
				MarkdownDescription: "The unique name identifier for the SMTP notification target.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server": schema.StringAttribute{
				Description:         "SMTP server hostname or IP address.",
				MarkdownDescription: "SMTP server hostname or IP address.",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				Description:         "SMTP server port.",
				MarkdownDescription: "SMTP server port. Common values: `25` (unencrypted), `465` (TLS), `587` (STARTTLS). Defaults to `25`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(25),
			},
			"mode": schema.StringAttribute{
				Description:         "Connection mode for SMTP.",
				MarkdownDescription: "Connection mode for SMTP. Valid values: `insecure` (no encryption), `starttls` (upgrade to TLS), `tls` (direct TLS). Defaults to `insecure`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("insecure"),
				Validators: []validator.String{
					stringvalidator.OneOf("insecure", "starttls", "tls"),
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
			"username": schema.StringAttribute{
				Description:         "SMTP authentication username.",
				MarkdownDescription: "SMTP authentication username. Required if the SMTP server requires authentication.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				Description:         "SMTP authentication password.",
				MarkdownDescription: "SMTP authentication password. Required if the SMTP server requires authentication.",
				Optional:            true,
				Sensitive:           true,
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
				Description:         "Origin of this configuration as reported by PBS (e.g., config file or built-in).",
				MarkdownDescription: "Origin of this configuration as reported by PBS (e.g., `user`, `builtin`).",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *smtpNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *smtpNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan smtpNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create SMTP target via API
	target := &notifications.SMTPTarget{
		Name:   plan.Name.ValueString(),
		Server: plan.Server.ValueString(),
		From:   plan.From.ValueString(),
	}

	// Set optional fields
	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt64())
		target.Port = &port
	}
	if !plan.Mode.IsNull() && !plan.Mode.IsUnknown() {
		target.Mode = plan.Mode.ValueString()
	}
	if plan.Mailto.IsNull() {
		target.To = []string{}
	} else if !plan.Mailto.IsUnknown() {
		target.To = decodeStringList(ctx, plan.Mailto, &resp.Diagnostics)
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
	if !plan.Username.IsNull() {
		target.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		target.Password = plan.Password.ValueString()
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

	err := r.client.Notifications.CreateSMTPTarget(ctx, target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SMTP notification target",
			fmt.Sprintf("Could not create SMTP notification target %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	// Read back from API to get computed values (like mode default)
	created, err := r.client.Notifications.GetSMTPTarget(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created SMTP notification target",
			fmt.Sprintf("Created SMTP notification target %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	applySMTPNotificationState(ctx, created, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if created.Mode == "" && (plan.Mode.IsUnknown() || plan.Mode.IsNull()) {
		plan.Mode = types.StringValue("insecure")
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *smtpNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state smtpNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get SMTP target from API
	target, err := r.client.Notifications.GetSMTPTarget(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SMTP notification target",
			fmt.Sprintf("Could not read SMTP notification target %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	applySMTPNotificationState(ctx, target, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	// Don't update password from API (sensitive field)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *smtpNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smtpNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update SMTP target via API
	target := &notifications.SMTPTarget{
		Name:   plan.Name.ValueString(),
		Server: plan.Server.ValueString(),
		From:   plan.From.ValueString(),
	}

	// Set optional fields
	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt64())
		target.Port = &port
	}
	if !plan.Mode.IsNull() {
		target.Mode = plan.Mode.ValueString()
	}
	if !plan.Mailto.IsNull() && !plan.Mailto.IsUnknown() {
		target.To = decodeStringList(ctx, plan.Mailto, &resp.Diagnostics)
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
	if !plan.Username.IsNull() {
		target.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		target.Password = plan.Password.ValueString()
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

	err := r.client.Notifications.UpdateSMTPTarget(ctx, plan.Name.ValueString(), target)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SMTP notification target",
			fmt.Sprintf("Could not update SMTP notification target %s: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	updated, err := r.client.Notifications.GetSMTPTarget(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated SMTP notification target",
			fmt.Sprintf("Updated SMTP notification target %s but could not read it back: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	applySMTPNotificationState(ctx, updated, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *smtpNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state smtpNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete SMTP target via API
	err := r.client.Notifications.DeleteSMTPTarget(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting SMTP notification target",
			fmt.Sprintf("Could not delete SMTP notification target %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *smtpNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func applySMTPNotificationState(ctx context.Context, target *notifications.SMTPTarget, state *smtpNotificationResourceModel, diags *diag.Diagnostics) {
	state.Server = types.StringValue(target.Server)
	state.From = types.StringValue(target.From)
	if target.Port != nil {
		state.Port = types.Int64Value(int64(*target.Port))
	} else {
		state.Port = types.Int64Null()
	}
	state.Mode = tfvalue.StringOrNull(target.Mode)

	to, listDiags := stringListState(ctx, target.To)
	diags.Append(listDiags...)
	state.Mailto = to

	mailtoUser, listDiags := stringListState(ctx, target.MailtoUser)
	diags.Append(listDiags...)
	state.MailtoUser = mailtoUser

	state.Username = tfvalue.StringOrNull(target.Username)
	state.Author = tfvalue.StringOrNull(target.Author)
	setNotificationCommonState(target.Comment, target.Disable, target.Origin, &state.Comment, &state.Disable, &state.Origin)
}
