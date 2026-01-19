package utils

import (
	"context"
	"fmt"

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
func ExecuteWithManagedCompanyContext(
	ctx context.Context,
	apiManager *api.ApiManager,
	managedCompany types.String,
	operation func() error,
) (err error) {
	// If managed company is provided, switch to it before the operation
	err = EnterpriseDown(ctx, apiManager)

	if !managedCompany.IsNull() && !managedCompany.IsUnknown() {
		if err := SwitchToManageCompany(ctx, apiManager, managedCompany.ValueString()); err != nil {
			return fmt.Errorf("Failed to switch to managed company: %w", err)
		}

		// Always switch back to MSP after the operation (even if it fails)
		defer func() {
			if switchErr := SwitchToMsp(ctx, apiManager); switchErr != nil {
				if err != nil {
					// Both operation and switch-back failed - combine errors so user knows about both
					err = fmt.Errorf("operation failed: %w; also failed to switch back to MSP: %w", err, switchErr)
				} else {
					// Operation succeeded but switch-back failed - this is critical, user must know
					err = fmt.Errorf("Failed to switch back to MSP: %w", switchErr)
				}
			}
		}()
	}

	// Execute the actual operation
	err = operation()
	return err
}
