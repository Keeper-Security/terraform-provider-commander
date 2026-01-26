// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EnterpriseRoleResource{}
var _ resource.ResourceWithConfigure = &EnterpriseRoleResource{}

type EnterpriseRoleResource struct {
	apiManager *api.ApiManager
}

// ManagingNodeModel represents a single managing node with its privileges and cascade option
// Note: The node name/ID is the map key, so it's not stored in this struct
type ManagingNodeModel struct {
	Privileges types.Set  `tfsdk:"privileges"`
	Cascade    types.Bool `tfsdk:"cascade"`
}

type EnterpriseRoleResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Node                types.String `tfsdk:"node"`
	Users               types.Set    `tfsdk:"users"`
	Teams               types.Set    `tfsdk:"teams"`
	ManagingNodes       types.Map    `tfsdk:"managing_nodes"`
	EnforcementPolicies types.Map    `tfsdk:"enforcement_policies"`
	ManagedCompany      types.String `tfsdk:"managed_company"`
}

func (r *EnterpriseRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_role"
}

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
				Optional: true,
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

func (r *EnterpriseRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiManager, ok := req.ProviderData.(*api.ApiManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			fmt.Sprintf("The provider was not configured correctly. Expected API manager, but got: %T. Please check your provider configuration.", req.ProviderData),
		)
		return
	}

	r.apiManager = apiManager
}

// ensureApiManager validates that apiManager is configured and returns an error if not
func (r *EnterpriseRoleResource) ensureApiManager() error {
	if r.apiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}

func (r *EnterpriseRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseRoleResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that only one of teams or managing_nodes is provided
	if err := validateTeamsAndManagingNodesMutualExclusivity(data.Teams, data.ManagingNodes); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			err.Error(),
		)
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, data.ManagedCompany, func() error {
		// Step 1: Build and execute the command to create the role
		command := buildEnterpriseRoleAddCommand(data)

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise role")
		if err != nil {
			return fmt.Errorf("create enterprise role failed: %w", err)
		}

		// Set the role ID (using name as ID)

		// TODO: Get the role ID from the response when Commander CLI command is updated
		roleId := data.Name.ValueString()
		data.Id = types.StringValue(roleId)

		// Step 2: Validate managing nodes before processing
		// Fetch all available nodes for current scope
		currentScopeNodes, err := r.apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json", "Unable to fetch enterprise nodes for the managed company")
		if err != nil {
			return err
		}

		// Validate that all managing nodes exist in the available nodes
		if err := validateManagingNodes(ctx, data.ManagingNodes, currentScopeNodes.Data); err != nil {
			return fmt.Errorf("managing nodes validation failed: %w", err)
		}

		// Step 3: Process managing nodes if provided
		if err := processManagingNodes(ctx, r.apiManager, roleId, data.ManagingNodes); err != nil {
			return err
		}

		// Step 4: Process enforcement policies if provided
		if err := processEnforcementPolicies(ctx, r.apiManager, roleId, data.EnforcementPolicies); err != nil {
			return err
		}

		// Set the ID in the state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Create Enterprise Role Failed",
			err.Error(),
		)
		return
	}
}

func (r *EnterpriseRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseRoleResourceModel

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {
		// Build the Commander command string
		command := fmt.Sprintf("enterprise-info '%s' -r --format json --columns='visible_below,default_role,admin,node,user_count,users,team_count,teams' -q", state.Id.ValueString())

		apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to read enterprise role")
		if err != nil {
			return fmt.Errorf("Read Enterprise Role Failed: %w", err)
		}

		// Parse the JSON response - it's an array of role objects
		var roles []RoleResponse

		// Unmarshal API response into roles struct
		if err := utils.UnmarshalApiResponse(apiResp.Data, &roles); err != nil {
			return fmt.Errorf("unable to parse enterprise roles list from API response: %w", err)
		}

		// Find the role matching state.Id
		var roleInfo *RoleResponse
		stateId := state.Id.ValueString()
		for i := range roles {
			if strconv.Itoa(roles[i].RoleId) == stateId || roles[i].Name == stateId {
				roleInfo = &roles[i]
				break
			}
		}

		if roleInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return nil
		}

		// Map the response to the state
		state.Id = types.StringValue(stateId)
		state.Name = types.StringValue(roleInfo.Name)
		state.Node = types.StringValue(utils.ExtractNodeName(roleInfo.Node))

		// TODO: Later we will add users, teams, enforcement policies, managing nodes to state

		// Set the updated state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Role Failed",
			err.Error(),
		)
		return
	}
}

