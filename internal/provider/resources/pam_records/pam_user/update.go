// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PamUserResourceModel
	var state PamUserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	plan.Id = state.Id
	if plan.Login.IsUnknown() || plan.Login.IsNull() {
		plan.Login = state.Login
	}
	if plan.Password.IsUnknown() || plan.Password.IsNull() {
		plan.Password = state.Password
	}
	recordUID := strings.TrimSpace(plan.Id.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, "PAM User record id is empty")
		return
	}

	// Phase 1: record fields (record-update).
	if updateHasMutations(plan, state) {
		cmd := buildRecordUpdatePamUserCommand(recordUID, plan, state)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Phase 2: rotation settings (`pam rotation edit`) when the nested block changed.
	if rotationSettingsNeedApply(plan.RotationSettings, state.RotationSettings) {
		editCmd := buildPamRotationEditCommand(recordUID, plan.RotationSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, ErrDetailRotationEditFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryRotationEditFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
