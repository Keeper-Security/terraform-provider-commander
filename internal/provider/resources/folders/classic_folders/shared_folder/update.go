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

	// Name is "parent/Shared Folder Name" or just "Shared Folder Name". Same parent + different leaf -> rndir; different parent -> mv.
	if !plan.Name.Equal(state.Name) {
		planPath := plan.Name.ValueString()
		statePath := state.Name.ValueString()
		planParent, planLeaf := SplitSharedFolderPath(planPath)
		stateParent, stateLeaf := SplitSharedFolderPath(statePath)

		if planParent == stateParent && planLeaf != stateLeaf {
			leaf := EscapeDoubleQuotesForCLI(planLeaf)
			command := fmt.Sprintf("%s '%s' %s \"%s\"", CmdRndir, folderUID, FlagName, leaf)
			if _, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpRename); err != nil {
				resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
				return
			}
		} else if planParent != stateParent {
			src := EscapeDoubleQuotesForCLI(MvPathForCommander(statePath))
			dst := EscapeDoubleQuotesForCLI(MvMoveTargetParent(planPath))
			command := fmt.Sprintf(`%s "%s" "%s" %s %s`, utils.CmdMv, src, dst, utils.FlagForce, FlagSharedFolder)
			if _, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpMove); err != nil {
				resp.Diagnostics.AddError(folderutils.ErrSummaryUpdateFailed, err.Error())
				return
			}
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
