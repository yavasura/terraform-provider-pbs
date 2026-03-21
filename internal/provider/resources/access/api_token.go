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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfschema"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfstate"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

var (
	_ resource.Resource                = &apiTokenResource{}
	_ resource.ResourceWithConfigure   = &apiTokenResource{}
	_ resource.ResourceWithImportState = &apiTokenResource{}
)

var tokenNameRegex = regexp.MustCompile(`^[^\s!/[:cntrl:]]+$`)

// NewAPITokenResource creates a new PBS API token resource.
func NewAPITokenResource() resource.Resource {
	return &apiTokenResource{}
}

type apiTokenResource struct {
	client *pbs.Client
}

type apiTokenResourceModel struct {
	UserID    types.String `tfsdk:"userid"`
	TokenName types.String `tfsdk:"token_name"`
	TokenID   types.String `tfsdk:"tokenid"`
	Comment   types.String `tfsdk:"comment"`
	Enable    types.Bool   `tfsdk:"enable"`
	Expire    types.Int64  `tfsdk:"expire"`
	Digest    types.String `tfsdk:"digest"`
	Value     types.String `tfsdk:"value"`
}

func (r *apiTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_token"
}

func (r *apiTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PBS API token for an existing user account.",
		MarkdownDescription: "Manages a PBS API token for an existing user account.\n\n" +
			"PBS returns the token secret only when the token is created. This resource exposes that one-time value via `value` " +
			"as a sensitive computed attribute and preserves it in state for subsequent reads.",
		Attributes: map[string]schema.Attribute{
			"userid": tfschema.RequiredReplaceStringAttribute(
				"The PBS user ID that owns the token.",
				"The PBS user ID that owns the token, for example `backup-operator@pbs`.",
				stringvalidator.LengthBetween(3, 128),
				stringvalidator.RegexMatches(
					userIDRegex,
					"must be in the format 'username@realm'",
				),
			),
			"token_name": schema.StringAttribute{
				Description:         "Token name segment appended to the user ID.",
				MarkdownDescription: "Token name segment appended to the user ID, for example `terraform` in `backup-operator@pbs!terraform`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					stringvalidator.RegexMatches(
						tokenNameRegex,
						"must not contain spaces, '!' or '/'",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tokenid": schema.StringAttribute{
				Description:         "Full PBS token ID in `userid!token_name` format.",
				MarkdownDescription: "Full PBS token ID in `userid!token_name` format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"comment": schema.StringAttribute{
				Description:         "Optional token comment.",
				MarkdownDescription: "Optional token comment. PBS token metadata is treated as immutable here, so changes force replacement.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description:         "Whether the token is enabled.",
				MarkdownDescription: "Whether the token is enabled. PBS token metadata is treated as immutable here, so changes force replacement.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"expire": schema.Int64Attribute{
				Description:         "Token expiration time as a Unix timestamp.",
				MarkdownDescription: "Token expiration time as a Unix timestamp. PBS token metadata is treated as immutable here, so changes force replacement.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"digest": tfschema.ComputedDigestAttribute(
				"Opaque digest returned by PBS for optimistic locking.",
				"Opaque digest returned by PBS for optimistic locking.",
			),
			"value": schema.StringAttribute{
				Description:         "One-time token secret returned by PBS when the token is created.",
				MarkdownDescription: "One-time token secret returned by PBS when the token is created. PBS does not return this value again on subsequent reads.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *apiTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

func (r *apiTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiTokenResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	token := buildAPITokenFromModel(&plan)
	generated, err := r.client.Access.CreateUserToken(ctx, plan.UserID.ValueString(), plan.TokenName.ValueString(), token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating API token",
			fmt.Sprintf("Could not create token %s for user %s: %s", plan.TokenName.ValueString(), plan.UserID.ValueString(), err.Error()),
		)
		return
	}

	createdToken, err := r.client.Access.GetUserToken(ctx, plan.UserID.ValueString(), plan.TokenName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading API token",
			fmt.Sprintf("Could not read token %s for user %s after creation: %s", plan.TokenName.ValueString(), plan.UserID.ValueString(), err.Error()),
		)
		return
	}

	var state apiTokenResourceModel
	setAPITokenState(createdToken, generated.Value, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *apiTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiTokenResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	token, err := r.client.Access.GetUserToken(ctx, state.UserID.ValueString(), state.TokenName.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading API token",
			fmt.Sprintf("Could not read token %s for user %s: %s", state.TokenName.ValueString(), state.UserID.ValueString(), err.Error()),
		)
		return
	}

	secret := ""
	if !state.Value.IsNull() && !state.Value.IsUnknown() {
		secret = state.Value.ValueString()
	}

	setAPITokenState(token, secret, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *apiTokenResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unexpected API token update",
		"PBS API token metadata is treated as immutable in this provider. Change token settings by replacing the resource.",
	)
}

func (r *apiTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiTokenResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	digest := ""
	if !state.Digest.IsNull() && !state.Digest.IsUnknown() {
		digest = state.Digest.ValueString()
	}

	if err := r.client.Access.DeleteUserToken(ctx, state.UserID.ValueString(), state.TokenName.ValueString(), digest); err != nil {
		if isNotFoundError(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting API token",
			fmt.Sprintf("Could not delete token %s for user %s: %s", state.TokenName.ValueString(), state.UserID.ValueString(), err.Error()),
		)
	}
}

func (r *apiTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	userID, tokenName, ok := pbsaccess.SplitAPITokenID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format 'userid!token_name', for example 'backup-operator@pbs!terraform'.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("userid"), userID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("token_name"), tokenName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tokenid"), req.ID)...)
}

func buildAPITokenFromModel(model *apiTokenResourceModel) *pbsaccess.APIToken {
	token := &pbsaccess.APIToken{
		UserID:    model.UserID.ValueString(),
		TokenName: model.TokenName.ValueString(),
		TokenID:   pbsaccess.FormatAPITokenID(model.UserID.ValueString(), model.TokenName.ValueString()),
	}

	if !model.Comment.IsNull() && !model.Comment.IsUnknown() {
		token.Comment = model.Comment.ValueString()
	}
	if !model.Enable.IsNull() && !model.Enable.IsUnknown() {
		enable := model.Enable.ValueBool()
		token.Enable = &enable
	}
	if !model.Expire.IsNull() && !model.Expire.IsUnknown() {
		expire := model.Expire.ValueInt64()
		token.Expire = &expire
	}
	if !model.Digest.IsNull() && !model.Digest.IsUnknown() {
		token.Digest = model.Digest.ValueString()
	}

	return token
}

func setAPITokenState(token *pbsaccess.APIToken, secret string, state *apiTokenResourceModel) {
	state.UserID = types.StringValue(token.UserID)
	state.TokenName = types.StringValue(token.TokenName)
	state.TokenID = types.StringValue(token.TokenID)
	state.Comment = tfvalue.StringOrNull(token.Comment)
	state.Enable = tfvalue.BoolPtrOrNull(token.Enable)
	state.Expire = tfvalue.Int64PtrOrNull(token.Expire)
	state.Digest = tfvalue.StringOrNull(token.Digest)
	state.Value = tfvalue.StringOrNull(secret)
}
