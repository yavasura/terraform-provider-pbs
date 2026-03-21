/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfplan"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfschema"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfstate"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

var userIDRegex = regexp.MustCompile(`^[^\s@/[:cntrl:]]+@[A-Za-z0-9_][A-Za-z0-9._\-]*$`)

// NewUserResource creates a new PBS user resource.
func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *pbs.Client
}

type userResourceModel struct {
	UserID    types.String `tfsdk:"userid"`
	Comment   types.String `tfsdk:"comment"`
	Enable    types.Bool   `tfsdk:"enable"`
	Expire    types.Int64  `tfsdk:"expire"`
	FirstName types.String `tfsdk:"firstname"`
	LastName  types.String `tfsdk:"lastname"`
	Email     types.String `tfsdk:"email"`
	Digest    types.String `tfsdk:"digest"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PBS user account.",
		MarkdownDescription: "Manages a PBS user account.\n\n" +
			"This resource manages account metadata exposed by `/access/users`, including enablement, expiration,\n" +
			"and directory profile information. It does not manage credentials or passwords.",
		Attributes: map[string]schema.Attribute{
			"userid": tfschema.RequiredReplaceStringAttribute(
				"The full PBS user ID in `username@realm` format.",
				"The full PBS user ID in `username@realm` format, for example `john@ldap` or `admin@pam`.",
				stringvalidator.LengthBetween(3, 128),
				stringvalidator.RegexMatches(
					userIDRegex,
					"must be in the format 'username@realm'",
				),
			),
			"comment": tfschema.OptionalCommentAttribute(
				"Optional comment for the user account.",
				"Optional comment for the user account.",
			),
			"enable": schema.BoolAttribute{
				Description:         "Whether the account is enabled.",
				MarkdownDescription: "Whether the account is enabled. Omit this to let PBS keep its current/default behavior.",
				Optional:            true,
			},
			"expire": schema.Int64Attribute{
				Description:         "User account expiration time as a Unix timestamp.",
				MarkdownDescription: "User account expiration time as a Unix timestamp. Omit to leave the account non-expiring or unchanged.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"firstname": schema.StringAttribute{
				Description:         "Given name for the user.",
				MarkdownDescription: "Given name for the user.",
				Optional:            true,
			},
			"lastname": schema.StringAttribute{
				Description:         "Family name for the user.",
				MarkdownDescription: "Family name for the user.",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				Description:         "Email address associated with the user.",
				MarkdownDescription: "Email address associated with the user.",
				Optional:            true,
			},
			"digest": tfschema.ComputedDigestAttribute(
				"Opaque digest returned by PBS for optimistic locking.",
				"Opaque digest returned by PBS for optimistic locking.",
			),
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	user := buildUserFromModel(&plan)
	if err := r.client.Access.CreateUser(ctx, user); err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			fmt.Sprintf("Could not create user %s: %s", plan.UserID.ValueString(), err.Error()),
		)
		return
	}

	createdUser, err := r.client.Access.GetUser(ctx, plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user %s after creation: %s", plan.UserID.ValueString(), err.Error()),
		)
		return
	}

	var state userResourceModel
	setUserState(createdUser, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	user, err := r.client.Access.GetUser(ctx, state.UserID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user %s: %s", state.UserID.ValueString(), err.Error()),
		)
		return
	}

	setUserState(user, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	var state userResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	if (plan.Digest.IsNull() || plan.Digest.IsUnknown()) && !state.Digest.IsNull() && !state.Digest.IsUnknown() {
		plan.Digest = state.Digest
	}

	user := buildUserFromModel(&plan)
	user.Delete = computeUserDeletes(&plan, &state)

	if err := r.client.Access.UpdateUser(ctx, plan.UserID.ValueString(), user); err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			fmt.Sprintf("Could not update user %s: %s", plan.UserID.ValueString(), err.Error()),
		)
		return
	}

	updatedUser, err := r.client.Access.GetUser(ctx, plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user %s after update: %s", plan.UserID.ValueString(), err.Error()),
		)
		return
	}

	setUserState(updatedUser, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	digest := ""
	if !state.Digest.IsNull() && !state.Digest.IsUnknown() {
		digest = state.Digest.ValueString()
	}

	if err := r.client.Access.DeleteUser(ctx, state.UserID.ValueString(), digest); err != nil {
		if isNotFoundError(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting user",
			fmt.Sprintf("Could not delete user %s: %s", state.UserID.ValueString(), err.Error()),
		)
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("userid"), req, resp)
}

func buildUserFromModel(model *userResourceModel) *pbsaccess.User {
	user := &pbsaccess.User{
		UserID: model.UserID.ValueString(),
	}

	if !model.Comment.IsNull() && !model.Comment.IsUnknown() {
		user.Comment = model.Comment.ValueString()
	}
	if !model.Enable.IsNull() && !model.Enable.IsUnknown() {
		enable := model.Enable.ValueBool()
		user.Enable = &enable
	}
	if !model.Expire.IsNull() && !model.Expire.IsUnknown() {
		expire := model.Expire.ValueInt64()
		user.Expire = &expire
	}
	if !model.FirstName.IsNull() && !model.FirstName.IsUnknown() {
		user.FirstName = model.FirstName.ValueString()
	}
	if !model.LastName.IsNull() && !model.LastName.IsUnknown() {
		user.LastName = model.LastName.ValueString()
	}
	if !model.Email.IsNull() && !model.Email.IsUnknown() {
		user.Email = model.Email.ValueString()
	}
	if !model.Digest.IsNull() && !model.Digest.IsUnknown() {
		user.Digest = model.Digest.ValueString()
	}

	return user
}

func setUserState(user *pbsaccess.User, state *userResourceModel) {
	state.UserID = types.StringValue(user.UserID)
	state.Comment = tfvalue.StringOrNull(user.Comment)
	state.Enable = tfvalue.BoolPtrOrNull(user.Enable)
	state.Expire = tfvalue.Int64PtrOrNull(user.Expire)
	state.FirstName = tfvalue.StringOrNull(user.FirstName)
	state.LastName = tfvalue.StringOrNull(user.LastName)
	state.Email = tfvalue.StringOrNull(user.Email)
	state.Digest = types.StringValue(user.Digest)
}

func computeUserDeletes(plan, state *userResourceModel) []string {
	var deletes tfplan.DeleteSet

	deletes.AddIf("comment", tfplan.ShouldDeleteString(plan.Comment, state.Comment))
	deletes.AddIf("enable", tfplan.ShouldDeleteBool(plan.Enable, state.Enable))
	deletes.AddIf("expire", tfplan.ShouldDeleteInt64(plan.Expire, state.Expire))
	deletes.AddIf("firstname", tfplan.ShouldDeleteString(plan.FirstName, state.FirstName))
	deletes.AddIf("lastname", tfplan.ShouldDeleteString(plan.LastName, state.LastName))
	deletes.AddIf("email", tfplan.ShouldDeleteString(plan.Email, state.Email))

	return deletes.Values()
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "404") || strings.Contains(msg, "not found")
}
