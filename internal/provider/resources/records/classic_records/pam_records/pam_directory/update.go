// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamDirectoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PamDirectoryResourceModel
	var state PamDirectoryResourceModel

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

	resp.Diagnostics.Append(commonpamrecords.ValidatePamSettingsFieldsNotRemoved(plan.PamSettings, state.PamSettings)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id
	recordUID := strings.TrimSpace(plan.Id.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryPamDirectoryRecordUpdateFailed, "PAM directory record id is empty")
		return
	}

	if err := commonpamrecords.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.Folder.ValueString(), state.Folder.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if commonpamdirectory.RecordUpdateHasMutations(plan.PamDirectoryResourceModel, state.PamDirectoryResourceModel) {
		cmd := commonpamdirectory.BuildUpdateCommand(utils.CmdRecordUpdate, recordUID, plan.PamDirectoryResourceModel, state.PamDirectoryResourceModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailPamDirectoryRecordUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamDirectoryRecordUpdateFailed, err.Error())
			return
		}
	}

	if plan.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, recordUID, plan.PamSettings, state.PamSettings); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamSettingsFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, recordUID, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryPamDirectoryRecordUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
