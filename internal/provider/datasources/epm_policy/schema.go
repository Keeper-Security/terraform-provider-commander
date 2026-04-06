// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EpmPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Look up an existing EPM (Endpoint Policy Management) policy by its policy ID.",
		MarkdownDescription: "Look up an existing **EPM (Endpoint Policy Management) policy** by its **policy ID**.",
		Attributes: map[string]schema.Attribute{
			"policy": schema.StringAttribute{
				Required:            true,
				Description:         "EPM policy ID to look up (same value as the commander_epm_policy resource id).",
				MarkdownDescription: "**EPM policy ID** to look up (same value as the `commander_epm_policy` resource `id`).",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "EPM policy ID (same as policy).",
				MarkdownDescription: "EPM policy **ID** (same as `policy`).",
			},
			"policy_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Display name of the EPM policy.",
				MarkdownDescription: "**Display name** of the EPM policy.",
			},
			"policy_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy type. One of: " + commonepm.PolicyTypeDescription() + ".",
				MarkdownDescription: "Policy type. One of: " + commonepm.PolicyTypeMarkdown() + ".",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy status. One of: " + commonepm.StatusDescription() + ".",
				MarkdownDescription: "Policy **status**. One of: " + commonepm.StatusMarkdown() + ".",
			},
			"control": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Control actions. Each value is one of: " + commonepm.ControlDescription() + ".",
				MarkdownDescription: "**Control** actions. Each value is one of: " + commonepm.ControlMarkdown() + ".",
			},
			"user_groups": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "User collection IDs or \"*\" for all users.",
				MarkdownDescription: "User collection IDs or **`\"*\"`** for all users.",
			},
			"machine_collections": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Machine collection IDs or \"*\" for all machines.",
				MarkdownDescription: "Machine collection IDs or **`\"*\"`** for all machines.",
			},
			"applications": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Application collection IDs or \"*\" for all applications.",
				MarkdownDescription: "Application collection IDs or **`\"*\"`** for all applications.",
			},
			"day_filter": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Day filter. Each value is one of: " + commonepm.DayFilterDescription() + ".",
				MarkdownDescription: "**Day filter**. Each value is one of: " + commonepm.DayFilterMarkdown() + ".",
			},
			"time_filter": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Time filter ranges as start-end hours (0–23), e.g. 9-12.",
				MarkdownDescription: "**Time filter** ranges as **start-end** hours (**0–23**), e.g. `9-12`.",
			},
			"date_filter": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Date filter ranges in ISO format YYYY-MM-DD:YYYY-MM-DD.",
				MarkdownDescription: "**Date filter** ranges in **ISO format** `YYYY-MM-DD:YYYY-MM-DD`.",
			},
		},
	}
}
