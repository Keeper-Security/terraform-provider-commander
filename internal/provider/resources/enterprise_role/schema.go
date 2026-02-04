package enterpriserole

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					nameValidator{},
				},
			},
			"node": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					nodeValidator{},
				},
			},
			"users": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					usersValidator{},
				},
			},
			"teams": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					teamsValidator{},
				},
			},
			"managing_nodes": schema.MapNestedAttribute{
				Optional: true,
				Validators: []validator.Map{
					managingNodesMapValidator{},
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"privileges": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								privilegesValidator{},
							},
							Description:         "Set of privileges to grant for this managing node. Valid values: manage_nodes, manage_user, manage_roles, manage_teams, run_reports, manage_bridge, approve_device, manage_record_types, run_compliance_reports, transfer_account, sharing_administrator",
							MarkdownDescription: "Set of privileges to grant for this managing node. Valid values: `manage_nodes`, `manage_user`, `manage_roles`, `manage_teams`, `run_reports`, `manage_bridge`, `approve_device`, `manage_record_types`, `run_compliance_reports`, `transfer_account`, `sharing_administrator`",
						},
						"cascade": schema.BoolAttribute{
							Optional:            true,
							Description:         "Whether to cascade privileges to child nodes.",
							MarkdownDescription: "Whether to cascade privileges to child nodes.",
						},
					},
				},
				Description:         "Map of managing nodes with their privileges and cascade options. The map key is the node name/ID.",
				MarkdownDescription: "Map of managing nodes with their privileges and cascade options. The map key is the node name/ID. Each managing node has `privileges` (optional) and `cascade` (optional) fields.",
			},
			"enforcement_policies": schema.MapAttribute{
				Optional: true,
				Validators: []validator.Map{
					enforcementPoliciesMapKeyValidator{},
				},
				ElementType:         types.StringType,
				Description:         "Map of enforcement policies to apply to the role. The map key is the enforcement policy key (e.g., REQUIRE_TWO_FACTOR) and the value is the policy value (e.g., false).",
				MarkdownDescription: "Map of enforcement policies to apply to the role. The map key is the enforcement policy key (e.g., `REQUIRE_TWO_FACTOR`) and the value is the policy value (e.g., `\"false\"`).",
			},
			"managed_company": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					managedCompanyValidator{},
				},
			},
		},
	}
}
