// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SplitSharedFolderPath splits a full vault path into parent (all but the last segment) and leaf (shared folder name).
// Example: "Templates/My Shared Folder 1" -> parent "Templates", leaf "My Shared Folder 1".
// A path with no "/" has an empty parent and the whole string as leaf.
func SplitSharedFolderPath(full string) (parent, leaf string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}
	i := strings.LastIndex(full, "/")
	if i < 0 {
		return "", full
	}
	return strings.TrimSpace(full[:i]), strings.TrimSpace(full[i+1:])
}

// EscapeDoubleQuotesForCLI escapes double quotes for use inside double-quoted shell arguments.
func EscapeDoubleQuotesForCLI(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// MvPathForCommander normalizes a vault path for Commander `mv`. Paths with no parent
// (no `/` — e.g. "My Shared Folder 1" at vault root) are prefixed with `/` so the CLI
// targets the root folder. Paths that already start with `/` or contain `/` are unchanged.
func MvPathForCommander(full string) string {
	full = strings.TrimSpace(full)
	if full == "" {
		return full
	}
	if strings.HasPrefix(full, "/") {
		return full
	}
	parent, leaf := SplitSharedFolderPath(full)
	if parent == "" {
		return "/" + leaf
	}
	return full
}

// MvMoveTargetParent returns the second argument to Commander `mv`: the destination parent folder only
// (not the shared folder leaf). Example: plan "Templates/test4/My Shared Folder 1" -> "Templates/test4".
// Plan "My Shared Folder 1" (vault root, no parent path) -> "/".
func MvMoveTargetParent(planPath string) string {
	planPath = strings.TrimSpace(planPath)
	if planPath == "" {
		return planPath
	}
	trim := planPath
	if strings.HasPrefix(trim, "/") {
		trim = strings.TrimSpace(trim[1:])
	}
	parent, _ := SplitSharedFolderPath(trim)
	parent = strings.TrimSpace(parent)
	if parent == "" {
		return "/"
	}
	return parent
}

// DefaultPermissionFlags holds the four default permission booleans (user_permissions + record_permissions).
type DefaultPermissionFlags struct {
	ManageUsers   bool
	ManageRecords bool
	CanShare      bool
	CanEdit       bool
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

// SyncSharedFolderRecords syncs records with the shared folder: grants only added/updated, removes removed.
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

// SyncSharedFolderUsers syncs users with the shared folder: grants only added/updated, removes removed.
// Skips grant for items that exist in state with the same value (no change).
func SyncSharedFolderUsers(ctx context.Context, apiManager *api.ApiManager, folderUID string, planUsers, stateUsers types.Map) error {
	planKeys := mapKeys(planUsers)
	stateKeys := mapKeys(stateUsers)
	stateElements := mapElements(stateUsers)

	// Remove: in state but not in plan
	for emailOrID := range stateKeys {
		if !planKeys[emailOrID] {
			cmd := buildSharedFolderUserCommand(ActionRemove, folderUID, emailOrID, false, false, "")
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
			manageUsers, manageRecords, expiration := getUserAttrs(planObj)
			cmd := buildSharedFolderUserCommand(ActionGrant, folderUID, emailOrID, manageUsers, manageRecords, expiration)
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

// userEntryEqual returns true if both objects have the same manage_users, manage_records, and expiration.
func userEntryEqual(a, b types.Object) bool {
	muA, mrA, expA := getUserAttrs(a)
	muB, mrB, expB := getUserAttrs(b)
	return muA == muB && mrA == mrB && strings.TrimSpace(expA) == strings.TrimSpace(expB)
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

func getUserAttrs(obj types.Object) (manageUsers, manageRecords bool, expiration string) {
	attrs := obj.Attributes()
	if v, ok := attrs[AttrManageUsers].(types.Bool); ok && !v.IsNull() {
		manageUsers = v.ValueBool()
	}
	if v, ok := attrs[AttrManageRecords].(types.Bool); ok && !v.IsNull() {
		manageRecords = v.ValueBool()
	}
	if v, ok := attrs[AttrExpiration].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
		expiration = v.ValueString()
	}
	return manageUsers, manageRecords, expiration
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

// buildSharedFolderUserCommand builds share-folder --action grant|remove SF_ID --email EMAIL_OR_ID [--manage-users on|off --manage-records on|off [--expire-at|--expire-in VALUE]].
// For remove, permission and expiration args are ignored.
func buildSharedFolderUserCommand(action, folderUID, emailOrID string, manageUsers, manageRecords bool, expiration string) string {
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
	expFlag, expVal := expirationFlagAndValue(expiration)
	if expFlag != "" && expVal != "" {
		parts = append(parts, expFlag, fmt.Sprintf("'%s'", expVal))
	}
	return strings.Join(parts, " ")
}

// expirationFlagAndValue returns (--expire-at, value) for yyyy-MM-ddTHH:mm:ss, or ("", "") when empty.
func expirationFlagAndValue(exp string) (flag string, value string) {
	exp = strings.TrimSpace(exp)
	if exp == "" {
		return "", ""
	}
	if expirationNever.MatchString(exp) {
		return FlagExpireAt, ValueNever
	}
	if t, err := time.Parse(TimeLayoutExpiration, exp); err == nil {
		return FlagExpireAt, t.Format(TimeLayoutExpiration)
	}
	return "", ""
}
