// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package shared_folder

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ErrSharedFolderNotFound is returned when the get command yields no usable folder (deleted, wrong id/path, or empty response).
var ErrSharedFolderNotFound = errors.New("classic shared folder not found")

// GetSharedFolderCommand builds the Commander CLI: get '<name-or-uid>' --format json (Commander accepts classic shared folder UID or vault path).
func GetSharedFolderCommand(nameOrUID string) string {
	return fmt.Sprintf("%s '%s' %s %s", cmdGet, nameOrUID, flagFormat, formatJSON)
}

// ParseSharedFolderResponse unmarshals API result data into SharedFolderResponse.
// If data is nil, returns (nil, nil) so callers can treat as missing resource.
func ParseSharedFolderResponse(data any) (*SharedFolderResponse, error) {
	if data == nil {
		return nil, nil
	}
	var out SharedFolderResponse
	if err := utils.UnmarshalApiResponse(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchSharedFolderByNameOrId loads a classic shared folder by vault path/name or UID; fails if the command errors or response has no folder.
func FetchSharedFolderByNameOrId(ctx context.Context, apiManager *api.ApiManager, nameOrUID string) (*SharedFolderResponse, error) {
	if nameOrUID == "" {
		return nil, fmt.Errorf("classic shared folder name or id is empty")
	}
	apiResp, err := apiManager.ExecuteCommand(ctx, GetSharedFolderCommand(nameOrUID), errOpGet)
	if err != nil {
		return nil, err
	}
	if apiResp == nil || apiResp.Data == nil {
		return nil, ErrSharedFolderNotFound
	}
	out, err := ParseSharedFolderResponse(apiResp.Data)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, ErrSharedFolderNotFound
	}
	if out.FolderUID == "" {
		return nil, fmt.Errorf("get response missing folder_uid")
	}
	return out, nil
}

// MapResponseToModel maps unmarshaled get classic shared folder API data into m. priorUsers is used to
// preserve user map keys from existing state (resource refresh); pass null Map for data sources.
func MapResponseToModel(api *SharedFolderResponse, m *Model, priorUsers types.Map, priorRecords types.Map) error {
	if api == nil {
		return fmt.Errorf("classic shared folder API response is nil")
	}

	m.Id = types.StringValue(api.FolderUID)
	m.Name = types.StringValue(api.Name)
	m.FolderLocation = utils.ExtractFolderValue(&api.FolderLocation, m.FolderLocation)

	m.UserPermissions = &UserPermissionsModel{
		ManageUsers:   types.BoolValue(api.ManageUsers),
		ManageRecords: types.BoolValue(api.ManageRecords),
	}
	m.RecordPermissions = &RecordPermissionsModel{
		CanShare: types.BoolValue(api.CanShare),
		CanEdit:  types.BoolValue(api.CanEdit),
	}

	recordsMap, err := buildRecordsMapFromAPIResponse(api.Records, priorRecords)
	if err != nil {
		return fmt.Errorf("records: %w", err)
	}
	m.Records = recordsMap

	usersMap, err := buildUsersMapFromAPIResponse(api.Users, api.Teams, priorUsers)
	if err != nil {
		return fmt.Errorf("users: %w", err)
	}
	m.Users = usersMap

	return nil
}

// recordEntryMapKey picks the records map key for an API row. If prior state already keys this record
// by record_uid or record_name, the same key is reused; otherwise record_name is preferred, then record_uid.
func recordEntryMapKey(rec SharedFolderRecordEntry, priorRecords types.Map) string {
	if !priorRecords.IsNull() && !priorRecords.IsUnknown() {
		for k := range priorRecords.Elements() {
			if rec.RecordUID != "" && k == rec.RecordUID {
				return k
			}
			if rec.RecordName != "" && k == rec.RecordName {
				return k
			}
		}
	}
	if rec.RecordName != "" {
		return rec.RecordName
	}
	return rec.RecordUID
}

func buildRecordsMapFromAPIResponse(entries []SharedFolderRecordEntry, priorRecords types.Map) (types.Map, error) {
	elements := make(map[string]attr.Value)
	for _, rec := range entries {
		key := recordEntryMapKey(rec, priorRecords)
		if key == "" {
			continue
		}

		elements[key] = types.ObjectValueMust(
			map[string]attr.Type{AttrCanShare: types.BoolType, AttrCanEdit: types.BoolType},
			map[string]attr.Value{
				AttrCanShare: types.BoolValue(rec.CanShare),
				AttrCanEdit:  types.BoolValue(rec.CanEdit),
			},
		)
	}
	mapVal, diags := types.MapValue(RecordEntryMapElemType, elements)
	if diags.HasError() {
		return types.MapNull(RecordEntryMapElemType), fmt.Errorf("failed to build records map: %v", diags)
	}
	return mapVal, nil
}

// userOrTeamEntryMapKey picks the users-map key for an API row that may represent
// either a user (Username/UserID populated) or a team (TeamName/TeamId populated).
// If prior state already keys this row by name or id, the same key is reused so
// the model stays stable across refreshes; otherwise the populated name is
// preferred, then the populated id.
func userOrTeamEntryMapKey(e SharedFolderUserTeamEntry, priorUsers types.Map) string {
	name := e.Username
	if name == "" {
		name = e.TeamName
	}
	id := e.UserID
	if id == "" {
		id = e.TeamId
	}
	if !priorUsers.IsNull() && !priorUsers.IsUnknown() {
		for k := range priorUsers.Elements() {
			if id != "" && k == id {
				return k
			}
			if name != "" && k == name {
				return k
			}
		}
	}
	if name != "" {
		return name
	}
	return id
}

// buildUsersMapFromAPIResponse projects the API's users and teams arrays into a
// single Terraform users map keyed by user email/UID or team name/UID. Commander's
// share-folder --email flag resolves either form, so users and teams share one
// block on both the read and write paths.
func buildUsersMapFromAPIResponse(usersEntries []SharedFolderUserTeamEntry, teamsEntries []SharedFolderUserTeamEntry, priorUsers types.Map) (types.Map, error) {
	elements := make(map[string]attr.Value, len(usersEntries)+len(teamsEntries))
	for _, e := range usersEntries {
		// Skip owner from the map
		if e.Owner {
			continue
		}

		addUserOrTeamEntry(elements, e, priorUsers)
	}
	for _, e := range teamsEntries {
		addUserOrTeamEntry(elements, e, priorUsers)
	}
	mapVal, diags := types.MapValue(UserEntryMapElemType, elements)
	if diags.HasError() {
		return types.MapNull(UserEntryMapElemType), fmt.Errorf("failed to build users map: %v", diags)
	}
	return mapVal, nil
}

// addUserOrTeamEntry inserts e into elements under its computed key.
func addUserOrTeamEntry(elements map[string]attr.Value, e SharedFolderUserTeamEntry, priorUsers types.Map) {
	key := userOrTeamEntryMapKey(e, priorUsers)
	if key == "" {
		return
	}
	elements[key] = types.ObjectValueMust(
		map[string]attr.Type{
			AttrManageUsers:   types.BoolType,
			AttrManageRecords: types.BoolType,
		},
		map[string]attr.Value{
			AttrManageUsers:   types.BoolValue(e.ManageUsers),
			AttrManageRecords: types.BoolValue(e.ManageRecords),
		},
	)
}
