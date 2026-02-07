// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *ManageCompanyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"managed_company": schema.StringAttribute{
				Required:            true,
				Description:         "Managed Company Name or ID to find the company.",
				MarkdownDescription: "Managed Company Name or ID to find the company.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found managed company.",
				MarkdownDescription: "ID of the found managed company.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found managed company.",
				MarkdownDescription: "Name of the found managed company.",
			},
			"node": schema.StringAttribute{
				Computed:            true,
				Description:         "Node of the found managed company.",
				MarkdownDescription: "Node of the found managed company.",
			},
			"plan": schema.StringAttribute{
				Computed:            true,
				Description:         "Base plan of the found managed company.",
				MarkdownDescription: "Base plan of the found managed company.",
			},
			"file_plan": schema.StringAttribute{
				Computed:            true,
				Description:         "File storage plan of the found managed company.",
				MarkdownDescription: "File storage plan of the found managed company.",
			},
		},
	}
}
