package enterpriseteam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseTeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"team": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise team name or ID to find the team.",
				MarkdownDescription: "Enterprise team name or ID to find the team.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise team.",
				MarkdownDescription: "ID of the found enterprise team.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise team.",
				MarkdownDescription: "Name of the found enterprise team.",
			},
			"users": schema.SetAttribute{
				Computed:            true,
				Description:         "Users of the found enterprise team.",
				MarkdownDescription: "Users of the found enterprise team.",
				ElementType:         types.StringType,
			},
			"roles": schema.SetAttribute{
				Computed:            true,
				Description:         "Roles of the found enterprise team.",
				MarkdownDescription: "Roles of the found enterprise team.",
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
