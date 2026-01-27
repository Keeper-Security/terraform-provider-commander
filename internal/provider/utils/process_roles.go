// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleInfo represents a role from the API response
type RoleInfo struct {
	RoleId int    `json:"role_id"`
	Name   string `json:"name"`
}

// ParseRolesResponse parses the JSON response from enterprise-info -r command
func ParseRolesResponse(data interface{}) ([]RoleInfo, error) {
	var roles []RoleInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &roles); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise roles list from Service Mode API response: %w", err)
	}

	return roles, nil
}

// BuildRoleLookupMaps creates lookup maps from API response
func BuildRoleLookupMaps(rolesRespData []RoleInfo) LookupMaps {
	identifierToId := make(map[string]string)
	idToIdentifier := make(map[string]string)

	for _, role := range rolesRespData {
		if role.RoleId > 0 && role.Name != "" {
			roleIdStr := strconv.Itoa(role.RoleId)
			identifierToId[role.Name] = roleIdStr
			idToIdentifier[roleIdStr] = role.Name
		}
	}

	return LookupMaps{
		IdentifierToId: identifierToId,
		IdToIdentifier: idToIdentifier,
	}
}

// ConvertRolesToIdMap converts a types.Set of roles to a map of role_id -> original input
func ConvertRolesToIdMap(roles types.Set, lookup LookupMaps, rolesRespData []RoleInfo) (map[string]string, error) {
	validateRole := func(userInput string) (bool, string) {
		for _, role := range rolesRespData {
			if role.Name == userInput && role.RoleId <= 0 {
				return false, "role '" + userInput + "' exists but has no valid role_id. This role cannot be used"
			}
		}
		return true, ""
	}

	return ConvertItemsToIdMap(roles, lookup, "role", validateRole)
}

// FetchAndProcessRoles processes roles for both create and update operations
// For create: stateRoles should be null/empty, planRoles contains roles to add
// For update: compares stateRoles (old) with planRoles (new) to determine additions and removals
// Returns a string with -ar "role_id" for additions and -rr "role_id" for removals
func FetchAndProcessRoles(ctx context.Context, apiManager *api.ApiManager, stateRoles types.Set, planRoles types.Set) (string, error) {
	// Early return if both are empty/null
	if (stateRoles.IsNull() || len(stateRoles.Elements()) == 0) &&
		(planRoles.IsNull() || len(planRoles.Elements()) == 0) {
		return "", nil
	}

	// Fetch roles from API
	rolesResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -r --format json", "Unable to fetch enterprise roles")
	if err != nil {
		return "", err
	}

	// Parse the roles response
	rolesRespData, err := ParseRolesResponse(rolesResp.Data)
	if err != nil {
		return "", err
	}

	// Build lookup maps
	lookup := BuildRoleLookupMaps(rolesRespData)

	// Convert state roles to role_id map (old roles)
	stateRoleIdMap, err := ConvertRolesToIdMap(stateRoles, lookup, rolesRespData)
	if err != nil {
		return "", err
	}

	// Convert plan roles to role_id map (new roles)
	planRoleIdMap, err := ConvertRolesToIdMap(planRoles, lookup, rolesRespData)
	if err != nil {
		return "", err
	}

	// Early return if no changes
	if len(stateRoleIdMap) == 0 && len(planRoleIdMap) == 0 {
		return "", nil
	}

	// Find roles to add and remove
	var parts []string

	// Add roles that are in plan but not in state
	for roleId := range planRoleIdMap {
		if _, exists := stateRoleIdMap[roleId]; !exists {
			parts = append(parts, fmt.Sprintf("-ar '%s'", roleId))
		}
	}

	// Remove roles that are in state but not in plan
	for roleId := range stateRoleIdMap {
		if _, exists := planRoleIdMap[roleId]; !exists {
			parts = append(parts, fmt.Sprintf("-rr '%s'", roleId))
		}
	}

	if len(parts) == 0 {
		return "", nil
	}

	return strings.Join(parts, " "), nil
}

// RestoreUserInputFormatForRoles converts role names from API response back to the format
// that the user originally provided in their Terraform configuration.
//
// This function preserves the original user input format to prevent false diffs in Terraform plans.
// If a user specified roles by ID (e.g., "123"), the function will return IDs. If they specified
// by name (e.g., "Admin Role"), it will return names.
//
// Parameters:
//   - roleNames: Role names returned by the API (from enterprise-info command)
//   - currentState: Current Terraform state containing roles (what user originally provided)
//
// Returns:
//   - types.Set: Set of roles in the original user input format (names or role IDs)
//   - error: Error if fetching roles or building lookup maps fails
//
// Example:
//   User config: roles = ["123", "456"]
//   API returns: ["Admin Role", "User Role"]
//   Function returns: ["123", "456"] (preserves original IDs)
func RestoreUserInputFormatForRoles(ctx context.Context, apiManager *api.ApiManager, roleNames []string, currentState types.Set) (types.Set, error) {
	return RestoreUserInputFormatFromApiResponse(
		ctx,
		apiManager,
		roleNames,
		currentState,
		"role",
		"enterprise-info -r --format json",
		func(data interface{}) (interface{}, error) { return ParseRolesResponse(data) },
		func(data interface{}) LookupMaps { return BuildRoleLookupMaps(data.([]RoleInfo)) },
	)
}
