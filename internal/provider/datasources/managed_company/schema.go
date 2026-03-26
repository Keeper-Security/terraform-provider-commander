// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *ManagedCompanyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to look up a managed company by name or ID so you can reference it from other resources.",
		MarkdownDescription: "Use this data source to look up a managed company by **name** or **ID** so you can reference it from other resources.",
		Attributes: map[string]schema.Attribute{
			"managed_company": schema.StringAttribute{
				Required:            true,
				Description:         "Managed Company Name or ID to find the company.",
				MarkdownDescription: "Managed Company **Name** or **ID** to find the company.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Company ID of the found managed company.",
				MarkdownDescription: "**Company ID** of the found managed company.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found managed company.",
				MarkdownDescription: "**Name** of the found managed company.",
			},
			"node": schema.StringAttribute{
				Computed:            true,
				Description:         "Node ID of the found managed company.",
				MarkdownDescription: "**Node ID** of the found managed company.",
			},
			"node_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Node name of the found managed company.",
				MarkdownDescription: "**Node name** of the found managed company.",
			},
			"plan": schema.StringAttribute{
				Computed:            true,
				Description:         "Keeper base plan of the found managed company.",
				MarkdownDescription: "Keeper **base plan** of the found managed company.",
			},
			"file_plan": schema.StringAttribute{
				Computed:            true,
				Description:         "Secure File Storage of the found managed company.",
				MarkdownDescription: "**Secure File Storage** of the found managed company.",
			},
			"seats": schema.Int64Attribute{
				Computed:            true,
				Description:         "Maximum number of user licenses of the found managed company.",
				MarkdownDescription: "Maximum number of **user licenses** of the found managed company.",
			},
			"add_ons": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "Secure Add-Ons of the found managed company.",
				MarkdownDescription: "**Secure Add-Ons** of the found managed company.",
			},
		},
	}
}
