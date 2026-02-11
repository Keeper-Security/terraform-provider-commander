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

// ParseUsersResponse parses the JSON response from enterprise-info -u command.
func ParseUsersResponse(data interface{}) ([]EnterpriseUserResponse, error) {
	var users []EnterpriseUserResponse

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &users); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise users list from Service Mode API response: %w", err)
	}

	return users, nil
}

// BuildUserLookupMaps creates lookup maps from API response.
func BuildUserLookupMaps(usersRespData []EnterpriseUserResponse) LookupMaps {
	identifierToId := make(map[string]string)
	idToIdentifier := make(map[string]string)

	for _, user := range usersRespData {
		if user.UserId > 0 && user.Email != "" {
			userIdStr := strconv.Itoa(user.UserId)
			identifierToId[user.Email] = userIdStr
			idToIdentifier[userIdStr] = user.Email
		}
	}

	return LookupMaps{
		IdentifierToId: identifierToId,
		IdToIdentifier: idToIdentifier,
	}
}

// ConvertUsersToIdMap converts a types.Set of users to a map of user_id -> original input.
func ConvertUsersToIdMap(users types.Set, lookup LookupMaps, usersRespData []EnterpriseUserResponse) (map[string]string, error) {
	validateUser := func(userInput string) (bool, string) {
		for _, user := range usersRespData {
			if user.Email == userInput && user.UserId <= 0 {
				return false, "user '" + userInput + "' exists but has no valid user_id. This user cannot be used"
			}
		}
		return true, ""
	}

	return ConvertItemsToIdMap(users, lookup, "user", validateUser)
}

// FetchAndProcessUsers processes users for both create and update operations
// For create: stateUsers should be null/empty, planUsers contains users to add
// For update: compares stateUsers (old) with planUsers (new) to determine additions and removals
// Returns a string with -au "user_id" for additions and -ru "user_id" for removals.
func FetchAndProcessUsers(ctx context.Context, apiManager *api.ApiManager, stateUsers types.Set, planUsers types.Set, optionalFlags ...string) (string, error) {
	// Early return if both are empty/null
	if (stateUsers.IsNull() || len(stateUsers.Elements()) == 0) &&
		(planUsers.IsNull() || len(planUsers.Elements()) == 0) {
		return "", nil
	}

	addUserFlagName := "-au"
	removeUserFlagName := "-ru"
	if len(optionalFlags) >= 1 && optionalFlags[0] != "" {
		addUserFlagName = optionalFlags[0]
	}
	if len(optionalFlags) >= 2 && optionalFlags[1] != "" {
		removeUserFlagName = optionalFlags[1]
	}

	// Fetch users from API
	usersResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -u --format json", "Unable to fetch enterprise users")
	if err != nil {
		return "", err
	}

	// Parse the users response
	usersRespData, err := ParseUsersResponse(usersResp.Data)
	if err != nil {
		return "", err
	}

	// Build lookup maps
	lookup := BuildUserLookupMaps(usersRespData)

	// Create a map of user_id -> EnterpriseUserResponse for status checking
	userIdToUserInfo := make(map[string]EnterpriseUserResponse)
	for _, user := range usersRespData {
		if user.UserId > 0 {
			userIdStr := strconv.Itoa(user.UserId)
			userIdToUserInfo[userIdStr] = user
		}
	}

	// Convert state users to user_id map (old users)
	stateUserIdMap, err := ConvertUsersToIdMap(stateUsers, lookup, usersRespData)
	if err != nil {
		return "", err
	}

	// Convert plan users to user_id map (new users)
	planUserIdMap, err := ConvertUsersToIdMap(planUsers, lookup, usersRespData)
	if err != nil {
		return "", err
	}

	// Early return if no changes
	if len(stateUserIdMap) == 0 && len(planUserIdMap) == 0 {
		return "", nil
	}

	// Find users to add and remove
	var parts []string

	// Add users that are in plan but not in state
	for userId := range planUserIdMap {
		if _, exists := stateUserIdMap[userId]; !exists {
			// Check if user has "Invited" status
			if userInfo, exists := userIdToUserInfo[userId]; exists {
				if userInfo.Status == "Invited" {
					userIdentifier := planUserIdMap[userId] // Get original input (email or user_id)
					return "", fmt.Errorf("user '%s' has status 'Invited'. Users must accept invitation before being added", userIdentifier)
				}
			}
			parts = append(parts, fmt.Sprintf("%s '%s'", addUserFlagName, userId))
		}
	}

	// Remove users that are in state but not in plan
	for userId := range stateUserIdMap {
		if _, exists := planUserIdMap[userId]; !exists {
			parts = append(parts, fmt.Sprintf("%s '%s'", removeUserFlagName, userId))
		}
	}

	if len(parts) == 0 {
		return "", nil
	}

	return strings.Join(parts, " "), nil
}

// RestoreUserInputFormatForUsers converts user emails from API response back to the format
// that the user originally provided in their Terraform configuration.
//
// This function preserves the original user input format to prevent false diffs in Terraform plans.
// If a user specified users by ID (e.g., "123"), the function will return IDs. If they specified
// by email (e.g., "user@example.com"), it will return emails.
//
// Parameters:
//   - userEmails: User emails returned by the API (from enterprise-info command)
//   - currentState: Current Terraform state containing users (what user originally provided)
//
// Returns:
//   - types.Set: Set of users in the original user input format (emails or user IDs)
//   - error: Error if fetching users or building lookup maps fails
//
// Example:
//
//	User config: users = ["123", "456"]
//	API returns: ["user1@example.com", "user2@example.com"]
//	Function returns: ["123", "456"] (preserves original IDs)
func RestoreUserInputFormatForUsers(ctx context.Context, apiManager *api.ApiManager, userEmails []string, currentState types.Set) (types.Set, error) {
	return RestoreUserInputFormatFromApiResponse(
		ctx,
		apiManager,
		userEmails,
		currentState,
		"user",
		"enterprise-info -u --format json",
		func(data interface{}) (interface{}, error) { return ParseUsersResponse(data) },
		func(data interface{}) LookupMaps {
			list, ok := data.([]EnterpriseUserResponse)
			if !ok {
				return LookupMaps{}
			}
			return BuildUserLookupMaps(list)
		},
	)
}
