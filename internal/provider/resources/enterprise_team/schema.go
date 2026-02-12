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
				Description:         "Restrict record editing.",
				MarkdownDescription: "Restrict record editing.",
			},
			"restrict_record_re_share": schema.BoolAttribute{
				Optional:            true,
				Description:         "Restrict record re-sharing.",
				MarkdownDescription: "Restrict record re-sharing.",
			},
			"enable_privacy_screen": schema.BoolAttribute{
				Optional:            true,
				Description:         "Enable privacy screen.",
				MarkdownDescription: "Enable privacy screen.",
			},
			"users": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.SetNoEmptyStringsValidator("User"),
				},
				Description:         "Set of users in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
				MarkdownDescription: "Set of users in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
			},
			"roles": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.RolesValidator,
				},
				Description:         "Set of roles in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
				MarkdownDescription: "Set of roles in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Managing Node name or ID.",
				MarkdownDescription: "Managing Node name or ID.",
				Validators: []validator.String{
					utils.NodeValidator,
				},
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed Company name or ID.",
				MarkdownDescription: "Managed Company name or ID.",
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
		},
	}
}
