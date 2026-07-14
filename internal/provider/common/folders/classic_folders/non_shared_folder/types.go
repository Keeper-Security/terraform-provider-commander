// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package non_shared_folder

import "github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"

// NonSharedFolderResponse is the data payload from `get FOLDER_UID --format json` for a non-shared folder.
type NonSharedFolderResponse struct {
	FolderUID      string                          `json:"folder_uid"`
	Name           string                          `json:"name"`
	Type           string                          `json:"type"`
	Path           string                          `json:"path,omitempty"`
	Records        []NonSharedFolderRecordResponse `json:"records,omitempty"`
	FolderLocation utils.FolderLocationResponse    `json:"folder"`
}

// NonSharedFolderRecordResponse represents a record entry in the folder get response.
type NonSharedFolderRecordResponse struct {
	RecordUID  string `json:"record_uid"`
	RecordName string `json:"record_name"`
}
