// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to look up an enterprise user by email / ID or managed company (MSP only) so you can reference it from other resources.",
		MarkdownDescription: "Use this data source to look up an enterprise user by **email** / **ID** or **managed company** (MSP only) so you can reference it from other resources.",
		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Description:         "Enterprise user email or ID to find the user.",
				Required:            true,
				MarkdownDescription: "**Enterprise user email** or **ID** to find the user.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise user.",
				MarkdownDescription: "**ID** of the found enterprise user.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise user.",
				MarkdownDescription: "**Name** of the found enterprise user.",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				Description:         "Email of the found enterprise user.",
				MarkdownDescription: "**Email** of the found enterprise user.",
			},
			"job_title": schema.StringAttribute{
				Computed:            true,
				Description:         "Job title of the found enterprise user.",
				MarkdownDescription: "**Job title** of the found enterprise user.",
			},
			"roles": schema.SetAttribute{
				Computed:            true,
				Description:         "Roles of the found enterprise user.",
				MarkdownDescription: "**Roles** of the found enterprise user.",
				ElementType:         types.StringType,
			},
			"teams": schema.SetAttribute{
				Computed:            true,
				Description:         "Teams of the found enterprise user.",
				MarkdownDescription: "**Teams** of the found enterprise user.",
				ElementType:         types.StringType,
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
			},
		},
	}
}
