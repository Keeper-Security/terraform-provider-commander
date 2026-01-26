package enterpriserole

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
		// Step 1: Fetch and process users/teams before creating the role
		// For create, stateUsers and stateTeams are null/empty, only planUsers and planTeams have items to add
		users, err := utils.FetchAndProcessUsers(ctx, r.apiManager, types.SetNull(types.StringType), data.Users)
		if err != nil {
			return err
		}

		teams, err := utils.FetchAndProcessTeams(ctx, r.apiManager, types.SetNull(types.StringType), data.Teams)
		if err != nil {
			return err
		}

		// Step 2: Build and execute the command to create the role
		command := buildEnterpriseRoleAddCommand(data)

		_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise role")
		if err != nil {
			return fmt.Errorf("create enterprise role failed: %w", err)
		}

		// Set the role ID (using name as ID)
		// TODO: Get the role ID from the response when Commander CLI command is updated
		roleId := data.Name.ValueString()
		data.Id = types.StringValue(roleId)

		// Step 3: Validate managing nodes before processing
		// Fetch all available nodes for current scope
		currentScopeNodes, err := r.apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json", "Unable to fetch enterprise nodes for the managed company")
		if err != nil {
			return err
		}

		// Validate that all managing nodes exist in the available nodes
		if err := validateManagingNodes(ctx, data.ManagingNodes, currentScopeNodes.Data); err != nil {
			return fmt.Errorf("managing nodes validation failed: %w", err)
		}

		// Step 4: Process managing nodes if provided
		if err := processManagingNodes(ctx, r.apiManager, roleId, data.ManagingNodes); err != nil {
			return err
		}

		// Step 5: Process enforcement policies if provided
		if err := processEnforcementPolicies(ctx, r.apiManager, roleId, data.EnforcementPolicies); err != nil {
			return err
		}

		// Step 6: Add users and teams to the recently created role
		var userTeamFlags []string
		if users != "" {
			userTeamFlags = append(userTeamFlags, users)
		}
		if teams != "" {
			userTeamFlags = append(userTeamFlags, teams)
		}

		if len(userTeamFlags) > 0 {
			command = fmt.Sprintf("enterprise-role '%s' -f %s", roleId, strings.Join(userTeamFlags, " "))
			_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to add users/teams to the enterprise role")
			if err != nil {
				return err
			}
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
