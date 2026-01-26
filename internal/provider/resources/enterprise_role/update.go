package enterpriserole

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

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

		// Update users if they have changed
		// 1. Removed users -> remove via -ru
		// 2. Added users -> add via -au
		if !plan.Users.Equal(state.Users) {
			users, err := utils.FetchAndProcessUsers(ctx, r.apiManager, state.Users, plan.Users)
			if err != nil {
				return fmt.Errorf("failed to process users: %w", err)
			}

			if users != "" {
				command := fmt.Sprintf("enterprise-role '%s' -f %s", state.Id.ValueString(), users)
				_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to update users/teams for the enterprise role")
				if err != nil {
					return fmt.Errorf("failed to update users and teams: %w", err)
				}
			}
		}

		// Update teams if they have changed
		// 1. Removed teams -> remove via -rt
		// 2. Added teams -> add via -at
		if !plan.Teams.Equal(state.Teams) {
			teams, err := utils.FetchAndProcessTeams(ctx, r.apiManager, state.Teams, plan.Teams)
			if err != nil {
				return fmt.Errorf("failed to process teams: %w", err)
			}

			if teams != "" {
				command := fmt.Sprintf("enterprise-role '%s' -f %s", state.Id.ValueString(), teams)
				_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to update teams for the enterprise role")
				if err != nil {
					return fmt.Errorf("failed to update teams: %w", err)
				}
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
