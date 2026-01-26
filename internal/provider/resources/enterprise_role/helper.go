// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildEnterpriseRoleAddCommand(data EnterpriseRoleResourceModel) string {
	var parts []string

	parts = append(parts, "enterprise-role")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	return strings.Join(parts, " ")
}

func buildEnterpriseRoleUpdateCommand(plan *EnterpriseRoleResourceModel, state *EnterpriseRoleResourceModel) string {
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

	return strings.Join(parts, " ")
}

// extractManagingNodes extracts managing nodes from a Map and returns them as a map with node names as keys
// This helper function can be used when building commands that require managing nodes
// Returns: map[nodeName]ManagingNodeModel
func extractManagingNodes(ctx context.Context, managingNodesMap types.Map) (map[string]ManagingNodeModel, error) {
	if managingNodesMap.IsNull() || managingNodesMap.IsUnknown() {
		return make(map[string]ManagingNodeModel), nil
	}

	var managingNodesMapValue map[string]ManagingNodeModel
	diags := managingNodesMap.ElementsAs(ctx, &managingNodesMapValue, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract managing nodes: %v", diags)
	}

	return managingNodesMapValue, nil
}

// extractEnforcementPolicies extracts enforcement policies from a Map and returns them as a map with policy keys as keys.
// Map value type is String (policy value).
// Returns: map[policyKey]types.String
func extractEnforcementPolicies(ctx context.Context, enforcementPoliciesMap types.Map) (map[string]types.String, error) {
	if enforcementPoliciesMap.IsNull() || enforcementPoliciesMap.IsUnknown() {
		return make(map[string]types.String), nil
	}

	var enforcementPoliciesMapValue map[string]types.String
	diags := enforcementPoliciesMap.ElementsAs(ctx, &enforcementPoliciesMapValue, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract enforcement policies: %v", diags)
	}

	return enforcementPoliciesMapValue, nil
}

// buildUpdateEnforcementPoliciesCommand builds a command to set enforcement policy values for a role.
// Format: enterprise-role 'Role ID/Name' --enforcement 'KEY:VALUE' --enforcement 'KEY2:VALUE2' ...
//
// Note: This function only supports setting/updating enforcement policies present in the input.
// If a key is removed from Terraform config, there is currently no known Commander CLI "unset" flag,
// so removals are treated as "stop managing" rather than actively unsetting on the remote side.
func buildUpdateEnforcementPoliciesCommand(roleId string, policies map[string]types.String) (string, error) {
	if len(policies) == 0 {
		return "", nil
	}

	keys := make([]string, 0, len(policies))
	for k := range policies {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	parts = append(parts, fmt.Sprintf("enterprise-role '%s'", roleId))

	for _, key := range keys {
		v := policies[key]
		if v.IsNull() || v.IsUnknown() {
			return "", fmt.Errorf("enforcement policy value for key '%s' is null/unknown", key)
		}
		if v.ValueString() == "" {
			return "", fmt.Errorf("enforcement policy value for key '%s' cannot be an empty string", key)
		}
		parts = append(parts, fmt.Sprintf("--enforcement '%s:%s'", key, v.ValueString()))
	}

	return strings.Join(parts, " "), nil
}

// buildAddManagingNodeCommand builds the command to add a managing node to a role
// Format: enterprise-role "Role ID/Name" -aa 'Managing Node Name' --cascade on/off
func buildAddManagingNodeCommand(roleId, managingNodeName string, cascade bool) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("enterprise-role '%s' -aa '%s'", roleId, managingNodeName))

	if cascade {
		parts = append(parts, "--cascade on")
	} else {
		parts = append(parts, "--cascade off")
	}

	return strings.Join(parts, " ")
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
func buildUpdateManagingNodeCascadeCommand(roleId, managingNodeName string, cascade bool) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("enterprise-role '%s' -aa '%s'", roleId, managingNodeName))

	if cascade {
		parts = append(parts, "--cascade on")
	} else {
		parts = append(parts, "--cascade off")
	}

	return strings.Join(parts, " ")
}

