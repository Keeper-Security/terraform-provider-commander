// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *EnterpriseNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a enterprise node. Use this resource to create and manage nodes in the MSP or Enterprise account.",
		MarkdownDescription: "Manages a enterprise node. Use this resource to create and manage nodes in the MSP or Enterprise account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the enterprise node.",
				MarkdownDescription: "ID of the enterprise node.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Set enterprise node's display name.",
				MarkdownDescription: "Set enterprise node's display name.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Node Name", 1, false),
				},
			},
			"parent": schema.StringAttribute{
				Required:            true,
				Description:         "Parent node name or ID. Make given node the parent of this node.",
				MarkdownDescription: "Parent node name or ID. Make given node the parent of this node.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true),
				},
			},
			"toggle_isolated": schema.BoolAttribute{
				Optional:            true,
				Description:         "Make node visible or invisible to people in other nodes.",
				MarkdownDescription: "Make node visible or invisible to people in other nodes.",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
		},
	}
}
