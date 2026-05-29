// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"
	"fmt"
	"strings"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *NonSharedFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NonSharedFolderResourceModel
	var state NonSharedFolderResourceModel

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
	folderUID := plan.Id.ValueString()

	nameChanged := !plan.Name.Equal(state.Name)

	locationChanged := !plan.FolderLocation.Equal(state.FolderLocation)

	// Move first (before rename) so the source path using the old name is still valid.
	if locationChanged {
		statePath := folderutils.BuildFolderPath(state.Name.ValueString(), state.FolderLocation.ValueString())
		planPath := folderutils.BuildFolderPath(state.Name.ValueString(), plan.FolderLocation.ValueString())
		src := folderutils.EscapeDoubleQuotesForCLI(folderutils.MvPathForCommander(statePath))
		dst := folderutils.EscapeDoubleQuotesForCLI(folderutils.MvMoveTargetParent(planPath))

		command := fmt.Sprintf(`%s "%s" "%s" %s`, utils.CmdMv, src, dst, utils.FlagForce)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpMove); err != nil {
			resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Rename via rndir (after move, so folder is already in the new location).
	if nameChanged {
		command := buildRndirCommand(folderUID, &plan)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpRename); err != nil {
			resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Sync records: link added, unlink removed
	folderName := folderutils.BuildFolderPath(plan.Name.ValueString(), plan.FolderLocation.ValueString())
	if err := SyncFolderRecords(ctx, r.ApiManager, folderUID, folderName, plan.Records, state.Records); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// buildRndirCommand builds: rndir '<folderUID>' -n "<newName>" -q.
func buildRndirCommand(folderUID string, plan *NonSharedFolderResourceModel) string {
	newName := folderutils.EscapeDoubleQuotesForCLI(plan.Name.ValueString())
	parts := []string{
		fmt.Sprintf("%s '%s'", CmdRndir, folderUID),
		FlagName,
		fmt.Sprintf(`"%s"`, newName),
		FlagQuiet,
	}
	return strings.Join(parts, " ")
}
