// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_cloud"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamCloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PamCloudResourceModel
	var state PamCloudResourceModel

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
		resp.Diagnostics.AddError(ErrSummaryPamCloudRecordUpdateFailed, "PAM cloud record id is empty")
		return
	}

	if err := commonrecordsutils.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if commonpamcloud.RecordUpdateHasMutations(plan.PamCloudResourceModel, state.PamCloudResourceModel) {
		cmd := commonpamcloud.BuildUpdateCommand(utils.CmdRecordUpdate, recordUID, plan.PamCloudResourceModel, state.PamCloudResourceModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailPamCloudRecordUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamCloudRecordUpdateFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, recordUID, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryPamCloudRecordUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
