// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up an enterprise role by name or ID. Returns the role's ID, name, users, teams, managing nodes, and enforcement policies so you can reference them in other resources.",
		Attributes: map[string]schema.Attribute{
			"role": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise role name or ID to find the role.",
				MarkdownDescription: "Enterprise role name or ID to find the role.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise role.",
				MarkdownDescription: "ID of the found enterprise role.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise role.",
				MarkdownDescription: "Name of the found enterprise role.",
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
			"managing_nodes": schema.MapAttribute{
				Computed:            true,
				Description:         "Managing nodes (admin privileges) for the role. Map key is node name/ID; value is object with privileges (set of strings) and cascade (bool).",
				MarkdownDescription: "Managing nodes (admin privileges) for the role. Map key is node name/ID; value is object with `privileges` (set of strings) and `cascade` (bool).",
				ElementType:         utils.ManagingNodesMapElemType,
			},
			"enforcement_policies": schema.MapAttribute{
				Computed:            true,
				Description:         "Enforcement policies for the role. Map key is policy name (e.g. REQUIRE_TWO_FACTOR); value is the policy value string.",
				MarkdownDescription: "Enforcement policies for the role. Map key is policy name (e.g. `REQUIRE_TWO_FACTOR`); value is the policy value string.",
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
