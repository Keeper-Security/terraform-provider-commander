// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"Invalid Keeper Enterprise Role Configuration",
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
		roleId := state.Id.ValueString()

		// Step 1: Update basic role attributes (name, node, etc.)
		if err := updateRoleBasicAttributes(ctx, r.apiManager, &plan, &state); err != nil {
			return err
		}

		// Step 2: Update users if changed
		// 1. Removed users -> remove via -ru
		// 2. Added users -> add via -au
		if err := updateRoleUsers(ctx, r.apiManager, roleId, &plan, &state); err != nil {
			return err
		}

		// Step 3: Update teams if changed
		// 1. Removed teams -> remove via -rt
		// 2. Added teams -> add via -at
		if err := updateRoleTeams(ctx, r.apiManager, roleId, &plan, &state); err != nil {
			return err
		}

		// Step 4: Update managing nodes if changed
		// 1. Removed nodes -> remove via -ra
		// 2. Added nodes -> add via -aa with privileges and cascade
		// 3. Changed cascade -> update via -aa with --cascade
		// 4. Changed privileges -> update via --node with -ap flags
		if !plan.ManagingNodes.Equal(state.ManagingNodes) {
			if err := updateRoleManagingNodes(ctx, r.apiManager, roleId, &plan, &state); err != nil {
				return err
			}
		}

		// Step 5: Update enforcement policies if changed
		if !plan.EnforcementPolicies.Equal(state.EnforcementPolicies) {
			if err := processEnforcementPoliciesUpdate(ctx, r.apiManager, roleId, plan.EnforcementPolicies, state.EnforcementPolicies); err != nil {
				return fmt.Errorf("failed to update enforcement policies: %w", err)
			}
		}

		// Keep the same ID
		plan.Id = state.Id

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise Role Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// validateTeamsAndManagingNodesMutualExclusivity validates that only one of teams or managing_nodes is provided
func validateTeamsAndManagingNodesMutualExclusivity(teams types.Set, managingNodes types.Map) error {
	// Check if teams is provided (not null, not unknown, and has elements)
	hasTeams := !teams.IsNull() && !teams.IsUnknown()
	if hasTeams {
		elements := teams.Elements()
		hasTeams = len(elements) > 0
	}

	// Check if managing_nodes is provided (not null, not unknown, and has elements)
	hasManagingNodes := !managingNodes.IsNull() && !managingNodes.IsUnknown()
	if hasManagingNodes {
		elements := managingNodes.Elements()
		hasManagingNodes = len(elements) > 0
	}

	if hasTeams && hasManagingNodes {
		return fmt.Errorf("cannot specify both 'teams' and 'managing_nodes'. Please provide only one of them")
	}

	return nil
}

// updateRoleBasicAttributes updates basic role attributes (name, node, etc.)
func updateRoleBasicAttributes(ctx context.Context, apiManager *api.ApiManager, plan, state *EnterpriseRoleResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-role")

	// Check for changes and build update command
	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	// Role to update
	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise role")
	if err != nil {
		return fmt.Errorf("failed to update enterprise role: %w", err)
	}
	return nil
}

// updateRoleManagingNodes handles updating managing nodes for a role
// 1. Removed nodes -> remove via -ra
// 2. Added nodes -> add via -aa with privileges and cascade
// 3. Changed cascade -> update via -aa with --cascade
// 4. Changed privileges -> update via --node with -ap flags
func updateRoleManagingNodes(ctx context.Context, apiManager *api.ApiManager, roleId string, plan, state *EnterpriseRoleResourceModel) error {
	/* NOTE: currently we dont need this logic bec when node name changes terraform will remove old managing node and add new managing node separately with its privileges and cascade option*/
	// TODO: Validate that no managing node names have been changed
	// if err := validateManagingNodeNamesUnchanged(ctx, plan.ManagingNodes, state.ManagingNodes); err != nil {
	// 	return fmt.Errorf("managing nodes update validation failed: %w", err)
	// }

	// Validate new managing nodes before processing by fetching all available nodes for current scope and validating them
	currentScopeNodes, err := apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json", "Unable to fetch enterprise nodes for validation")
	if err != nil {
		return err
	}

	// Validate that all new managing nodes exist in the available nodes
	if err := validateManagingNodes(ctx, plan.ManagingNodes, currentScopeNodes.Data); err != nil {
		return fmt.Errorf("managing nodes validation failed: %w", err)
	}

	// Process managing nodes update
	if err := processManagingNodesUpdate(ctx, apiManager, roleId, plan.ManagingNodes, state.ManagingNodes); err != nil {
		return fmt.Errorf("failed to update managing nodes: %w", err)
	}

	return nil
}

