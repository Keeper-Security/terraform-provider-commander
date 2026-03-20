// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseNodeResourceModel
	var state EnterpriseNodeResourceModel

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

	// Validate ApiManager is configured
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// NOTE: We should not allow user to update managed company, bec. once node is created in managed company, if we allow user to update managed company then switching to that MC we will not able to find that node, so command will fail.
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			"Managed Company Cannot Be Updated",
			"Cannot update the managed_company field. Once an enterprise node is created in a managed company, the managed company cannot be changed.",
		)
		return
	}

	// Use managed company from plan (or state if plan doesn't have it)
	managedCompany := plan.ManagedCompany
	if managedCompany.IsNull() || managedCompany.IsUnknown() {
		managedCompany = state.ManagedCompany
	}

	// Execute with managed company context if provided
	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, managedCompany, func() error {
		if err := updateEnterpriseNode(ctx, r.ApiManager, &plan, &state); err != nil {
			return err
		}
		plan.Id = state.Id
		return nil
	}, "Update Enterprise Node Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func updateEnterpriseNode(ctx context.Context, apiManager *api.ApiManager, plan *EnterpriseNodeResourceModel, state *EnterpriseNodeResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-node")

	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	if !state.Parent.Equal(plan.Parent) {
		parts = append(parts, fmt.Sprintf("--parent '%s'", plan.Parent.ValueString()))
	}

	// --toggle-isolated toggles current state; it does not accept true/false.
	// Only append when the effective "on" state changed. Treat null and false as "off" so
	// false→null does not send the flag (would incorrectly turn isolated on).
	stateIsolatedOn := !state.ToggleIsolated.IsNull() && state.ToggleIsolated.ValueBool()
	planIsolatedOn := !plan.ToggleIsolated.IsNull() && plan.ToggleIsolated.ValueBool()
	if stateIsolatedOn != planIsolatedOn {
		parts = append(parts, "--toggle-isolated")
	}

	// Node to update
	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise node")
	if err != nil {
		return err
	}

	return nil
}
