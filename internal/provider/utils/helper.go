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
