// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *EnterpriseNodesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to look up an enterprise node by name / ID or managed company (MSP only) so you can reference it from other resources.",
		MarkdownDescription: "Use this data source to look up an enterprise node by **name** / **ID** or **managed company** (MSP only) so you can reference it from other resources.",
		Attributes: map[string]schema.Attribute{
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise node name or ID to find the node.",
				MarkdownDescription: "**Enterprise node name** or **ID** to find the node.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise node.",
				MarkdownDescription: "**ID** of the found enterprise node.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise node.",
				MarkdownDescription: "**Name** of the found enterprise node.",
			},
			"parent": schema.StringAttribute{
				Computed:            true,
				Description:         "Parent node name of the found enterprise node.",
				MarkdownDescription: "**Parent node name** of the found enterprise node.",
			},
			"parent_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Parent node ID of the found enterprise node.",
				MarkdownDescription: "**Parent node ID** of the found enterprise node.",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
			},
		},
	}
}
