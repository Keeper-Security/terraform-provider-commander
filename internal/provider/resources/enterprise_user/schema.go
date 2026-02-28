// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description:         "Creates and manages an enterprise user in your Keeper MSP or Enterprise account.<br><br>" + "All employees or users you choose to deploy Keeper to are responsible for managing their own encrypted vault. Every user's vault can be made up of private records or shared records. Users can be provisioned many different ways. Users can be required to set up a Master Password or they can be provisioned and authenticated through your SSO provider.<br><br>" + "For more information, see https://docs.keeper.io/en/enterprise-guide/getting-started-with-keeper-admin-console#users",
		MarkdownDescription: "Creates and manages an **enterprise user** in your Keeper MSP or Enterprise account.<br><br>" + "All employees or users you choose to deploy Keeper to are responsible for managing their own encrypted vault. Every user's vault can be made up of private records or shared records. Users can be provisioned many different ways. Users can be required to set up a Master Password or they can be provisioned and authenticated through your SSO provider.<br><br>" + "For more information, see [Enterprise Users documentation](https://docs.keeper.io/en/enterprise-guide/getting-started-with-keeper-admin-console#users).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "User ID assigned by Keeper to the role after it is created. Use this value to import an existing user into Terraform state or to reference the user from other resources.",
				MarkdownDescription: "**User ID** assigned by Keeper to the role after it is created. Use this value to **import** an existing user into Terraform state or to reference the user from other resources.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				Description:         "Set the email address for the enterprise user.",
				MarkdownDescription: "Set the **email address** for the enterprise user.",
				Validators: []validator.String{
					emailValidator{},
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the display name for the enterprise user. Must be at least one character.",
				MarkdownDescription: "Set the **display name** for the enterprise user. Must be at least **one character**.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Name", 1, true),
				},
			},
			"job_title": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the job title for the enterprise user. Must be at least one character.",
				MarkdownDescription: "Set the **job title** for the enterprise user. Must be at least **one character**.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Job title", 1, true),
				},
			},
			"roles": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Set of roles assigned to this enterprise user. Provide role names or role IDs. ",
				MarkdownDescription: "Set of **roles** assigned to this enterprise user. Provide **role names** or **role IDs**. ",
				Validators: []validator.Set{
					utils.RolesValidator,
				},
			},
			"teams": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Set of teams assigned to this enterprise user. Provide team names or team IDs. ",
				MarkdownDescription: "Set of **teams** assigned to this enterprise user. Provide **team names** or **team IDs**. ",
				Validators: []validator.Set{
					utils.TeamsValidator,
				},
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "The node that will manage this enterprise user. Provide the node name or node ID. ",
				MarkdownDescription: "The **node** that will manage this enterprise user. Provide the **node name** or **node ID**. ",
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
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "User current status (e.g. Active, Inactive). Set by the provider from the API and used for internal tracking; do not set in configuration.",
				MarkdownDescription: "User **current status** (e.g. Active, Inactive). Set by the provider from the API and used for internal tracking; **do not set in configuration**.",
			},
		},
	}
}
