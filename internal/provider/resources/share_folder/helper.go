// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RecordEntryMapElemType is the object type for each entry in the records map (can_share, can_edit).
// Used for MapNull in Create when state has no records.
var RecordEntryMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		AttrCanShare: types.BoolType,
		AttrCanEdit:  types.BoolType,
	},
}

// UserEntryMapElemType is the object type for each entry in the users map (manage_users, manage_records, expiration).
var UserEntryMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		AttrManageUsers:   types.BoolType,
		AttrManageRecords: types.BoolType,
		AttrExpiration:    types.StringType,
	},
}

// DefaultPermissionFlags holds the four default permission booleans (user_permissions + record_permissions).
type DefaultPermissionFlags struct {
	ManageUsers   bool
	ManageRecords bool
	CanShare      bool
	CanEdit       bool
}

// GetDefaultPermissions returns the default permission flags from the resource model (nil/null = false).
func GetDefaultPermissions(data *ShareFolderResourceModel) DefaultPermissionFlags {
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

// SyncShareFolderRecords syncs records with the shared folder: grants only added/updated, removes removed.
// Skips grant for items that exist in state with the same value (no change).
func SyncShareFolderRecords(ctx context.Context, apiManager *api.ApiManager, folderUID string, planRecords, stateRecords types.Map) error {
	planKeys := mapKeys(planRecords)
	stateKeys := mapKeys(stateRecords)
	stateElements := mapElements(stateRecords)

	// Remove: in state but not in plan
	for recordID := range stateKeys {
		if !planKeys[recordID] {
			cmd := buildShareFolderRecordCommand(ActionRemove, folderUID, recordID, false, false)
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
			cmd := buildShareFolderRecordCommand(ActionGrant, folderUID, recordID, canShare, canEdit)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpAddUpdateRecord); err != nil {
				return err
			}
		}
	}
	return nil
}

