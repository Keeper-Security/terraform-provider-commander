// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (d *EnterpriseScimDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to look up an enterprise SCIM configuration by ID / Node ID or managed company (MSP only) so you can reference it from other resources.",
		MarkdownDescription: "Use this data source to look up an enterprise **SCIM configuration** by **ID** / **Node ID** or **managed company** (MSP only) so you can reference it from other resources.",
		Attributes: map[string]schema.Attribute{
			"scim": schema.StringAttribute{
				Required:            true,
				Description:         "SCIM configuration ID or Node ID to look up.",
				MarkdownDescription: "**SCIM configuration ID or Node ID** to look up.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("SCIM attribute requires either SCIM ID or Node ID", 1, false),
				},
			},
			"scim_id": schema.StringAttribute{
				Computed:            true,
				Description:         "SCIM configuration ID to look up.",
				MarkdownDescription: "**SCIM configuration ID** to look up.",
			},
			"scim_url": schema.StringAttribute{
				Computed:            true,
				Description:         "The SCIM endpoint URL for this configuration.",
				MarkdownDescription: "The SCIM **endpoint URL** for this configuration.",
			},
			"node_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The node that manages this SCIM configuration.",
				MarkdownDescription: "The **node** that manages this SCIM configuration.",
			},
			"node_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the node that manages this SCIM configuration.",
				MarkdownDescription: "The **name** of the node that manages this SCIM configuration.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "Current status of the SCIM configuration (for example, active or inactive).",
				MarkdownDescription: "**Current status** of the SCIM configuration (for example, `active` or `inactive`).",
			},
			"prefix": schema.StringAttribute{
				Computed:            true,
				Description:         "Role prefix. SCIM groups starting with this prefix are imported to Keeper as roles.",
				MarkdownDescription: "**Role prefix**. SCIM groups starting with this prefix are imported to Keeper as roles.",
			},
			"unique_groups": schema.BoolAttribute{
				Computed:            true,
				Description:         "Whether unique groups are used.",
				MarkdownDescription: "Whether **unique groups** are used.",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
			},
		},
	}
}
