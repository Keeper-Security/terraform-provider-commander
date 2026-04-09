// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *SharedFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SharedFolderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
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

	if err := validateSharedFolderRecordRefs(ctx, r.ApiManager, data.Records); err != nil {
		resp.Diagnostics.AddError(ErrSummaryInvalidConfig, err.Error())
		return
	}

	// Phase 1: create shared folder with name, folder_location, user_permissions, record_permissions
	command, err := buildCreateSharedFolderCommand(&data)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryInvalidConfig, err.Error())
		return
	}

	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpCreateSF)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	folderUID, err := extractFolderUIDFromCreateSharedFolderResponse(apiResp.Data)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}
	data.Id = types.StringValue(folderUID)

	// Phase 2: sync records and users (grant all; no prior state so nothing to remove)
	if err := SyncSharedFolderRecords(ctx, r.ApiManager, folderUID, data.Records, types.MapNull(RecordEntryMapElemType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}
	if err := SyncSharedFolderUsers(ctx, r.ApiManager, folderUID, data.Users, types.MapNull(UserEntryMapElemType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildMkdirCommand builds the "mkdir --shared-folder" command and the folder path used as NAME.
// Folder path: if folder_location is set, "folder_location/name", otherwise "name".
func buildCreateSharedFolderCommand(data *SharedFolderResourceModel) (command string, err error) {
	name := data.Name.ValueString()
	if name == "" {
		return "", fmt.Errorf("name is required")
	}

	parts := []string{CmdMkdir, FlagSharedFolder, fmt.Sprintf(`"%s"`, name)}
	permFlags := GetDefaultPermissions(data)
	parts = append(parts, DefaultPermissionFlagsForMkdir(permFlags)...)

	command = strings.Join(parts, " ")
	return command, nil
}

// extractFolderUIDFromResponse extracts folder_uid from the API response Data (map from JSON).
// If folder_uid is not a string (e.g. number from JSON), it is converted to string.
func extractFolderUIDFromCreateSharedFolderResponse(data any) (string, error) {
	m, _ := data.(map[string]interface{})
	v, ok := m[KeyFolderUID]
	if !ok || v == nil {
		return "", fmt.Errorf("API response missing %s", KeyFolderUID)
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return "", fmt.Errorf("API response %s is empty", KeyFolderUID)
		}
		return s, nil
	}
	return fmt.Sprintf("%v", v), nil
}

// DefaultPermissionFlagsForMkdir returns flag names only for true values (for mkdir --manage-users etc.).
func DefaultPermissionFlagsForMkdir(f DefaultPermissionFlags) []string {
	var parts []string
	if f.ManageUsers {
		parts = append(parts, FlagManageUsers)
	}
	if f.ManageRecords {
		parts = append(parts, FlagManageRecords)
	}
	if f.CanShare {
		parts = append(parts, FlagCanShare)
	}
	if f.CanEdit {
		parts = append(parts, FlagCanEdit)
	}
	return parts
}
