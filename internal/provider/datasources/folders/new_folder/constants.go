// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

// Schema descriptions specific to the Keeper Drive (new) folder data source.
const (
	// Data source attribute name for the lookup field. The value is either
	// the folder UID or the folder name; nsf-get accepts both.
	AttrNewFolder = "new_folder"

	DescDataSource            = "Look up an existing Keeper Drive folder by UID or name."
	DescDataSourceMD          = "Look up an existing Keeper Drive folder by **UID** or **name**."
	DescDataSourceNewFolder   = "Nested shared folder UID or name to look up."
	DescDataSourceNewFolderMD = "Nested shared folder **UID** or **name** to look up."
	DescDataSourceRecords     = "Set of record UIDs linked to this folder."
)