// updateRoleUsers updates users for a role
// 1. Removed users -> remove via -ru
// 2. Added users -> add via -au
func updateRoleUsers(ctx context.Context, apiManager *api.ApiManager, roleId string, plan, state *EnterpriseRoleResourceModel) error {
	if plan.Users.Equal(state.Users) {
		return nil // No change
	}

	users, err := utils.FetchAndProcessUsers(ctx, apiManager, state.Users, plan.Users)
	if err != nil {
		return fmt.Errorf("failed to process users: %w", err)
	}

	if users != "" {
		command := fmt.Sprintf("enterprise-role '%s' -f %s", roleId, users)
		_, err = apiManager.ExecuteCommand(ctx, command, "Unable to update users for the enterprise role")
		if err != nil {
			return fmt.Errorf("failed to update users: %w", err)
		}
	}

	return nil
}

// updateRoleTeams updates teams for a role
// 1. Removed teams -> remove via -rt
// 2. Added teams -> add via -at
func updateRoleTeams(ctx context.Context, apiManager *api.ApiManager, roleId string, plan, state *EnterpriseRoleResourceModel) error {
	if plan.Teams.Equal(state.Teams) {
		return nil // No change
	}

	teams, err := utils.FetchAndProcessTeams(ctx, apiManager, state.Teams, plan.Teams)
	if err != nil {
		return fmt.Errorf("failed to process teams: %w", err)
	}

	if teams != "" {
		command := fmt.Sprintf("enterprise-role '%s' -f %s", roleId, teams)
		_, err = apiManager.ExecuteCommand(ctx, command, "Unable to update teams for the enterprise role")
		if err != nil {
			return fmt.Errorf("failed to update teams: %w", err)
		}
	}

	return nil
}

// buildRemoveManagingNodeCommand builds the command to remove multiple managing nodes from a role in a single call
// Format: enterprise-role "Role ID/Name" -ra "Node1" -ra "Node2" -ra "Node3" ...
func buildRemoveManagingNodeCommand(roleId string, managingNodeNames []string) string {
	if len(managingNodeNames) == 0 {
		return ""
	}

	var parts []string
	parts = append(parts, fmt.Sprintf("enterprise-role '%s'", roleId))

	// Add each -ra flag for each node to remove
	for _, nodeName := range managingNodeNames {
		parts = append(parts, fmt.Sprintf("-ra '%s'", nodeName))
	}

	return strings.Join(parts, " ")
}

// buildUpdateManagingNodeCascadeCommand builds the command to update cascade option for a managing node
// Format: enterprise-role "Role ID/Name" -aa "Managing Node Name" --cascade on/off
func buildUpdateManagingNodeCascadeCommand(roleId string, managingNodeName string, cascade bool) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("enterprise-role '%s' -aa '%s'", roleId, managingNodeName))

	if cascade {
		parts = append(parts, "--cascade on")
	} else {
		parts = append(parts, "--cascade off")
	}

	return strings.Join(parts, " ")
}

// privilegesEqual compares two privilege lists to see if they're equal
func privilegesEqual(privs1, privs2 []string) bool {
	if len(privs1) != len(privs2) {
		return false
	}

	// Create maps for comparison
	map1 := make(map[string]bool)
	map2 := make(map[string]bool)

	for _, p := range privs1 {
		map1[p] = true
	}
	for _, p := range privs2 {
		map2[p] = true
	}

	if len(map1) != len(map2) {
		return false
	}

	for k := range map1 {
		if !map2[k] {
			return false
		}
	}

	return true
}

