// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *SharedFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
		resp.Diagnostics.AddError(ErrSummarySyncDownFailed, err.Error())
		return
	}

	plan.Id = state.Id
	folderUID := plan.Id.ValueString()

	// Rename shared folder if name changed
	if !plan.Name.Equal(state.Name) {
		name := strings.ReplaceAll(plan.Name.ValueString(), `"`, `\"`)
		command := fmt.Sprintf("%s '%s' %s \"%s\"", CmdRndir, folderUID, FlagName, name)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpRenameSF); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Update default user_permissions / record_permissions if changed
	planPerms := GetDefaultPermissions(&plan)
	statePerms := GetDefaultPermissions(&state)
	if planPerms != statePerms {
		command := BuildSharedFolderDefaultPermissionsCommand(folderUID, planPerms)
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpUpdateDefaultPerms); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	// Sync records: remove removed, grant new/updated
	if err := SyncSharedFolderRecords(ctx, r.ApiManager, folderUID, plan.Records, state.Records); err != nil {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
		return
	}
	// Sync users: remove removed, grant new/updated
	if err := SyncSharedFolderUsers(ctx, r.ApiManager, folderUID, plan.Users, state.Users); err != nil {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// BuildSharedFolderDefaultPermissionsCommand builds share-folder FOLDER_UID --email '*' --record '*' with --manage-users/--manage-records/--can-share/--can-edit on|off.
// Used to apply default user_permissions and record_permissions to the shared folder (create or update).
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
