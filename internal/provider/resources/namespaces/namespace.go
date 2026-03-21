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
	pbsdatastores "github.com/yavasura/terraform-provider-pbs/pbs/datastores"
)

var (
	_ resource.Resource                = &namespaceResource{}
	_ resource.ResourceWithConfigure   = &namespaceResource{}
	_ resource.ResourceWithImportState = &namespaceResource{}
)

// NewNamespaceResource creates a new PBS namespace resource.
func NewNamespaceResource() resource.Resource {
	return &namespaceResource{}
}

type namespaceResource struct {
	client *pbs.Client
}

type namespaceResourceModel struct {
	Store     types.String `tfsdk:"store"`
	Namespace types.String `tfsdk:"namespace"`
	Comment   types.String `tfsdk:"comment"`
}

func (r *namespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *namespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a datastore namespace in PBS.",
		MarkdownDescription: "Manages a datastore namespace in PBS.\n\n" +
			"Namespace paths are hierarchical and limited to 7 segments. Parent namespaces must exist before a child namespace can be created.",
		Attributes: map[string]schema.Attribute{
			"store": schema.StringAttribute{
				Description:         "Datastore that owns the namespace.",
				MarkdownDescription: "Datastore that owns the namespace.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"namespace": schema.StringAttribute{
				Description:         "Namespace path within the datastore.",
				MarkdownDescription: "Namespace path within the datastore, for example `production/vms`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					namespacePathValidator{},
				},
			},
			"comment": schema.StringAttribute{
				Description:         "Optional namespace comment.",
				MarkdownDescription: "Optional namespace comment. Namespace comments are treated as immutable and changes force replacement.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *namespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config.ConfigureResourceClient(&r.client, req.ProviderData, &resp.Diagnostics)
}

func (r *namespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan namespaceResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	comment := ""
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		comment = plan.Comment.ValueString()
	}

	if err := r.client.Datastores.CreateNamespace(ctx, plan.Store.ValueString(), plan.Namespace.ValueString(), comment); err != nil {
		message := fmt.Sprintf("Could not create namespace %s in datastore %s: %s", plan.Namespace.ValueString(), plan.Store.ValueString(), err.Error())
		if strings.Contains(strings.ToLower(err.Error()), "parent") || strings.Contains(strings.ToLower(err.Error()), "not found") {
			message += "\n\nNamespaces are hierarchical. Create parent namespaces first, or use explicit Terraform dependencies to enforce parent-before-child ordering."
		}

		resp.Diagnostics.AddError("Error creating namespace", message)
		return
	}

	namespace, err := r.client.Datastores.GetNamespace(ctx, plan.Store.ValueString(), plan.Namespace.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading namespace",
			fmt.Sprintf("Could not read namespace %s in datastore %s after creation: %s", plan.Namespace.ValueString(), plan.Store.ValueString(), err.Error()),
		)
		return
	}

	var state namespaceResourceModel
	setNamespaceState(namespace, plan.Comment, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *namespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state namespaceResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	namespace, err := r.client.Datastores.GetNamespace(ctx, state.Store.ValueString(), state.Namespace.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading namespace",
			fmt.Sprintf("Could not read namespace %s in datastore %s: %s", state.Namespace.ValueString(), state.Store.ValueString(), err.Error()),
		)
		return
	}

	setNamespaceState(namespace, state.Comment, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *namespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan namespaceResourceModel
	if !tfstate.Decode(ctx, req.Plan, &plan, &resp.Diagnostics) {
		return
	}

	var state namespaceResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	namespace, err := r.client.Datastores.GetNamespace(ctx, state.Store.ValueString(), state.Namespace.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading namespace",
			fmt.Sprintf("Could not read namespace %s in datastore %s while refreshing state: %s", state.Namespace.ValueString(), state.Store.ValueString(), err.Error()),
		)
		return
	}

	setNamespaceState(namespace, plan.Comment, &state)
	tfstate.Encode(ctx, &resp.State, &state, &resp.Diagnostics)
}

func (r *namespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state namespaceResourceModel
	if !tfstate.Decode(ctx, req.State, &state, &resp.Diagnostics) {
		return
	}

	if err := r.client.Datastores.DeleteNamespace(ctx, state.Store.ValueString(), state.Namespace.ValueString()); err != nil {
		if isNotFoundError(err) {
			return
		}

		message := fmt.Sprintf("Could not delete namespace %s in datastore %s: %s", state.Namespace.ValueString(), state.Store.ValueString(), err.Error())
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "child") || strings.Contains(lower, "not empty") || strings.Contains(lower, "backup") {
			message += "\n\nPBS prevents deleting namespaces that still contain child namespaces or backups."
		}

		resp.Diagnostics.AddError("Error deleting namespace", message)
	}
}

func (r *namespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format 'store:namespace', for example 'backups:production/vms'.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("store"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("namespace"), parts[1])...)
}

func setNamespaceState(namespace *pbsdatastores.Namespace, comment types.String, state *namespaceResourceModel) {
	state.Store = types.StringValue(namespace.Store)
	state.Namespace = types.StringValue(namespace.Path)
	if !comment.IsNull() && !comment.IsUnknown() {
		state.Comment = comment
		return
	}
	state.Comment = tfvalue.StringOrNull(namespace.Comment)
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "404") || strings.Contains(msg, "not found")
}
