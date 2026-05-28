// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package non_shared_folder

// errOpGet is the short description passed to ExecuteCommand for `get`
// failures (logs / error context).
const errOpGet = "Unable to get non-shared folder"

// API response keys (get FOLDER_UID --format json).
const (
	KeyFolderUID = "folder_uid"
	KeyName      = "name"
	KeyPath      = "path"
)
