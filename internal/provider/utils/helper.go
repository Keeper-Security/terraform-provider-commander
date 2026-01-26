// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SwitchToManageCompany switches to the specified managed company
func SwitchToManageCompany(ctx context.Context, apiManager *api.ApiManager, manageCompany string) error {
	command := fmt.Sprintf("switch-to-mc '%s'", manageCompany)
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to switch to manage company")
	return err
}

// SwitchToMsp switches back to MSP context
func SwitchToMsp(ctx context.Context, apiManager *api.ApiManager) error {
	command := "switch-to-msp"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to switch to msp")

	// NOTE: For now we are commenting bec commander cli sending 500 status code
	if err != nil && strings.Contains(err.Error(), "Already MSP") {
		return nil
	}
	return err
}

// Perform msp down
func MspDown(ctx context.Context, apiManager *api.ApiManager) error {
	command := "msp-down"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to perform msp down")
	return err
}

// Perform enterprise down
func EnterpriseDown(ctx context.Context, apiManager *api.ApiManager) error {
	command := "enterprise-down"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to perform enterprise down")
	return err
}

// ExecuteWithManagedCompanyContext executes a function with managed company context switching
// If managedCompany is provided and not null, it switches to MC before execution and back to MSP after
// If managedCompany is not provided, it ensures we're in the correct base context (MSP for MSP accounts, Enterprise for Enterprise)
// This prevents race conditions when Terraform runs resources in parallel or different orders
func ExecuteWithManagedCompanyContext(
	ctx context.Context,
	apiManager *api.ApiManager,
	managedCompany types.String,
	operation func() error,
) (err error) {
	// Track whether we switched to MC so we know to switch back
	switchedToMC := false

	if !managedCompany.IsNull() && !managedCompany.IsUnknown() {
		// Has managed_company - switch to it
		// Sync enterprise data first
		if err := EnterpriseDown(ctx, apiManager); err != nil {
			return fmt.Errorf("Failed to sync enterprise data: %w", err)
		}

		// Switch to managed company
		if err := SwitchToManageCompany(ctx, apiManager, managedCompany.ValueString()); err != nil {
			return fmt.Errorf("Failed to switch to managed company: %w", err)
		}
		switchedToMC = true
	} else {
		// No managed_company provided - ensure we're in the correct base context
		// This prevents race conditions where another resource might have switched to MC
		// For MSP accounts: switch to MSP context
		// For Enterprise accounts: just sync data (already in Enterprise context)
		if apiManager.IsMspAccount {
			// MSP account - explicitly switch to MSP to ensure clean state
			if err := SwitchToMsp(ctx, apiManager); err != nil {
				return fmt.Errorf("Failed to switch to MSP context: %w", err)
			}
		}
		// For both MSP and Enterprise accounts, sync enterprise data
		if err := EnterpriseDown(ctx, apiManager); err != nil {
			return fmt.Errorf("Failed to sync enterprise data: %w", err)
		}
	}

	// Always switch back to MSP after the operation if we switched to MC
	// This ensures clean state for subsequent operations
	defer func() {
		if switchedToMC {
			// Only switch back to MSP if it's an MSP account
			if apiManager.IsMspAccount {
				if switchErr := SwitchToMsp(ctx, apiManager); switchErr != nil {
					if err != nil {
						// Both operation and switch-back failed - combine errors
						err = fmt.Errorf("operation failed: %w; also failed to switch back to MSP: %w", err, switchErr)
					} else {
						// Operation succeeded but switch-back failed - this is critical
						err = fmt.Errorf("Failed to switch back to MSP: %w", switchErr)
					}
				}
			}
		}

	}()

	// Execute the actual operation
	err = operation()
	return err
}

