// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// extractManagingNodes extracts managing nodes from a Map and returns them as a map with node names as keys.
// This helper function can be used when building commands that require managing nodes.
// Returns: map[nodeName]ManagingNodeModel.
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
// Returns: map[policyKey]types.String.
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
// Format: enterprise-role "Role ID/Name" -aa 'Managing Node Name' --cascade on/off.
func buildAddManagingNodeCommand(roleId string, managingNodeName string, cascade bool) string {
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
func buildAddRemoveManagingNodePrivilegesCommand(roleId string, managingNodeName string, addPrivileges []string, removePrivileges []string) string {
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

// getPrivilegesList extracts privileges from a ManagingNodeModel and returns them as a sorted string slice.
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

// getCascadeValue returns the cascade boolean value from a ManagingNodeModel.
func getCascadeValue(node ManagingNodeModel) bool {
	if node.Cascade.IsNull() || node.Cascade.IsUnknown() {
		return false
	}
	return node.Cascade.ValueBool()
}

// addManagingNodeWithPrivileges adds a managing node to a role and sets its privileges
// managingNodeName is the node name/ID (map key).
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
// 2. Adding privileges to each managing node.
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

// validateManagingNodes validates that all managing nodes exist in the available nodes list.
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

	var availableNodes []utils.EnterpriseNodeResponse
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

// extractRoleIdFromCreateRoleResponse extracts the role id from the response.
func extractRoleIdFromCreateRoleResponse(s string) (string, bool) {
	re := regexp.MustCompile(`Role ID :\s*(\d+)`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}
