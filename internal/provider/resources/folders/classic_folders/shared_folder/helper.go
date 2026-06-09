// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultPermissionFlags holds the four default permission booleans (user_permissions + record_permissions).
type DefaultPermissionFlags struct {
	ManageUsers   bool
	ManageRecords bool
	CanShare      bool
	CanEdit       bool
}

// validateSharedFolderRecordRefs ensures each records map key matches a vault record_uid or title (see `list --format json`).
func validateSharedFolderRecordRefs(ctx context.Context, apiManager *api.ApiManager, records types.Map) error {
	keys := recordKeysSlice(records)
	if len(keys) == 0 {
		return nil
	}
	return utils.ValidateVaultRecordIdentifiers(ctx, apiManager, keys)
}

// recordKeysSlice returns map keys of m in arbitrary order; empty if m is null/unknown.
func recordKeysSlice(m types.Map) []string {
	km := mapKeys(m)
	out := make([]string, 0, len(km))
	for k := range km {
		out = append(out, k)
	}
	return out
}

// GetDefaultPermissions returns the default permission flags from the resource model (nil/null = false).
func GetDefaultPermissions(data *SharedFolderResourceModel) DefaultPermissionFlags {
	var f DefaultPermissionFlags
	if data.UserPermissions != nil {
		if !data.UserPermissions.ManageUsers.IsNull() {
			f.ManageUsers = data.UserPermissions.ManageUsers.ValueBool()
		}
		if !data.UserPermissions.ManageRecords.IsNull() {
			f.ManageRecords = data.UserPermissions.ManageRecords.ValueBool()
		}
	}
	if data.RecordPermissions != nil {
		if !data.RecordPermissions.CanShare.IsNull() {
			f.CanShare = data.RecordPermissions.CanShare.ValueBool()
		}
		if !data.RecordPermissions.CanEdit.IsNull() {
			f.CanEdit = data.RecordPermissions.CanEdit.ValueBool()
		}
	}
	return f
}

// SyncSharedFolderRecords syncs records with the classic shared folder: grants only added/updated, removes removed.
// Skips grant for items that exist in state with the same value (no change).
func SyncSharedFolderRecords(ctx context.Context, apiManager *api.ApiManager, folderUID string, planRecords, stateRecords types.Map) error {
	planKeys := mapKeys(planRecords)
	stateKeys := mapKeys(stateRecords)
	stateElements := mapElements(stateRecords)

	// Remove: in state but not in plan
	for recordID := range stateKeys {
		if !planKeys[recordID] {
			cmd := buildSharedFolderRecordCommand(ActionRemove, folderUID, recordID, false, false)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpRemoveRecord); err != nil {
				return err
			}
		}
	}

	// Grant: only if added (not in state) or updated (in state but value changed)
	if !planRecords.IsNull() && !planRecords.IsUnknown() {
		for recordID, planVal := range planRecords.Elements() {
			planObj, ok := planVal.(types.Object)
			if !ok {
				return fmt.Errorf("invalid record entry for key %q", recordID)
			}
			if stateVal, inState := stateElements[recordID]; inState {
				if stateObj, ok := stateVal.(types.Object); ok && recordEntryEqual(planObj, stateObj) {
					continue // unchanged, skip grant
				}
			}
			canShare, canEdit := getRecordAttrs(planObj)
			cmd := buildSharedFolderRecordCommand(ActionGrant, folderUID, recordID, canShare, canEdit)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpAddUpdateRecord); err != nil {
				return err
			}
		}
	}
	return nil
}

// SyncSharedFolderUsers syncs users with the classic shared folder: grants only added/updated, removes removed.
// Skips grant for items that exist in state with the same value (no change).
func SyncSharedFolderUsers(ctx context.Context, apiManager *api.ApiManager, folderUID string, planUsers, stateUsers types.Map) error {
	planKeys := mapKeys(planUsers)
	stateKeys := mapKeys(stateUsers)
	stateElements := mapElements(stateUsers)

	// Remove: in state but not in plan
	for emailOrID := range stateKeys {
		if !planKeys[emailOrID] {
			cmd := buildSharedFolderUserCommand(ActionRemove, folderUID, emailOrID, false, false)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpRemoveUser); err != nil {
				return err
			}
		}
	}

	// Grant: only if added (not in state) or updated (in state but value changed)
	if !planUsers.IsNull() && !planUsers.IsUnknown() {
		for emailOrID, planVal := range planUsers.Elements() {
			planObj, ok := planVal.(types.Object)
			if !ok {
				return fmt.Errorf("invalid user entry for key %q", emailOrID)
			}
			if stateVal, inState := stateElements[emailOrID]; inState {
				if stateObj, ok := stateVal.(types.Object); ok && userEntryEqual(planObj, stateObj) {
					continue // unchanged, skip grant
				}
			}
			manageUsers, manageRecords := getUserAttrs(planObj)
			cmd := buildSharedFolderUserCommand(ActionGrant, folderUID, emailOrID, manageUsers, manageRecords)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpAddUpdateUser); err != nil {
				return err
			}
		}
	}
	return nil
}

