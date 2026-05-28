// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package non_shared_folder

// errOpGet is the short description passed to ExecuteCommand for `get`
// failures (logs / error context).
const errOpGet = "Unable to get non-shared folder"

// API response keys (get FOLDER_UID --format json). KeyFolderUID lives in
// folderutils since it is shared with the classic shared folder create response.
const (
	KeyName = "name"
	KeyPath = "path"
)
