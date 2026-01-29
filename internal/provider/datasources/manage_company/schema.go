package managecompany

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *ManageCompanyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.NumberAttribute{
				Optional:            true,
				Description:         "Managed Company ID of selected managed company.",
				MarkdownDescription: "Managed Company ID of selected managed company.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed Company Name of selected managed company.",
				MarkdownDescription: "Managed Company Name of selected managed company.",
			},
			"node": schema.StringAttribute{
				Computed:            true,
				Description:         "Managing Node name or ID of selected managed company.",
				MarkdownDescription: "Managing Node name or ID of selected managed company.",
			},
			"plan": schema.StringAttribute{
				Computed:            true,
				Description:         "Base plan of selected managed company.",
				MarkdownDescription: "Base plan of selected managed company.",
			},
			"file_plan": schema.StringAttribute{
				Computed:            true,
				Description:         "File storage plan of selected managed company.",
				MarkdownDescription: "File storage plan of selected managed company.",
			},
		},
	}
}
