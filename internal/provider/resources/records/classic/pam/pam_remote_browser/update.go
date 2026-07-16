// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_remote_browser"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamRemoteBrowserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PamRemoteBrowserResourceModel
	var state PamRemoteBrowserResourceModel

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
		resp.Diagnostics.AddError(ErrSummaryPamRemoteBrowserRecordUpdateFailed, "PAM remote browser record id is empty")
		return
	}

	if err := commonrecordsutils.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if commonpamremotebrowser.RecordUpdateHasMutations(plan.PamRemoteBrowserResourceModel, state.PamRemoteBrowserResourceModel) {
		cmd := commonpamremotebrowser.BuildUpdateCommand(utils.CmdRecordUpdate, recordUID, plan.PamRemoteBrowserResourceModel, state.PamRemoteBrowserResourceModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailPamRemoteBrowserRecordUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamRemoteBrowserRecordUpdateFailed, err.Error())
			return
		}
	}

	if commonpamremotebrowser.PamRemoteBrowserSettingsNeedApply(plan.PamRemoteBrowserSettings, state.PamRemoteBrowserSettings) {
		editCmd := commonpamremotebrowser.BuildPamRbiEditCommand(recordUID, plan.PamRemoteBrowserSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, ErrDetailPamRbiEditFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamRbiEditFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, recordUID, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryPamRemoteBrowserRecordUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
