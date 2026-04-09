// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RunWithManagedCompanyContext runs op inside ExecuteWithManagedCompanyContext.
// If op returns an error, it is added to diags with the given summary, unless the error is ErrResourceRemoved
// (caller should return without setting state in that case).
// Returns the error from op so the caller can check for ErrResourceRemoved.
func RunWithManagedCompanyContext(
	ctx context.Context,
	apiManager *api.ApiManager,
	managedCompany types.String,
	op func() error,
	errorSummary string,
	diags *diag.Diagnostics,
) error {
	err := ExecuteWithManagedCompanyContext(ctx, apiManager, managedCompany, op)
	if err != nil {
		if !errors.Is(err, ErrResourceRemoved) {
			diags.AddError(errorSummary, err.Error())
		}
		return err
	}
	return nil
}

// RunWithMspContext ensures MSP context (when needed) then runs op.
// Serializes with other context operations and skips switch-to-msp when already in MSP.
func RunWithMspContext(
	ctx context.Context,
	apiManager *api.ApiManager,
	op func() error,
	errorSummary string,
	diags *diag.Diagnostics,
) error {
	apiManager.LockContext()
	defer apiManager.UnlockContext()

	if apiManager.IsMspAccount && apiManager.GetCurrentContext() != "" {
		if err := SwitchToMsp(ctx, apiManager); err != nil {
			diags.AddError(errorSummary, fmt.Sprintf("Failed to switch to MSP: %s", err.Error()))
			return err
		}
		apiManager.SetCurrentContext("")
	}
	if err := op(); err != nil {
		if !errors.Is(err, ErrResourceRemoved) {
			diags.AddError(errorSummary, err.Error())
		}
		return err
	}
	return nil
}
