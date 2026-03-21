/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tfschema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// RequiredReplaceStringAttribute builds a required string attribute that forces replacement.
func RequiredReplaceStringAttribute(description, markdown string, validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Description:         description,
		MarkdownDescription: markdown,
		Required:            true,
		Validators:          validators,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

// OptionalCommentAttribute builds a conventional optional comment field.
func OptionalCommentAttribute(description, markdown string) schema.StringAttribute {
	return schema.StringAttribute{
		Description:         description,
		MarkdownDescription: markdown,
		Optional:            true,
	}
}

// ComputedDigestAttribute builds a standard optimistic-lock digest field.
func ComputedDigestAttribute(description, markdown string) schema.StringAttribute {
	return schema.StringAttribute{
		Description:         description,
		MarkdownDescription: markdown,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// OptionalComputedBoolDefaultFalseAttribute builds a standard optional bool defaulting to false.
func OptionalComputedBoolDefaultFalseAttribute(description, markdown string) schema.BoolAttribute {
	return schema.BoolAttribute{
		Description:         description,
		MarkdownDescription: markdown,
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	}
}
