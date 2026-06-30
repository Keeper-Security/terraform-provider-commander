// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

// Commander CLI commands for non-shared folder operations.
const (
	CmdRmdir = "rmdir"
)

// API response keys returned by Commander folder commands (mkdir/get/...).
const (
	KeyFolderUID = "folder_uid"
)

// Command flags.
const (
	FlagName = "--name"
)

// Terraform attribute names for folder identity fields.
const (
	AttrId             = "id"
	AttrName           = "name"
	AttrFolderLocation = "folder_location"
)

// Schema descriptions for folder id, name and folder_location (resource and data source).
const (
	DescId             = "The folder ID assigned by Keeper."
	DescIdMD           = "The folder **ID** assigned by Keeper."
	DescName           = "Folder name."
	DescNameMD         = "**Folder name**."
	DescFolderLocation = "Parent folder path where the folder will be created. Leave empty for vault root."

	// NameValidatorLabel is the human-readable field name passed to StringMinLengthValidator.
	NameValidatorLabel = "Folder Name"
)

// Generic folder error summaries (first argument to AddError). Shared across all
// folder resources and data sources (classic shared folder, classic non-shared
// folder, Keeper Drive new folder, etc.).
const (
	ErrSummaryCreateFailed  = "Create Folder Failed"
	ErrSummaryReadFailed    = "Read Folder Failed"
	ErrSummaryUpdateFailed  = "Update Folder Failed"
	ErrSummaryDeleteFailed  = "Delete Folder Failed"
	ErrSummaryInvalidConfig = "Invalid Folder Configuration"
)

// Generic folder error operation messages (passed to ExecuteCommand and
// AddError; short description for logs). Shared across all folder resources
// and data sources.
const (
	ErrOpCreate       = "Unable to create folder"
	ErrOpRead         = "Unable to read folder"
	ErrOpUpdate       = "Unable to update folder"
	ErrOpRename       = "Unable to rename folder"
	ErrOpMove         = "Unable to move folder"
	ErrOpDelete       = "Unable to delete folder"
	ErrOpGet          = "Unable to get folder"
	ErrOpLinkRecord   = "Unable to link record to folder"
	ErrOpUnlinkRecord = "Unable to remove record from folder"
)
