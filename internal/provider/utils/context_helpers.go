// Copyright (c) Keeper Security, Inc.
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

// RunWithMspContext switches to MSP context if apiManager.IsMspAccount, then runs op.
// If switch or op fails, adds error to diags with the given summary and returns the error.
// ErrResourceRemoved is not added to diags (caller should return without setting state).
func RunWithMspContext(
	ctx context.Context,
	apiManager *api.ApiManager,
	op func() error,
	errorSummary string,
	diags *diag.Diagnostics,
) error {
	if apiManager.IsMspAccount {
		if err := SwitchToMsp(ctx, apiManager); err != nil {
			diags.AddError(errorSummary, fmt.Sprintf("Failed to switch to MSP: %s", err.Error()))
			return err
		}
	}
	if err := op(); err != nil {
		if !errors.Is(err, ErrResourceRemoved) {
			diags.AddError(errorSummary, err.Error())
		}
		return err
	}
	return nil
}
