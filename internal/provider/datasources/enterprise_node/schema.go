package enterprisenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *EnterpriseNodesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise node name or ID to find the node.",
				MarkdownDescription: "Enterprise node name or ID to find the node.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise node.",
				MarkdownDescription: "ID of the found enterprise node.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise node.",
				MarkdownDescription: "Name of the found enterprise node.",
			},
			"parent": schema.StringAttribute{
				Computed:            true,
				Description:         "Parent of selected enterprise node.",
				MarkdownDescription: "Parent of selected enterprise node.",
			},
			"parent_id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the parent of the found enterprise node.",
				MarkdownDescription: "ID of the parent of the found enterprise node.",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed company name or ID to scope the lookup (used for API context only; not returned in the result).",
				MarkdownDescription: "Managed company name or ID to scope the lookup (used for API context only; not returned in the result).",
			},
		},
	}
}