// buildAddRemoveManagingNodePrivilegesCommand builds the command to add/remove privileges for a managing node.
// Format: enterprise-role 'Role ID/Name' --node 'Managing Node Name' -ap 'priv1' -ap 'priv2' -rp 'priv3' -rp 'priv4' ...
func buildAddRemoveManagingNodePrivilegesCommand(roleId, managingNodeName string, addPrivileges []string, removePrivileges []string) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("enterprise-role '%s' --node '%s'", roleId, managingNodeName))

	// Add privileges with -ap
	for _, privilege := range addPrivileges {
		parts = append(parts, fmt.Sprintf("-ap '%s'", privilege))
	}

	// Remove privileges with -rp
	for _, privilege := range removePrivileges {
		parts = append(parts, fmt.Sprintf("-rp '%s'", privilege))
	}

	return strings.Join(parts, " ")
}

// Note: getManagingNodeKey is no longer needed since map keys are the node names

// getPrivilegesList extracts privileges from a ManagingNodeModel and returns them as a sorted string slice
func getPrivilegesList(ctx context.Context, node ManagingNodeModel) ([]string, error) {
	if node.Privileges.IsNull() || node.Privileges.IsUnknown() {
		return []string{}, nil
	}

	var privileges []types.String
	diags := node.Privileges.ElementsAs(ctx, &privileges, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract privileges: %v", diags)
	}

	privilegeList := make([]string, 0, len(privileges))
	for _, priv := range privileges {
		if !priv.IsNull() && !priv.IsUnknown() {
			privilegeList = append(privilegeList, priv.ValueString())
		}
	}

	// Make command generation deterministic
	sort.Strings(privilegeList)

	return privilegeList, nil
}