// processManagingNodesUpdate processes managing nodes updates for a role (for UPDATE operation)
// Handles all cases:
// 1. Removed nodes -> remove via -ra
// 2. Added nodes -> add via -aa with privileges and cascade
// 3. Changed cascade -> update via -aa with --cascade
// 4. Changed privileges -> update via --node with -ap flags
func processManagingNodesUpdate(ctx context.Context, apiManager *api.ApiManager, roleId string, planNodesMap, stateNodesMap types.Map) error {
	// Extract managing nodes from plan and state maps
	planNodesMapValue, err := extractManagingNodes(ctx, planNodesMap)
	if err != nil {
		return fmt.Errorf("failed to extract plan managing nodes: %w", err)
	}

	stateNodesMapValue, err := extractManagingNodes(ctx, stateNodesMap)
	if err != nil {
		return fmt.Errorf("failed to extract state managing nodes: %w", err)
	}

	// Case 1: Find nodes that were removed (in state but not in plan)
	// Collect all nodes to be removed and remove them in a single API call
	var nodesToRemove []string
	for stateKey := range stateNodesMapValue {
		if _, exists := planNodesMapValue[stateKey]; !exists {
			nodesToRemove = append(nodesToRemove, stateKey)
		}
	}

	// Remove all nodes in a single command if any need to be removed
	if len(nodesToRemove) > 0 {
		removeCommand := buildRemoveManagingNodeCommand(roleId, nodesToRemove)
		_, err := apiManager.ExecuteCommand(ctx, removeCommand, fmt.Sprintf("Unable to remove managing nodes from role: %s", strings.Join(nodesToRemove, ", ")))
		if err != nil {
			return fmt.Errorf("failed to remove managing nodes: %w", err)
		}
	}

	// Case 2, 3, 4: Process nodes in plan
	for planKey, planNode := range planNodesMapValue {
		stateNode, exists := stateNodesMapValue[planKey]

		if !exists {
			// Case 2: New node added - add it with privileges and cascade
			if err := addManagingNodeWithPrivileges(ctx, apiManager, roleId, planKey, planNode); err != nil {
				return err
			}
		} else {
			// Node exists in both - check for changes
			planCascade := getCascadeValue(planNode)
			stateCascade := getCascadeValue(stateNode)

			planPrivileges, err := getPrivilegesList(ctx, planNode)
			if err != nil {
				return fmt.Errorf("failed to extract plan privileges for node '%s': %w", planKey, err)
			}

			statePrivileges, err := getPrivilegesList(ctx, stateNode)
			if err != nil {
				return fmt.Errorf("failed to extract state privileges for node '%s': %w", planKey, err)
			}

			// Case 3: Cascade changed - update cascade
			if planCascade != stateCascade {
				updateCascadeCommand := buildUpdateManagingNodeCascadeCommand(roleId, planKey, planCascade)
				_, err := apiManager.ExecuteCommand(ctx, updateCascadeCommand, fmt.Sprintf("Unable to update cascade for managing node '%s'", planKey))
				if err != nil {
					return fmt.Errorf("failed to update cascade for managing node '%s': %w", planKey, err)
				}
			}

			// Case 4: Privileges changed - update privileges (add + remove in one command)
			if !privilegesEqual(planPrivileges, statePrivileges) {
				planSet := make(map[string]struct{}, len(planPrivileges))
				for _, p := range planPrivileges {
					planSet[p] = struct{}{}
				}

				stateSet := make(map[string]struct{}, len(statePrivileges))
				for _, p := range statePrivileges {
					stateSet[p] = struct{}{}
				}

				var toAdd []string
				for p := range planSet {
					if _, exists := stateSet[p]; !exists {
						toAdd = append(toAdd, p)
					}
				}

				var toRemove []string
				for p := range stateSet {
					if _, exists := planSet[p]; !exists {
						toRemove = append(toRemove, p)
					}
				}

				sort.Strings(toAdd)
				sort.Strings(toRemove)

				if len(toAdd) > 0 || len(toRemove) > 0 {
					updatePrivilegesCommand := buildAddRemoveManagingNodePrivilegesCommand(roleId, planKey, toAdd, toRemove)
					_, err := apiManager.ExecuteCommand(ctx, updatePrivilegesCommand, fmt.Sprintf("Unable to update privileges for managing node '%s'", planKey))
					if err != nil {
						return fmt.Errorf("failed to update privileges for managing node '%s': %w", planKey, err)
					}
				}
			}
		}
	}

	return nil
}

// processEnforcementPoliciesUpdate sets/updates enforcement policy values for keys present in the plan (for UPDATE operation).
// It sends a single command containing all changed/added enforcement policies.
func processEnforcementPoliciesUpdate(ctx context.Context, apiManager *api.ApiManager, roleId string, planPoliciesMap, statePoliciesMap types.Map) error {
	planPolicies, err := extractEnforcementPolicies(ctx, planPoliciesMap)
	if err != nil {
		return err
	}

	statePolicies, err := extractEnforcementPolicies(ctx, statePoliciesMap)
	if err != nil {
		return err
	}

	// Only send the policies that changed (or are newly added).
	toSet := make(map[string]types.String)
	for key, planValue := range planPolicies {
		stateValue, exists := statePolicies[key]
		if !exists || !stateValue.Equal(planValue) {
			toSet[key] = planValue
		}
	}

	cmd, err := buildUpdateEnforcementPoliciesCommand(roleId, toSet)
	if err != nil {
		return err
	}
	if cmd == "" {
		return nil
	}

	_, err = apiManager.ExecuteCommand(ctx, cmd, "Unable to update enforcement policies for role")
	if err != nil {
		return fmt.Errorf("failed to update enforcement policies for role '%s': %w", roleId, err)
	}

	return nil
}
