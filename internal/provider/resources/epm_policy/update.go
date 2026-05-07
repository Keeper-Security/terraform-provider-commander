// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EpmPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state EpmPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	utils.RestrictAttributeUpdate(&resp.Diagnostics, []utils.ImmutableAttributeCheck{
		{
			PlanValue:  plan.ManagedCompany,
			StateValue: state.ManagedCompany,
			Summary:    utils.ErrSummaryManagedCompanyCannotBeUpdated,
			Detail:     utils.ErrDetailManagedCompany,
		},
	})
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.IsUnknown() || state.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			ErrSummaryUpdateFailed,
			"Cannot update EPM policy: state has no policy id. Recreate the resource or import an existing policy.",
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, plan.ManagedCompany, func() error {
		if err := utils.EpmSyncDown(ctx, r.ApiManager); err != nil {
			return err
		}

		command := buildUpdateCommand(state.Id.ValueString(), &plan)
		_, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpUpdateEpmPolicy)

		if err != nil {
			return err
		}

		plan.Id = state.Id
		return nil
	}, ErrSummaryUpdateFailed, &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