// mapKeys returns the set of keys in m; empty map if m is null/unknown.
func mapKeys(m types.Map) map[string]bool {
	out := make(map[string]bool)
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	for k := range m.Elements() {
		out[k] = true
	}
	return out
}

// mapElements returns the elements map; nil if m is null/unknown.
func mapElements(m types.Map) map[string]attr.Value {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	return m.Elements()
}

// recordEntryEqual returns true if both objects have the same can_share and can_edit.
func recordEntryEqual(a, b types.Object) bool {
	canShareA, canEditA := getRecordAttrs(a)
	canShareB, canEditB := getRecordAttrs(b)
	return canShareA == canShareB && canEditA == canEditB
}

// userEntryEqual returns true if both objects have the same manage_users and manage_records.
func userEntryEqual(a, b types.Object) bool {
	muA, mrA := getUserAttrs(a)
	muB, mrB := getUserAttrs(b)
	return muA == muB && mrA == mrB
}

func getRecordAttrs(obj types.Object) (canShare, canEdit bool) {
	attrs := obj.Attributes()
	if v, ok := attrs[AttrCanShare].(types.Bool); ok && !v.IsNull() {
		canShare = v.ValueBool()
	}
	if v, ok := attrs[AttrCanEdit].(types.Bool); ok && !v.IsNull() {
		canEdit = v.ValueBool()
	}
	return canShare, canEdit
}

func getUserAttrs(obj types.Object) (manageUsers, manageRecords bool) {
	attrs := obj.Attributes()
	if v, ok := attrs[AttrManageUsers].(types.Bool); ok && !v.IsNull() {
		manageUsers = v.ValueBool()
	}
	if v, ok := attrs[AttrManageRecords].(types.Bool); ok && !v.IsNull() {
		manageRecords = v.ValueBool()
	}
	return manageUsers, manageRecords
}

// buildSharedFolderRecordCommand builds share-folder --action grant|remove SF_ID --record RECORD_ID [--can-share on|off --can-edit on|off].
// For remove, canShare/canEdit are ignored.
func buildSharedFolderRecordCommand(action, folderUID, recordID string, canShare, canEdit bool) string {
	base := fmt.Sprintf("%s %s %s '%s' %s '%s'", CmdShareFolder, FlagAction, action, folderUID, FlagRecord, recordID)
	if action != ActionGrant {
		return base
	}
	canShareVal := ValueOff
	if canShare {
		canShareVal = ValueOn
	}
	canEditVal := ValueOff
	if canEdit {
		canEditVal = ValueOn
	}
	return fmt.Sprintf("%s %s %s %s %s", base, FlagCanShare, canShareVal, FlagCanEdit, canEditVal)
}

// buildSharedFolderUserCommand builds share-folder --action grant|remove SF_ID --email EMAIL_OR_ID [--manage-users on|off --manage-records on|off].
// For remove, permission args are ignored.
func buildSharedFolderUserCommand(action, folderUID, emailOrID string, manageUsers, manageRecords bool) string {
	base := fmt.Sprintf("%s %s %s '%s' %s '%s'", CmdShareFolder, FlagAction, action, folderUID, FlagEmail, emailOrID)
	if action != ActionGrant {
		return base
	}
	parts := []string{base}
	mu := ValueOff
	if manageUsers {
		mu = ValueOn
	}
	parts = append(parts, FlagManageUsers, mu)
	mr := ValueOff
	if manageRecords {
		mr = ValueOn
	}
	parts = append(parts, FlagManageRecords, mr)
	return strings.Join(parts, " ")
}
