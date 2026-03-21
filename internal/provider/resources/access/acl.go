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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/config"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfstate"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	"github.com/yavasura/terraform-provider-pbs/pbs"
	pbsaccess "github.com/yavasura/terraform-provider-pbs/pbs/access"
)

var (
	_ resource.Resource                = &aclResource{}
	_ resource.ResourceWithConfigure   = &aclResource{}
	_ resource.ResourceWithImportState = &aclResource{}
)

var (
	aclPathRegex  = regexp.MustCompile(`^(/|(/[A-Za-z0-9._-]+)+)$`)
	groupIDRegex  = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9._-]*$`)
	knownACLRoles = []string{
		"Admin",
		"Audit",
		"DatastoreAdmin",
		"DatastoreBackup",
		"DatastorePowerUser",
		"DatastoreReader",
		"NoAccess",
		"RemoteAdmin",
		"RemoteSyncOperator",
		"TapeAdmin",
		"TapeAudit",
		"TapeOperator",
	}
)

// NewACLResource creates a new PBS ACL resource.
func NewACLResource() resource.Resource {
	return &aclResource{}
}

type aclResource struct {
	client *pbs.Client
}

type aclResourceModel struct {
	Path      types.String `tfsdk:"path"`
	UGID      types.String `tfsdk:"ugid"`
	RoleID    types.String `tfsdk:"role_id"`
	Propagate types.Bool   `tfsdk:"propagate"`
}

func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl"
}

func (r *aclResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PBS access control list entry.",
		MarkdownDescription: "Manages a PBS access control list entry.\n\n" +
			"ACLs grant roles to users, groups, or API tokens on PBS paths. More specific paths override more general paths, " +
			"and `propagate = false` stops inheritance to child paths.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description:         "PBS ACL path such as `/`, `/datastore`, `/datastore/backups`, or `/remote/secondary`.",
				MarkdownDescription: "PBS ACL path such as `/`, `/datastore`, `/datastore/backups`, or `/remote/secondary`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(aclPathRegex, "must be a valid PBS ACL path"),
					aclPathValidator{},
				},
			},
			"ugid": schema.StringAttribute{
				Description:         "User, group, or token identifier receiving the role.",
				MarkdownDescription: "User, group, or token identifier receiving the role, for example `admin@pam`, `admin@pam!terraform`, or `backup-operators`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					ugidValidator{},
				},
			},
			"role_id": schema.StringAttribute{
				Description:         "PBS role to assign.",
				MarkdownDescription: "PBS role to assign.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(knownACLRoles...),
				},
			},
			"propagate": schema.BoolAttribute{
				Description:         "Whether the ACL propagates to child paths.",
				MarkdownDescription: "Whether the ACL propagates to child paths.",
				Required:            true,
			},
		},
	}
}

func (r *aclResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aclResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	acl := buildACLFromModel(&plan)
	if err := r.client.Access.SetACL(ctx, acl); err != nil {
		resp.Diagnostics.AddError(
			"Error creating ACL",
			fmt.Sprintf("Could not set ACL for %s on %s: %s", plan.UGID.ValueString(), plan.Path.ValueString(), err.Error()),
		)
		return
	}

	createdACL, err := r.client.Access.GetACLForRole(ctx, plan.Path.ValueString(), plan.UGID.ValueString(), plan.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ACL",
			fmt.Sprintf("Could not read ACL for %s on %s after creation: %s", plan.UGID.ValueString(), plan.Path.ValueString(), err.Error()),
		)
		return
	}

	var state aclResourceModel
	setACLState(createdACL, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *aclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aclResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	acl, err := r.client.Access.GetACLForRole(ctx, state.Path.ValueString(), state.UGID.ValueString(), state.RoleID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading ACL",
			fmt.Sprintf("Could not read ACL for %s on %s: %s", state.UGID.ValueString(), state.Path.ValueString(), err.Error()),
		)
		return
	}

	setACLState(acl, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aclResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	var state aclResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	if plan.RoleID.ValueString() != state.RoleID.ValueString() {
		if err := r.client.Access.DeleteACL(ctx, state.Path.ValueString(), state.UGID.ValueString(), state.RoleID.ValueString()); err != nil && !isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Error replacing ACL role",
				fmt.Sprintf("Could not remove previous ACL role %s for %s on %s: %s", state.RoleID.ValueString(), state.UGID.ValueString(), state.Path.ValueString(), err.Error()),
			)
			return
		}
	}

	acl := buildACLFromModel(&plan)
	if err := r.client.Access.SetACL(ctx, acl); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ACL",
			fmt.Sprintf("Could not set ACL for %s on %s: %s", plan.UGID.ValueString(), plan.Path.ValueString(), err.Error()),
		)
		return
	}

	updatedACL, err := r.client.Access.GetACLForRole(ctx, plan.Path.ValueString(), plan.UGID.ValueString(), plan.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ACL",
			fmt.Sprintf("Could not read ACL for %s on %s after update: %s", plan.UGID.ValueString(), plan.Path.ValueString(), err.Error()),
		)
		return
	}

	var newState aclResourceModel
	setACLState(updatedACL, &newState)
	tfstate.Encode(ctx, &resp.State, &newState, &resp.Diagnostics)
}

func (r *aclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aclResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	if err := r.client.Access.DeleteACL(ctx, state.Path.ValueString(), state.UGID.ValueString(), state.RoleID.ValueString()); err != nil {
		if isNotFoundError(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting ACL",
			fmt.Sprintf("Could not delete ACL for %s on %s: %s", state.UGID.ValueString(), state.Path.ValueString(), err.Error()),
		)
	}
}

func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	sep := strings.LastIndex(req.ID, ":")
	if sep <= 0 || sep == len(req.ID)-1 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format '{path}:{ugid}', for example '/datastore/backups:admin@pam'.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), req.ID[:sep])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ugid"), req.ID[sep+1:])...)
}

func buildACLFromModel(model *aclResourceModel) *pbsaccess.ACL {
	acl := &pbsaccess.ACL{
		Path:   model.Path.ValueString(),
		UGID:   model.UGID.ValueString(),
		RoleID: model.RoleID.ValueString(),
	}
	if !model.Propagate.IsNull() && !model.Propagate.IsUnknown() {
		propagate := model.Propagate.ValueBool()
		acl.Propagate = &propagate
	}
	return acl
}

func setACLState(acl *pbsaccess.ACL, state *aclResourceModel) {
	state.Path = types.StringValue(acl.Path)
	state.UGID = types.StringValue(acl.UGID)
	state.RoleID = types.StringValue(acl.RoleID)
	state.Propagate = tfvalue.BoolPtrOrDefault(acl.Propagate, true)
}
