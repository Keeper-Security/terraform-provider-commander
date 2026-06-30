// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdatabase

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_database"
	commonrecordutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PamDatabaseResourceModel
	var state PamDatabaseResourceModel

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

	// Throw error if user tries to change the folder location as it is not supported
	if !plan.FolderLocation.Equal(state.FolderLocation) {
		resp.Diagnostics.AddError(commonrecordutils.ErrSummaryInvalidConfig, commonrecordutils.ErrSummaryMoveNotSupported)
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(commonpamrecords.ValidateDatabasePamSettingsFieldsNotRemoved(plan.PamSettings, state.PamSettings)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id
	recordUID := strings.TrimSpace(plan.Id.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryPamDatabaseRecordUpdateFailed, "new PAM database record id is empty")
		return
	}

	if err := commonrecordutils.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if commonpamdatabase.RecordUpdateHasMutations(plan.PamDatabaseResourceModel, state.PamDatabaseResourceModel) {
		cmd := commonpamdatabase.BuildUpdateCommand(utils.CmdNsfRecordUpdate, recordUID, plan.PamDatabaseResourceModel, state.PamDatabaseResourceModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailPamDatabaseRecordUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamDatabaseRecordUpdateFailed, err.Error())
			return
		}
	}

	if plan.PamSettings != nil {
		if err := commonpamrecords.ApplyDatabasePamSettings(ctx, r.ApiManager, recordUID, plan.PamSettings, state.PamSettings); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamSettingsFailed, err.Error())
			return
		}
	}

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdNsfShareRecord, recordUID, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryPamDatabaseRecordUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
