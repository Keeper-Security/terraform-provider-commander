// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
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

	recordUID := strings.TrimSpace(plan.Id.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryUpdateFailed, "PAM User record id is empty")
		return
	}

	if err := commonrecordsutils.MoveRecordFromSourceToDestination(ctx, r.ApiManager, recordUID, plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if commonpamuser.UpdateHasMutations(plan.PamUserSharedModel, state.PamUserSharedModel) {
		cmd := commonpamuser.BuildUpdateCommand(utils.CmdRecordUpdate, recordUID, plan.PamUserSharedModel, state.PamUserSharedModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, commonpamuser.ErrDetailUpdateFailed); err != nil {
			resp.Diagnostics.AddError(commonpamuser.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	if commonpamuser.RotationSettingsNeedApply(plan.RotationSettings, state.RotationSettings) {
		editCmd := commonpamuser.BuildPamRotationEditCommand(recordUID, plan.RotationSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, commonpamuser.ErrDetailRotationEditFailed); err != nil {
			resp.Diagnostics.AddError(commonpamuser.ErrSummaryRotationEditFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, recordUID, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
