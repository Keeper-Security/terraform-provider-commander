// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseTeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to look up an enterprise team by name / ID or managed company (MSP only) so you can reference it from other resources.",
		MarkdownDescription: "Use this data source to look up an enterprise team by **name** / **ID** or **managed company** (MSP only) so you can reference it from other resources.",
		Attributes: map[string]schema.Attribute{
			"team": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise team name or ID to find the team.",
				MarkdownDescription: "**Enterprise team name** or **ID** to find the team.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise team.",
				MarkdownDescription: "**ID** of the found enterprise team.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise team.",
				MarkdownDescription: "**Name** of the found enterprise team.",
			},
			"users": schema.SetAttribute{
				Computed:            true,
				Description:         "Users of the found enterprise team.",
				MarkdownDescription: "**Users** of the found enterprise team.",
				ElementType:         types.StringType,
			},
			"roles": schema.SetAttribute{
				Computed:            true,
				Description:         "Roles of the found enterprise team.",
				MarkdownDescription: "**Roles** of the found enterprise team.",
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
