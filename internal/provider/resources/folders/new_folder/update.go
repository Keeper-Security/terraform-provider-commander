// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"
	"fmt"

	commonfolderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *NewFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NewFolderResourceModel
	var state NewFolderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	plan.Id = state.Id

	if err := utils.MoveNsfFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(ErrSummaryNsfMoveFolderFailed, err.Error())
		return
	}

	if !plan.Name.Equal(state.Name) {
		command := fmt.Sprintf(`%s "%s" %s="%s"`, CmdNsfRndir, state.Id.ValueString(), commonfolderutils.FlagName, plan.Name.ValueString())
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, commonfolderutils.ErrOpRename); err != nil {
			resp.Diagnostics.AddError(commonfolderutils.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Sync the share permissions - share the folder to the users in the share block.
	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdNsfShareFolder, plan.Id.ValueString(), plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(commonfolderutils.ErrSummaryUpdateFailed, err.Error())
		return
	}

	// Link the records to the folder.
	folderUID := plan.Id.ValueString()
	unlinkFn := func(record string) error {
		return commonfolderutils.UnlinkRecordNsf(ctx, r.ApiManager, CmdNsfRm, folderUID, record)
	}
	if err := commonfolderutils.SyncFolderRecords(ctx, r.ApiManager, CmdNsfLn, folderUID, plan.Records, state.Records, unlinkFn); err != nil {
		resp.Diagnostics.AddError(commonfolderutils.ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
