// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package shared_folder

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SharedFolderRecordEntry is one element of the records array from get classic shared folder --format json.
type SharedFolderRecordEntry struct {
	RecordUID  string `json:"record_uid"`
	RecordName string `json:"record_name"`
	CanShare   bool   `json:"can_share"`
	CanEdit    bool   `json:"can_edit"`
}

// SharedFolderUserTeamEntry is one element of the users array from get classic shared folder --format json.
type SharedFolderUserTeamEntry struct {
	// User specific fields
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Owner    bool   `json:"owner"`

	// Team specific fields
	TeamName string `json:"name"`
	TeamId   string `json:"team_uid"`

	// Common for both users and teams
	ManageUsers   bool `json:"manage_users"`
	ManageRecords bool `json:"manage_records"`
}

// SharedFolderResponse is the data payload from get SHARED_FOLDER_ID --format json.
type SharedFolderResponse struct {
	FolderUID      string                       `json:"folder_uid"`
	Name           string                       `json:"name"`
	Path           string                       `json:"path"`
	ManageUsers    bool                         `json:"manage_users"`
	ManageRecords  bool                         `json:"manage_records"`
	CanShare       bool                         `json:"can_share"`
	CanEdit        bool                         `json:"can_edit"`
	Records        []SharedFolderRecordEntry    `json:"records"`
	Users          []SharedFolderUserTeamEntry  `json:"users"`
	Teams          []SharedFolderUserTeamEntry  `json:"teams,omitempty"`
	FolderLocation utils.FolderLocationResponse `json:"folder"`
}

// RecordEntryMapElemType is the object type for each entry in the records map.
var RecordEntryMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		AttrCanShare: types.BoolType,
		AttrCanEdit:  types.BoolType,
	},
}

// UserEntryMapElemType is the object type for each entry in the users map.
var UserEntryMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		AttrManageUsers:   types.BoolType,
		AttrManageRecords: types.BoolType,
	},
}
