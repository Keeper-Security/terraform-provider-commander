// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *EnterpriseScimResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Creates and manages an enterprise SCIM configuration in your Keeper Enterprise or MSP account.<br><br>" + "Automatically provision users and teams through Entra ID, Okta, Google Workspace and other popular identity platforms by establishing a SCIM connection.<br><br>" + "For more information, see https://docs.keeper.io/en/enterprise-guide/user-and-team-provisioning/automated-provisioning-with-scim#what-is-scim",
		MarkdownDescription: "Creates and manages an enterprise SCIM configuration in your Keeper Enterprise or MSP account.<br><br>" + "Automatically provision users and teams through Entra ID, Okta, Google Workspace and other popular identity platforms by establishing a SCIM connection.<br><br>" + "For more information, see [KeeperSCIM documentation](https://docs.keeper.io/en/enterprise-guide/user-and-team-provisioning/automated-provisioning-with-scim#what-is-scim).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "SCIM configuration ID assigned by Keeper to the scim after it is created.",
				MarkdownDescription: "SCIM configuration **ID** assigned by Keeper to the scim after it is created.",
			},
			"scim_url": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "The SCIM endpoint URL for this configuration. Use this URL to configure your SCIM client.",
				MarkdownDescription: "The SCIM **endpoint URL** for this configuration. Use this URL to configure your SCIM client.",
			},
			"node": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.NodeValidator,
				},
				Description:         "The node that will manage this SCIM configuration. Provide the node name or node ID. ",
				MarkdownDescription: "The **node** that will manage this SCIM configuration. Provide the **node name** or **node ID**. ",
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Current status of the SCIM configuration (for example, active or inactive)",
				MarkdownDescription: "**Current status** of the SCIM configuration (for example, `active` or `inactive`)",
			},
			"prefix": schema.StringAttribute{
				Required:            true,
				Description:         "Role Prefix. SCIM groups staring with prefix will be imported to Keeper as Roles.",
				MarkdownDescription: "**Role Prefix**. SCIM groups staring with prefix will be imported to Keeper as Roles.",
			},
			"unique_groups": schema.BoolAttribute{
				Required:            true,
				Description:         "Whether to use unique groups.",
				MarkdownDescription: "Whether to use **unique groups**.",
			},
			"provisioning_token": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Provisioning token for the SCIM configuration. Use this token to configure your SCIM client.",
				MarkdownDescription: "Provisioning token for the SCIM configuration. Use this token to configure your SCIM client.",
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
