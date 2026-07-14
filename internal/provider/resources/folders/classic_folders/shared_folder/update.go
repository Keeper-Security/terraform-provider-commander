// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"
	"fmt"
	"strings"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ClassicSharedFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SharedFolderResourceModel
	var state SharedFolderResourceModel

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

	if err := validateSharedFolderRecordRefs(ctx, r.ApiManager, plan.Records); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryInvalidConfig, err.Error())
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

		command := fmt.Sprintf(`%s "%s" "%s" %s %s`, utils.CmdMv, src, dst, utils.FlagForce, FlagSharedFolder)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpMove); err != nil {
			resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Rename via rndir (after move, so folder is already in the new location).
	if nameChanged {
		leaf := folderutils.EscapeDoubleQuotesForCLI(plan.Name.ValueString())
		command := fmt.Sprintf("%s '%s' %s \"%s\"", CmdRndir, folderUID, FlagName, leaf)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpRename); err != nil {
			resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Update default user_permissions / record_permissions if changed
	planPerms := GetDefaultPermissions(&plan)
	statePerms := GetDefaultPermissions(&state)
	if planPerms != statePerms {
		command := BuildSharedFolderDefaultPermissionsCommand(folderUID, planPerms)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpUpdateDefaultPerms); err != nil {
			resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Sync records: remove removed, grant new/updated
	if err := SyncSharedFolderRecords(ctx, r.ApiManager, folderUID, plan.Records, state.Records); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
		return
	}
	// Sync users: remove removed, grant new/updated
	if err := SyncSharedFolderUsers(ctx, r.ApiManager, folderUID, plan.Users, state.Users); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// BuildSharedFolderDefaultPermissionsCommand builds share-folder FOLDER_UID --email '*' --record '*' with --manage-users/--manage-records/--can-share/--can-edit on|off.
// Used to apply default user_permissions and record_permissions to the classic shared folder (create or update).
func BuildSharedFolderDefaultPermissionsCommand(folderUID string, f DefaultPermissionFlags) string {
	onOff := func(b bool) string {
		if b {
			return ValueOn
		}
		return ValueOff
	}
	parts := []string{
		fmt.Sprintf("%s '%s' %s '%s' %s '%s'", CmdShareFolder, folderUID, FlagEmail, WildcardAll, FlagRecord, WildcardAll),
		FlagManageUsers, onOff(f.ManageUsers),
		FlagManageRecords, onOff(f.ManageRecords),
		FlagCanShare, onOff(f.CanShare),
		FlagCanEdit, onOff(f.CanEdit),
	}
	return strings.Join(parts, " ")
}
