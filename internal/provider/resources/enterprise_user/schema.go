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
		Description:         "Manages an enterprise user. Use this resource to create and manage users in the MSP or Enterprise account",
		MarkdownDescription: "Manages an enterprise user. Use this resource to create and manage users in the MSP or Enterprise account",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the user.",
				MarkdownDescription: "The ID of the user.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				Description:         "Email address of the enterprise user.",
				MarkdownDescription: "Email address of the enterprise user.",
				Validators: []validator.String{
					emailValidator{},
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Name of the enterprise user.",
				MarkdownDescription: "Name of the enterprise user.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Name", 1, true),
				},
			},
			"job_title": schema.StringAttribute{
				Optional:            true,
				Description:         "Job title of the enterprise user.",
				MarkdownDescription: "Job title of the enterprise user.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Job title", 1, true),
				},
			},
			"roles": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Manage roles to the enterprise user.",
				MarkdownDescription: "Manage roles to the enterprise user.",
				Validators: []validator.Set{
					utils.RolesValidator,
				},
			},
			"teams": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Manage teams to the enterprise user.",
				MarkdownDescription: "Manage teams to the enterprise user.",
				Validators: []validator.Set{
					utils.TeamsValidator,
				},
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Managing node name or ID of the enterprise user.",
				MarkdownDescription: "Managing node name or ID of the enterprise user.",
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
				Description:         "User status (e.g. Active, Inactive). Set by the provider from the API; do not set in configuration.",
				MarkdownDescription: "User status (e.g. Active, Inactive). Set by the provider from the API; do not set in configuration.",
			},
		},
	}
}
