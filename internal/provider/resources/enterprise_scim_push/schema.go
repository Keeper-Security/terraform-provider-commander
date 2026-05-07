// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *EnterpriseScimPushResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Push SCIM (Google, AD, or record-based) data to a Keeper SCIM endpoint in one step. Changing any value runs the push again.",
		MarkdownDescription: "Use this resource to **push SCIM data** to a Keeper SCIM endpoint in a single step. You choose where the data comes from (Google Workspace, Active Directory, or a record) and whether to auto-approve teams.\n\n" +
			"## What this resource does\n\n" +
			"- When you **apply**, it runs a one-time SCIM push with your settings.\n" +
			"- Terraform does **not** read or delete anything on the server. It only runs the push and tracks that it happened.\n" +
			"- If you change **scim_id**, **source**, **record**, **auto_approve**, or **managed_company**, Terraform will run the push again with the new values.\n\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Deterministic ID computed from scim_id, source, record, auto_approve, and managed_company.",
				MarkdownDescription: "Deterministic ID computed from scim_id, source, record, auto_approve, and managed_company.",
			},
			"scim_id": schema.StringAttribute{
				Required:            true,
				Description:         "SCIM ID",
				MarkdownDescription: "SCIM ID",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Scim ID", 1, false),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				Required:            true,
				Description:         "Source of SCIM data. Must be one of: google, ad, record.",
				MarkdownDescription: "Source of SCIM data. Must be one of: `google`, `ad`, `record`.",
				Validators: []validator.String{
					sourceValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"record": schema.StringAttribute{
				Required:            true,
				Description:         "Record UID with SCIM configuration",
				MarkdownDescription: "**Record UID** with SCIM configuration",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Record", 1, false),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_approve": schema.BoolAttribute{
				Required:            true,
				Description:         "Auto approve SCIM teams",
				MarkdownDescription: "**Auto approve** SCIM teams",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
