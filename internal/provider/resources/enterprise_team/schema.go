// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseTeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Creates and manages an enterprise team in your Keeper MSP or Enterprise account.<br><br>" + "The purpose of creating Teams is to give users the ability to share the records and folders within their vaults with logical groupings of individuals. The administrator simply creates the team, sets any Team Restrictions (edit/viewing/sharing of passwords) and adds individual users to the team. Teams can also be used to easily assign Roles to entire groups of users to ensure the consistency of enforcement policies across a collective group of individuals.<br><br>" + "For more information, see https://docs.keeper.io/en/enterprise-guide/getting-started-with-keeper-admin-console#teams",
		MarkdownDescription: "Creates and manages an enterprise team in your Keeper MSP or Enterprise account.<br><br>" + "The purpose of creating Teams is to give users the ability to share the records and folders within their vaults with logical groupings of individuals. The administrator simply creates the team, sets any Team Restrictions (edit/viewing/sharing of passwords) and adds individual users to the team. Teams can also be used to easily assign Roles to entire groups of users to ensure the consistency of enforcement policies across a collective group of individuals.<br><br>" + "For more information, see [Enterprise Teams documentation](https://docs.keeper.io/en/enterprise-guide/getting-started-with-keeper-admin-console#teams).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Team ID assigned by Keeper to the team after it is created. " +
					"Use this value to import an existing team into Terraform state or to reference the team from other resources.",
				MarkdownDescription: "**Team ID** assigned by Keeper to the team after it is created. " +
					"Use this value to **import** an existing team into Terraform state or to reference the team from other resources.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Set the display name for the enterprise team. Must be at least one character.",
				MarkdownDescription: "Set the **display name** for the enterprise team. Must be at least **one character**.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Team Name", 1, false),
				},
			},
			"restrict_record_edit": schema.BoolAttribute{
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				Description:         "Restrict record editing. Decide if users in this team can edit records. Defaults to false.",
				MarkdownDescription: "Restrict record editing. Decide if **users in this team can edit records**. Defaults to `false`.",
			},
			"restrict_record_re_share": schema.BoolAttribute{
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				Description:         "Restrict record re-sharing. Decide if users in this team can share records. Defaults to false.",
				MarkdownDescription: "Restrict record re-sharing. Decide if **users in this team can share records**. Defaults to `false`.",
			},
			"enable_privacy_screen": schema.BoolAttribute{
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				Description:         "Enable privacy screen. Decide if users in this team can view record passwords. Defaults to false.",
				MarkdownDescription: "Enable privacy screen. Decide if **users in this team can view record passwords**. Defaults to `false`.",
			},
			"users": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.SetNoEmptyStringsValidator("User"),
				},
				Description:         "Set of users assigned to this enterprise team. Provide user email addresses or user IDs. ",
				MarkdownDescription: "Set of **users** assigned to this enterprise team. Provide **user email addresses** or **user IDs**. ",
			},
			"roles": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.RolesValidator,
				},
				Description:         "Set of roles assigned to this enterprise team. Provide role names or role IDs. ",
				MarkdownDescription: "Set of **roles** assigned to this enterprise team. Provide **role names** or **role IDs**. ",
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "The node that will manage this enterprise team. Provide the node name or node ID. ",
				MarkdownDescription: "The **node** that will manage this enterprise team. Provide the **node name** or **node ID**. ",
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
