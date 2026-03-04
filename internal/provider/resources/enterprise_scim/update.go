// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update allows only prefix and unique_groups to be changed.
func (r *EnterpriseScimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state EnterpriseScimResourceModel

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

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	// managed_company cannot be changed on update (schema has RequiresReplace; this is a safety check).
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			ErrSummaryManagedCompanyCannotBeUpdated,
			ErrDetailManagedCompany,
		)
		return
	}

	// Use state's managed_company for API context (plan equals state here; state is source of truth for identity).
	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, plan.ManagedCompany, func() error {
		command := buildUpdateCommand(state.Id.ValueString(), &plan)
		_, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpUpdateScim)
		return err
	}, ErrSummaryUpdateFailed, &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve id, node, managed_company from state
	plan.Id = state.Id
	plan.Node = state.Node
	plan.ScimURL = state.ScimURL

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