// getCascadeValue returns the cascade boolean value from a ManagingNodeModel
func getCascadeValue(node ManagingNodeModel) bool {
	if node.Cascade.IsNull() || node.Cascade.IsUnknown() {
		return false
	}
	return node.Cascade.ValueBool()
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

// addManagingNodeWithPrivileges adds a managing node to a role and sets its privileges
// managingNodeName is the node name/ID (map key)
func addManagingNodeWithPrivileges(ctx context.Context, apiManager *api.ApiManager, roleId string, managingNodeName string, node ManagingNodeModel) error {
	if managingNodeName == "" {
		return fmt.Errorf("managing node name cannot be empty")
	}

	cascade := getCascadeValue(node)

	// Step 1: Add the managing node to the role
	addNodeCommand := buildAddManagingNodeCommand(roleId, managingNodeName, cascade)
	_, err := apiManager.ExecuteCommand(ctx, addNodeCommand, fmt.Sprintf("Unable to add managing node '%s' to role", managingNodeName))
	if err != nil {
		return fmt.Errorf("failed to add managing node '%s': %w", managingNodeName, err)
	}

	// Step 2: Add privileges to the managing node if provided
	privilegeList, err := getPrivilegesList(ctx, node)
	if err != nil {
		return fmt.Errorf("failed to extract privileges for managing node '%s': %w", managingNodeName, err)
	}

	if len(privilegeList) > 0 {
		addPrivilegesCommand := buildAddRemoveManagingNodePrivilegesCommand(roleId, managingNodeName, privilegeList, nil)
		_, err := apiManager.ExecuteCommand(ctx, addPrivilegesCommand, fmt.Sprintf("Unable to add privileges to managing node '%s'", managingNodeName))
		if err != nil {
			return fmt.Errorf("failed to add privileges to managing node '%s': %w", managingNodeName, err)
		}
	}

	return nil
}

// processManagingNodes processes all managing nodes for a role (for CREATE operation)
// 1. Adding each managing node to the role with cascade option
// 2. Adding privileges to each managing node
func processManagingNodes(ctx context.Context, apiManager *api.ApiManager, roleId string, managingNodesMap types.Map) error {
	if managingNodesMap.IsNull() || managingNodesMap.IsUnknown() {
		return nil
	}

	managingNodes, err := extractManagingNodes(ctx, managingNodesMap)
	if err != nil {
		return fmt.Errorf("failed to extract managing nodes: %w", err)
	}

	// Process each managing node - map key is the node name
	for nodeName, managingNode := range managingNodes {
		if err := addManagingNodeWithPrivileges(ctx, apiManager, roleId, nodeName, managingNode); err != nil {
			return err
		}
	}

	return nil
}

// processEnforcementPolicies sets enforcement policies on a role (for CREATE operation).
// This will set all policies present in the config on the role.
func processEnforcementPolicies(ctx context.Context, apiManager *api.ApiManager, roleId string, enforcementPoliciesMap types.Map) error {
	policies, err := extractEnforcementPolicies(ctx, enforcementPoliciesMap)
	if err != nil {
		return err
	}

	if len(policies) == 0 {
		return nil
	}

	cmd, err := buildUpdateEnforcementPoliciesCommand(roleId, policies)
	if err != nil {
		return err
	}
	if cmd == "" {
		return nil
	}

	_, err = apiManager.ExecuteCommand(ctx, cmd, "Unable to set enforcement policies for role")
	if err != nil {
		return fmt.Errorf("failed to set enforcement policies for role '%s': %w", roleId, err)
	}

	return nil
}

/* NOTE: currently we dont need this logic bec when node name changes terraform will remove old managing node and add new managing node separately with its privileges and cascade option*/
// validateManagingNodeNamesUnchanged validates that no managing node names have been changed
// Users cannot change managing node names in a single update - they must remove the old node and add a new one in separate operations
// With MapNestedAttribute, the map keys are the node names, so we can directly compare keys
// func validateManagingNodeNamesUnchanged(ctx context.Context, planNodesMap, stateNodesMap types.Map) error {
// 	// Get map keys (node names) - Elements() returns map[string]ManagingNodeModel where keys are node names
// 	var planNodeNames map[string]bool
// 	var stateNodeNames map[string]bool

// 	if !planNodesMap.IsNull() && !planNodesMap.IsUnknown() {
// 		planElements := planNodesMap.Elements()
// 		planNodeNames = make(map[string]bool, len(planElements))
// 		for key := range planElements {
// 			planNodeNames[key] = true
// 		}
// 	} else {
// 		planNodeNames = make(map[string]bool)
// 	}

// 	if !stateNodesMap.IsNull() && !stateNodesMap.IsUnknown() {
// 		stateElements := stateNodesMap.Elements()
// 		stateNodeNames = make(map[string]bool, len(stateElements))
// 		for key := range stateElements {
// 			stateNodeNames[key] = true
// 		}
// 	} else {
// 		stateNodeNames = make(map[string]bool)
// 	}

// 	// Find removed and added node names
// 	var removedNames []string
// 	var addedNames []string

// 	for stateName := range stateNodeNames {
// 		if !planNodeNames[stateName] {
// 			removedNames = append(removedNames, stateName)
// 		}
// 	}

// 	for planName := range planNodeNames {
// 		if !stateNodeNames[planName] {
// 			addedNames = append(addedNames, planName)
// 		}
// 	}

// 	// If nodes were both removed and added in the same update, it might be a rename attempt
// 	// We prevent this to force users to do remove and add in separate operations
// 	if len(removedNames) > 0 && len(addedNames) > 0 {
// 		return fmt.Errorf(
// 			"cannot change managing node names. To change a managing node name, you must first remove the old node(s) '%s' in one update, then add the new node(s) '%s' in a separate update. Please remove the old managing node(s) and apply, then add the new managing node(s) in a subsequent update",
// 			strings.Join(removedNames, ", "),
// 			strings.Join(addedNames, ", "),
// 		)
// 	}

// 	return nil
// }

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

// validateManagingNodes validates that all managing nodes exist in the available nodes list
func validateManagingNodes(ctx context.Context, managingNodesMap types.Map, apiResponseData interface{}) error {
	if managingNodesMap.IsNull() || managingNodesMap.IsUnknown() {
		return nil
	}

	// Extract managing nodes from the map - map keys are the node names
	managingNodes, err := extractManagingNodes(ctx, managingNodesMap)
	if err != nil {
		return fmt.Errorf("failed to extract managing nodes: %w", err)
	}

	if len(managingNodes) == 0 {
		return nil
	}

	// Parse the API response to get available nodes
	dataBytes, err := json.Marshal(apiResponseData)
	if err != nil {
		return fmt.Errorf("unable to process the response from Keeper Commander API: %w", err)
	}

	var availableNodes []NodeInfo
	if err := json.Unmarshal(dataBytes, &availableNodes); err != nil {
		return fmt.Errorf("unable to parse nodes list from API response: %w", err)
	}

	// Create a map of available node names for quick lookup
	availableNodeNames := make(map[string]bool)
	for _, node := range availableNodes {
		availableNodeNames[node.Name] = true
	}

	// Validate each managing node - map keys are the node names
	var missingNodes []string
	for nodeName := range managingNodes {
		if !availableNodeNames[nodeName] {
			missingNodes = append(missingNodes, nodeName)
		}
	}

	// If any nodes are missing, return an error
	if len(missingNodes) > 0 {
		return fmt.Errorf("the following managing nodes are not present: %s", strings.Join(missingNodes, ", "))
	}

	return nil
}
