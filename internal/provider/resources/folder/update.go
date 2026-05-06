// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FolderResourceModel
	var state FolderResourceModel

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
		resp.Diagnostics.AddError(ErrSummarySyncDownFailed, err.Error())
		return
	}

	plan.Id = state.Id
	folderUID := plan.Id.ValueString()

	nameChanged := !plan.Name.Equal(state.Name)
	locationChanged := !plan.FolderLocation.Equal(state.FolderLocation)
	colorChanged := !plan.Color.Equal(state.Color)

	// Move first (before rename) so the source path using the old name is still valid.
	if locationChanged {
		statePath := BuildFolderPath(state.Name.ValueString(), state.FolderLocation.ValueString())
		planPath := BuildFolderPath(state.Name.ValueString(), plan.FolderLocation.ValueString())
		src := EscapeDoubleQuotesForCLI(MvPathForCommander(statePath))
		dst := EscapeDoubleQuotesForCLI(MvMoveTargetParent(planPath))
		command := fmt.Sprintf(`%s "%s" "%s" %s`, CmdMv, src, dst, FlagForce)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpMoveFolder); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Rename and/or color change via rndir (after move, so folder is already in the new location).
	if nameChanged || colorChanged {
		command := buildRndirCommand(folderUID, &plan, nameChanged, colorChanged)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpRenameFolder); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Sync records: link added, unlink removed
	folderName := BuildFolderPath(plan.Name.ValueString(), plan.FolderLocation.ValueString())
	if err := SyncFolderRecords(ctx, r.ApiManager, folderUID, folderName, plan.Records, state.Records); err != nil {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// buildRndirCommand builds: rndir '<folderUID>' [-n <newName>] [--color <color>] -q.
func buildRndirCommand(folderUID string, plan *FolderResourceModel, nameChanged, colorChanged bool) string {
	parts := []string{fmt.Sprintf("%s '%s'", CmdRndir, folderUID)}

	if nameChanged {
		newName := EscapeDoubleQuotesForCLI(plan.Name.ValueString())
		parts = append(parts, FlagName, fmt.Sprintf(`"%s"`, newName))
	}

	if colorChanged {
		color := ColorNone
		if !plan.Color.IsNull() && !plan.Color.IsUnknown() && plan.Color.ValueString() != "" {
			color = plan.Color.ValueString()
		}
		parts = append(parts, FlagColor, color)
	}

	parts = append(parts, FlagQuiet)
	return strings.Join(parts, " ")
}
