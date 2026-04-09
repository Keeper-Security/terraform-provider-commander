// Copyright Keeper Security, Inc. 2026
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
		Description:         "Use this data source to look up an enterprise role by name / ID or managed company (MSP only) so you can reference it from other resources.",
		MarkdownDescription: "Use this data source to look up an enterprise role by **name** / **ID** or **managed company** (MSP only) so you can reference it from other resources.",
		Attributes: map[string]schema.Attribute{
			"role": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise role name or ID to find the role.",
				MarkdownDescription: "**Enterprise role name** or **ID** to find the role.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the found enterprise role.",
				MarkdownDescription: "**ID** of the found enterprise role.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the found enterprise role.",
				MarkdownDescription: "**Name** of the found enterprise role.",
			},
			"users": schema.SetAttribute{
				Computed:            true,
				Description:         "Users of the found enterprise role.",
				MarkdownDescription: "**Users** of the found enterprise role.",
				ElementType:         types.StringType,
			},
			"teams": schema.SetAttribute{
				Computed:            true,
				Description:         "Teams of the found enterprise role.",
				MarkdownDescription: "**Teams** of the found enterprise role.",
				ElementType:         types.StringType,
			},
			"managing_nodes": schema.MapAttribute{
				Computed:            true,
				Description:         "Managing nodes with privileges and cascade options of the found enterprise role. Map key is node name/ID, value is object with privileges (set of strings) and cascade (bool).",
				MarkdownDescription: "Managing nodes with **privileges** and **cascade** options of the found enterprise role. Map **key** is node name/ID, value is object with **privileges** (set of strings) and **cascade** (bool).",
				ElementType:         utils.ManagingNodesMapElemType,
			},
			"enforcement_policies": schema.MapAttribute{
				Computed:    true,
				Description: "Enforcement policies of the found enterprise role. Map key is policy name, value is the policy value string. For valid keys and value types, see the commander_enterprise_role resource documentation.",
				MarkdownDescription: "Enforcement policies of the found enterprise role. Map **key** is policy name, **value** is the policy value string.<br>" +
					"For valid keys and value types, see the <i>**commander_enterprise_role**</i> **resource** documentation.",
				ElementType: types.StringType,
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
			},
		},
	}
}
