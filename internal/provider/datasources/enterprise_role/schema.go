package enterpriserole

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise role name or ID to find the role.",
				MarkdownDescription: "Enterprise role name or ID to find the role.",
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
			"users": schema.SetAttribute{
				Computed:            true,
				Description:         "Users of the found enterprise role.",
				MarkdownDescription: "Users of the found enterprise role.",
				ElementType:         types.StringType,
			},
			"teams": schema.SetAttribute{
				Computed:            true,
				Description:         "Teams of the found enterprise role.",
				MarkdownDescription: "Teams of the found enterprise role.",
				ElementType:         types.StringType,
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed company name or ID to scope the lookup (used for API context only; not returned in the result).",
				MarkdownDescription: "Managed company name or ID to scope the lookup (used for API context only; not returned in the result).",
			},
		},
	}
}
