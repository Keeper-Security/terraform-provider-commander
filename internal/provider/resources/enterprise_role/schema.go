// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages an enterprise role. Use this resource to create and manage roles in the MSP or Enterprise account",
		MarkdownDescription: "Manages an enterprise role. Use this resource to create and manage roles in the MSP or Enterprise account",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the enterprise role.",
				MarkdownDescription: "ID of the enterprise role.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise role name.",
				MarkdownDescription: "Enterprise role name.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Role Name", 1, false),
				},
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Managing node name or ID of the enterprise role.",
				MarkdownDescription: "Managing node name or ID of the enterprise role.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Node", 1, true),
				},
			},
			"users": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Manage users to the enterprise role.",
				MarkdownDescription: "Manage users to the enterprise role.",
				Validators: []validator.Set{
					utils.SetNoEmptyStringsValidator("User"),
				},
			},
			"teams": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Manage teams to the enterprise role.",
				MarkdownDescription: "Manage teams to the enterprise role.",
				Validators: []validator.Set{
					utils.TeamsValidator,
				},
			},
			"managing_nodes": schema.MapNestedAttribute{
				Optional: true,
				Validators: []validator.Map{
					utils.MapKeysMinLengthValidator("managing node name", 1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"privileges": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								privilegesValidator{},
							},
							Description:         "Manage privileges to grant for this managing node. Valid values: manage_nodes, manage_user, manage_roles, manage_teams, run_reports, manage_bridge, approve_device, manage_record_types, run_compliance_reports, transfer_account, sharing_administrator, manage_companies",
							MarkdownDescription: "Manage privileges to grant for this managing node. Valid values: `manage_nodes`, `manage_user`, `manage_roles`, `manage_teams`, `run_reports`, `manage_bridge`, `approve_device`, `manage_record_types`, `run_compliance_reports`, `transfer_account`, `sharing_administrator`, `manage_companies`",
						},
						"cascade": schema.BoolAttribute{
							Optional:            true,
							Description:         "Manage extending admin-privileges for the specified role(s) to child nodes as well ",
							MarkdownDescription: "Manage extending admin-privileges for the specified role(s) to child nodes as well ",
						},
					},
				},
				Description:         "Manage an admin privilege to the role. The map key is the node name/ID. Each managing node has `privileges` (optional) and `cascade` (optional) fields.",
				MarkdownDescription: "Manage an admin privilege to the role. The map key is the node name/ID. Each managing node has `privileges` (optional) and `cascade` (optional) fields.",
			},
			"enforcement_policies": schema.MapAttribute{
				Optional: true,
				Validators: []validator.Map{
					enforcementPoliciesMapKeyValidator{},
				},
				PlanModifiers: []planmodifier.Map{
					enforcementPoliciesGPCPlanModifier{},
				},
				ElementType:         types.StringType,
				Description:         "Manage enforcement policies for the given role. The map key is the enforcement policy key (e.g., REQUIRE_TWO_FACTOR) and the value is the policy value (e.g., false). For GENERATED_PASSWORD_COMPLEXITY use a JSON string (e.g. jsonencode([...])).",
				MarkdownDescription: "Manage enforcement policies for the given role. The map key is the enforcement policy key (e.g., `REQUIRE_TWO_FACTOR`) and the value is the policy value (e.g., `\"false\"`). For `GENERATED_PASSWORD_COMPLEXITY` use a JSON string (e.g. `jsonencode([...])`).",
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
