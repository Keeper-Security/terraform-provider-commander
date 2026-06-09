// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_folder

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

// NewFolderGetResponse is the nsf-get JSON payload for a Keeper Drive folder.
// Embeds new_share.ShareResponseFragment so the user_permissions array can be
// parsed by the new_share package's MapResponseToModel helper.
type NewFolderGetResponse struct {
	FolderUID      string                       `json:"folder_uid"`
	Name           string                       `json:"name"`
	FolderLocation utils.FolderLocationResponse `json:"folder"`
	Records        []NewFolderRecordResponse    `json:"records,omitempty"`
	new_share.ShareResponseFragment
}

// NewFolderRecordResponse represents a record entry in the folder get response.
type NewFolderRecordResponse struct {
	RecordUID  string `json:"record_uid"`
	RecordName string `json:"record_name"`
}
