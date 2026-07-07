// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *SshKeysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SshKeysResourceModel
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
	uid := strings.TrimSpace(plan.Id.ValueString())
	if uid == "" {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, "SSH keys record id is empty")
		return
	}

	if err := commonrecordsutils.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if commonrecordsshkeys.UpdateHasMutations(plan.SshKeysModel, state.SshKeysModel) {
		cmd := commonrecordsshkeys.BuildUpdateCommand(utils.CmdRecordUpdate, uid, plan.SshKeysModel, state.SshKeysModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, uid, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