// Note: After creating a node, service mode api returns message like: "Node is created with Node ID: 1169425105420462"
// This function extracts the node id from the response
func ExtractNodeIDFromCreateNodeResponse(s string) (string, bool) {
	re := regexp.MustCompile(`Node ID:\s*(\d+)`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}

// Function to extract the node name from the input string like "Metronlabs\\Aditya Dev Inc" -> "Aditya Dev Inc"
// msp-info retuns node_name as "Metronlabs\\Aditya Dev Inc" if present in child node or node_name as "Metronlabs" if present in root node
func ExtractNodeName(input string) string {
	if idx := strings.LastIndex(input, `\`); idx != -1 {
		return input[idx+1:]
	}
	return input
}

// UnmarshalApiResponse unmarshals API response data into a target struct.
// It handles the common pattern of marshaling interface{} to JSON bytes and then unmarshaling into the target type.
// Parameters:
//   - data: The API response data (typically apiResp.Data from ExecuteCommand)
//   - target: A pointer to the struct/slice that should receive the unmarshaled data
//
// Returns an error if marshaling or unmarshaling fails.
//
// Example usage:
//
//	var roles []RoleInfo
//	if err := utils.UnmarshalApiResponse(apiResp.Data, &roles); err != nil {
//	    return fmt.Errorf("failed to parse roles: %w", err)
//	}
func UnmarshalApiResponse(data interface{}, target interface{}) error {
	// Convert apiResp.Data to JSON bytes
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to process the response from Keeper Commander API: %w", err)
	}

	// Unmarshal JSON bytes into target struct
	if err := json.Unmarshal(dataBytes, target); err != nil {
		return fmt.Errorf("unable to parse API response: %w", err)
	}

	return nil
}

// LookupMaps holds the mappings between identifiers (name/email) and IDs
type LookupMaps struct {
	IdentifierToId map[string]string // identifier (name/email) -> id
	IdToIdentifier map[string]string // id -> identifier (name/email)
}

// ConvertItemsToIdMap is a generic function that converts a types.Set to a map of id -> original input
// It works for roles, users, and teams by accepting lookup maps and validation functions
func ConvertItemsToIdMap(
	items types.Set,
	lookup LookupMaps,
	itemType string, // "role", "user", or "team"
	validateItem func(string) (bool, string), // returns (isValid, errorMessage)
) (map[string]string, error) {
	result := make(map[string]string)

	if items.IsNull() || items.IsUnknown() {
		return result, nil
	}

	elements := items.Elements()
	if len(elements) == 0 {
		return result, nil
	}

	seenIds := make(map[string]string) // id -> original input

	for _, itemElem := range elements {
		itemStr := itemElem.(types.String)
		userInput := itemStr.ValueString()

		if userInput == "" {
			continue
		}

		var itemId string
		var itemIdentifier string

		// Check if input is an id
		if existingIdentifier, isId := lookup.IdToIdentifier[userInput]; isId {
			itemId = userInput
			itemIdentifier = existingIdentifier
		} else if id, isIdentifier := lookup.IdentifierToId[userInput]; isIdentifier {
			// Input is an identifier (name/email), convert to id
			itemId = id
			itemIdentifier = userInput
		} else {
			// Validate if item exists but has invalid id
			isValid, errMsg := validateItem(userInput)
			if !isValid {
				return nil, fmt.Errorf("%s", errMsg)
			}
			return nil, fmt.Errorf("%s '%s' not found. Please provide a valid %s identifier or %s Id", itemType, userInput, itemType, itemType)
		}

		if itemId == "" {
			return nil, fmt.Errorf("%s '%s' resulted in an empty %s_id. This should not happen - please report this issue", itemType, userInput, itemType)
		}

		// Check for duplicates
		if originalInput, exists := seenIds[itemId]; exists {
			return nil, fmt.Errorf("duplicate %s detected: '%s' and '%s' both map to the same %s Id '%s' (%s identifier: '%s')",
				itemType, originalInput, userInput, itemType, itemId, itemType, itemIdentifier)
		}

		seenIds[itemId] = userInput
		result[itemId] = userInput
	}

	return result, nil
}