func (r *EnterpriseRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseRoleResourceModel
	var state EnterpriseRoleResourceModel

	// Get planned data (new values)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that only one of teams or managing_nodes is provided
	if err := validateTeamsAndManagingNodesMutualExclusivity(plan.Teams, plan.ManagingNodes); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			err.Error(),
		)
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// NOTE: We should not allow user to update managed company, bec. once role is created in managed company, if we allow user to update managed company then switching to that MC we will not able to find that role, so command will fail.
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			"Managed Company Cannot Be Updated",
			"Cannot update the managed_company field. Once an enterprise role is created in a managed company, the managed company cannot be changed.",
		)
		return
	}

	// Use managed company from plan (or state if plan doesn't have it)
	managedCompany := plan.ManagedCompany
	if managedCompany.IsNull() || managedCompany.IsUnknown() {
		managedCompany = state.ManagedCompany
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, managedCompany, func() error {
		command := buildEnterpriseRoleUpdateCommand(&plan, &state)

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise role")
		if err != nil {
			return fmt.Errorf("update enterprise role failed: %w", err)
		}

		// Update managing nodes if they have changed
		// 1. Removed nodes -> remove via -ra
		// 2. Added nodes -> add via -aa with privileges and cascade
		// 3. Changed cascade -> update via -aa with --cascade
		// 4. Changed privileges -> update via --node with -ap flags
		// Note: Changing managing node names is not allowed - users must remove old and add new separately
		if !plan.ManagingNodes.Equal(state.ManagingNodes) {
			/* NOTE: currently we dont need this logic bec when node name changes terraform will remove old managing node and add new managing node separately with its privileges and cascade option*/
			// Validate that no managing node names have been changed
			// if err := validateManagingNodeNamesUnchanged(ctx, plan.ManagingNodes, state.ManagingNodes); err != nil {
			// 	return fmt.Errorf("managing nodes update validation failed: %w", err)
			// }

			// Validate new managing nodes before processing by fetching all available nodes for current scope and validating them
			currentScopeNodes, err := r.apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json", "Unable to fetch enterprise nodes for validation")
			if err != nil {
				return err
			}

			// Validate that all new managing nodes exist in the available nodes
			if err := validateManagingNodes(ctx, plan.ManagingNodes, currentScopeNodes.Data); err != nil {
				return fmt.Errorf("managing nodes validation failed: %w", err)
			}

			// Process managing nodes update
			if err := processManagingNodesUpdate(ctx, r.apiManager, state.Id.ValueString(), plan.ManagingNodes, state.ManagingNodes); err != nil {
				return fmt.Errorf("failed to update managing nodes: %w", err)
			}
		}

		// Update enforcement policies if they have changed
		if !plan.EnforcementPolicies.Equal(state.EnforcementPolicies) {
			if err := processEnforcementPoliciesUpdate(ctx, r.apiManager, state.Id.ValueString(), plan.EnforcementPolicies, state.EnforcementPolicies); err != nil {
				return fmt.Errorf("failed to update enforcement policies: %w", err)
			}
		}

		// Keep the same ID
		plan.Id = state.Id
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise Role Failed",
			err.Error(),
		)
		return
	}
}

func (r *EnterpriseRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnterpriseRoleResourceModel

	// Get state from Terraform
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {

		command := fmt.Sprintf("enterprise-role --delete '%s'", state.Id.ValueString())

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to delete enterprise role")
		if err != nil {
			return fmt.Errorf("Delete Enterprise Role Failed: %w", err)
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Enterprise Role Failed",
			err.Error(),
		)
		return
	}
}

func NewEnterpriseRoleResource() resource.Resource {
	return &EnterpriseRoleResource{}
}
