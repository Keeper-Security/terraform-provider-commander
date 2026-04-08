// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *EnterpriseNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages an enterprise node in your Keeper MSP or Enterprise account.<br><br>" +
			"Nodes organize users, roles, teams and administrators into distinct groupings, similar to organizational units in Active Directory. " +
			"You can structure nodes by location, department, division, or any other hierarchy that fits your organization. " +
			"The top-level node (Root) is set to the organization name; create child nodes under it or under other nodes as needed.<br><br>" +
			"For more information, see https://docs.keeper.io/en/enterprise-guide/getting-started-with-keeper-admin-console#nodes",
		MarkdownDescription: "Creates and manages an **enterprise node** in your Keeper MSP or Enterprise account.<br><br>" +
			"Nodes organize users, roles, teams and administrators into distinct groupings, similar to organizational units in Active Directory. " +
			"You can structure nodes by location, department, division, or any other hierarchy that fits your organization. " +
			"The top-level node (**Root**) is set to the organization name; create child nodes under it or under other nodes as needed.<br><br>" +
			"For more information, see [Enterprise Nodes documentation](https://docs.keeper.io/en/enterprise-guide/getting-started-with-keeper-admin-console#nodes).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Node ID assigned by Keeper to the node after it is created. " +
					"Use this value to import an existing node into Terraform state or to reference the node from other resources.",
				MarkdownDescription: "**Node ID** assigned by Keeper to the node after it is created. " +
					"Use this value to **import** an existing node into Terraform state or to reference the node from other resources.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Set the display name for the enterprise node. Must be at least one character.",
				MarkdownDescription: "Set the **display name** for the enterprise node. Must be at least one character.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Node Name", 1, false),
				},
			},
			"parent": schema.StringAttribute{
				Required:            true,
				Description:         "The parent node that will manage this enterprise node. Provide the node name or node ID. ",
				MarkdownDescription: "The **parent node** that will manage this enterprise node. Provide the **node name** or **node ID**. ",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true),
				},
			},
			"toggle_isolated": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         "When true, this node is isolated: users in this node cannot see or be seen by users in other nodes. When false, the node is visible across the organization. Defaults to false.<br>" + "Not supported on create; set on update to turn isolation on or off.",
				MarkdownDescription: "When `true`, this node is **isolated**: users in this node cannot see or be seen by users in other nodes. When `false`, the node is visible across the organization. Defaults to `false`.<br>" + "**Not supported** on **create**; set on **update** to turn isolation on or off.",
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