// SyncShareFolderUsers syncs users with the shared folder: grants only added/updated, removes removed.
// Skips grant for items that exist in state with the same value (no change).
func SyncShareFolderUsers(ctx context.Context, apiManager *api.ApiManager, folderUID string, planUsers, stateUsers types.Map) error {
	planKeys := mapKeys(planUsers)
	stateKeys := mapKeys(stateUsers)
	stateElements := mapElements(stateUsers)

	// Remove: in state but not in plan
	for emailOrID := range stateKeys {
		if !planKeys[emailOrID] {
			cmd := buildShareFolderUserCommand(ActionRemove, folderUID, emailOrID, false, false, "")
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
			cmd := buildShareFolderUserCommand(ActionGrant, folderUID, emailOrID, manageUsers, manageRecords, expiration)
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

// buildShareFolderRecordCommand builds share-folder --action grant|remove SF_ID --record RECORD_ID [--can-share on|off --can-edit on|off].
// For remove, canShare/canEdit are ignored.
func buildShareFolderRecordCommand(action, folderUID, recordID string, canShare, canEdit bool) string {
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

// buildShareFolderUserCommand builds share-folder --action grant|remove SF_ID --email EMAIL_OR_ID [--manage-users on|off --manage-records on|off [--expire-at|--expire-in VALUE]].
// For remove, permission and expiration args are ignored.
func buildShareFolderUserCommand(action, folderUID, emailOrID string, manageUsers, manageRecords bool, expiration string) string {
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

// expirationFlagAndValue returns (--expire-at, value) for ISO date/datetime, (--expire-in, value) for relative or "never", or ("", "") when empty.
func expirationFlagAndValue(exp string) (flag string, value string) {
	exp = strings.TrimSpace(exp)
	if exp == "" {
		return "", ""
	}
	if expirationNever.MatchString(exp) {
		return FlagExpireIn, ValueNever
	}
	if expirationRelative.MatchString(exp) {
		return FlagExpireIn, exp
	}
	if expirationISO.MatchString(exp) {
		return FlagExpireAt, exp
	}
	return "", ""
}

// getStringFromMap returns a string value from m[key], or empty string if missing/invalid.
func getStringFromMap(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// getBoolFromMap returns a bool from m[key]; JSON may give bool or number. Missing/nil => false.
func getBoolFromMap(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	if n, ok := v.(float64); ok {
		return n != 0
	}
	return false
}

// MapGetResponseToModel maps the "get ID --format json" API response (Data) to ShareFolderResourceModel.
// data is expected to be map[string]interface{} with shared_folder_uid, name, manage_users, manage_records,
// can_edit, can_share, records ([]interface{}), users ([]interface{}). Folder_location is not in API and set null.
func MapGetResponseToModel(ctx context.Context, data any) (*ShareFolderResourceModel, error) {
	m, ok := data.(map[string]interface{})
	if !ok || m == nil {
		return nil, fmt.Errorf("get response data is not a map")
	}

	uid := getStringFromMap(m, KeySharedFolderUID)
	if uid == "" {
		return nil, fmt.Errorf("get response missing %s", KeySharedFolderUID)
	}

	model := &ShareFolderResourceModel{
		Id:             types.StringValue(uid),
		Name:           types.StringValue(getStringFromMap(m, KeyName)),
		FolderLocation: types.StringNull(),
		UserPermissions: &UserPermissionsModel{
			ManageUsers:   types.BoolValue(getBoolFromMap(m, AttrManageUsers)),
			ManageRecords: types.BoolValue(getBoolFromMap(m, AttrManageRecords)),
		},
		RecordPermissions: &RecordPermissionsModel{
			CanShare: types.BoolValue(getBoolFromMap(m, AttrCanShare)),
			CanEdit:  types.BoolValue(getBoolFromMap(m, AttrCanEdit)),
		},
	}

	// records: key = record_uid, value = { can_share, can_edit }
	recordsRaw, _ := m[KeyRecords]
	recordsMap, err := buildRecordsMapFromGetResponse(recordsRaw)
	if err != nil {
		return nil, fmt.Errorf("records: %w", err)
	}
	model.Records = recordsMap

	// users: key = username, value = { manage_users, manage_records, expiration }; expiration not in API => null
	usersRaw, _ := m[KeyUsers]
	usersMap, err := buildUsersMapFromGetResponse(usersRaw)
	if err != nil {
		return nil, fmt.Errorf("users: %w", err)
	}
	model.Users = usersMap

	return model, nil
}

func buildRecordsMapFromGetResponse(recordsRaw interface{}) (types.Map, error) {
	elements := make(map[string]attr.Value)
	if recordsRaw == nil {
		mapVal, diags := types.MapValue(RecordEntryMapElemType, elements)
		if diags.HasError() {
			return types.MapNull(RecordEntryMapElemType), fmt.Errorf("failed to build records map: %v", diags)
		}
		return mapVal, nil
	}
	slice, ok := recordsRaw.([]interface{})
	if !ok {
		return types.MapNull(RecordEntryMapElemType), fmt.Errorf("records is not an array")
	}
	for _, item := range slice {
		rec, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		key := getStringFromMap(rec, KeyRecordUID)
		if key == "" {
			continue
		}
		canShare := getBoolFromMap(rec, AttrCanShare)
		canEdit := getBoolFromMap(rec, AttrCanEdit)
		elements[key] = types.ObjectValueMust(
			map[string]attr.Type{AttrCanShare: types.BoolType, AttrCanEdit: types.BoolType},
			map[string]attr.Value{
				AttrCanShare: types.BoolValue(canShare),
				AttrCanEdit:  types.BoolValue(canEdit),
			},
		)
	}
	mapVal, diags := types.MapValue(RecordEntryMapElemType, elements)
	if diags.HasError() {
		return types.MapNull(RecordEntryMapElemType), fmt.Errorf("failed to build records map: %v", diags)
	}
	return mapVal, nil
}

func buildUsersMapFromGetResponse(usersRaw interface{}) (types.Map, error) {
	elements := make(map[string]attr.Value)
	if usersRaw == nil {
		mapVal, diags := types.MapValue(UserEntryMapElemType, elements)
		if diags.HasError() {
			return types.MapNull(UserEntryMapElemType), fmt.Errorf("failed to build users map: %v", diags)
		}
		return mapVal, nil
	}
	slice, ok := usersRaw.([]interface{})
	if !ok {
		return types.MapNull(UserEntryMapElemType), fmt.Errorf("users is not an array")
	}
	for _, item := range slice {
		u, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		key := getStringFromMap(u, KeyUsername)
		if key == "" {
			continue
		}
		manageUsers := getBoolFromMap(u, AttrManageUsers)
		manageRecords := getBoolFromMap(u, AttrManageRecords)
		elements[key] = types.ObjectValueMust(
			map[string]attr.Type{
				AttrManageUsers:   types.BoolType,
				AttrManageRecords: types.BoolType,
				AttrExpiration:    types.StringType,
			},
			map[string]attr.Value{
				AttrManageUsers:   types.BoolValue(manageUsers),
				AttrManageRecords: types.BoolValue(manageRecords),
				AttrExpiration:    types.StringNull(),
			},
		)
	}
	mapVal, diags := types.MapValue(UserEntryMapElemType, elements)
	if diags.HasError() {
		return types.MapNull(UserEntryMapElemType), fmt.Errorf("failed to build users map: %v", diags)
	}
	return mapVal, nil
}
