package enterpriseuser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Enterprise user data source",
		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Description:         "Enterprise user email or ID to find the user.",
				Required:            true,
				MarkdownDescription: "Enterprise user email or ID to find the user.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise user.",
				MarkdownDescription: "ID of the found enterprise user.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise user.",
				MarkdownDescription: "Name of the found enterprise user.",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				Description:         "Email of the found enterprise user.",
				MarkdownDescription: "Email of the found enterprise user.",
			},
			"job_title": schema.StringAttribute{
				Computed:            true,
				Description:         "Job title of the found enterprise user.",
				MarkdownDescription: "Job title of the found enterprise user.",
			},
			"roles": schema.SetAttribute{
				Computed:            true,
				Description:         "Roles of the found enterprise user.",
				MarkdownDescription: "Roles of the found enterprise user.",
				ElementType:         types.StringType,
			},
			"teams": schema.SetAttribute{
				Computed:            true,
				Description:         "Teams of the found enterprise user.",
				MarkdownDescription: "Teams of the found enterprise user.",
				ElementType:         types.StringType,
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "Status of the found enterprise user.",
				MarkdownDescription: "Status of the found enterprise user.",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed company name or ID to scope the lookup (used for API context only; not returned in the result).",
				MarkdownDescription: "Managed company name or ID to scope the lookup (used for API context only; not returned in the result).",
			},
		},
	}
}
