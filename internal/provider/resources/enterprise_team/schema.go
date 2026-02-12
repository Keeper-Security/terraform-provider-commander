// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseTeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages an enterprise team. Use this resource to create and manage teams in the MSP or Enterprise account",
		MarkdownDescription: "Manages an enterprise team. Use this resource to create and manage teams in the MSP or Enterprise account",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the enterprise team.",
				MarkdownDescription: "The ID of the enterprise team.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise Team Name.",
				MarkdownDescription: "Enterprise Team Name.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Team Name", 1, false),
				},
			},
			"restrict_record_edit": schema.BoolAttribute{
				Optional:            true,
				Description:         "Restrict record editing. Decide if users in this team can edit records",
				MarkdownDescription: "Restrict record editing. Decide if users in this team can edit records",
			},
			"restrict_record_re_share": schema.BoolAttribute{
				Optional:            true,
				Description:         "Restrict record re-sharing. Decide if users in this team can share records",
				MarkdownDescription: "Restrict record re-sharing. Decide if users in this team can share records",
			},
			"enable_privacy_screen": schema.BoolAttribute{
				Optional:            true,
				Description:         "Enable privacy screen. Decide if users in this team can view record passwords",
				MarkdownDescription: "Enable privacy screen. Decide if users in this team can view record passwords",
			},
			"users": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.SetNoEmptyStringsValidator("User"),
				},
				Description:         "Manage users to the enterprise team.",
				MarkdownDescription: "Manage users to the enterprise team.",
			},
			"roles": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.RolesValidator,
				},
				Description:         "Manage roles to the enterprise team.",
				MarkdownDescription: "Manage roles to the enterprise team.",
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Managing Node name or ID of the enterprise team.",
				MarkdownDescription: "Managing Node name or ID of the enterprise team.",
				Validators: []validator.String{
					utils.NodeValidator,
				},
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
		},
	}
}
