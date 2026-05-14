// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

// Commander CLI commands for folder operations.
const (
	CmdMkdir = "mkdir"
	CmdRndir = "rndir"
	CmdRmdir = "rmdir"
	CmdMv    = "mv"
	CmdGet   = "get"
	CmdLn    = "ln"
	CmdRm    = "rm"
)

// Command flags.
const (
	FlagName       = "-n"
	FlagColor      = "--color"
	FlagFormat     = "--format"
	FlagForce      = "-f"
	FlagQuiet      = "-q"
	FormatJSON     = "json"
	FlagUserFolder = "--user-folder"
)

// Valid folder color values.
const (
	ColorNone   = "none"
	ColorRed    = "red"
	ColorGreen  = "green"
	ColorBlue   = "blue"
	ColorOrange = "orange"
	ColorYellow = "yellow"
	ColorGray   = "gray"
)

// ValidColors is the set of accepted --color values.
var ValidColors = []string{ColorNone, ColorRed, ColorGreen, ColorBlue, ColorOrange, ColorYellow, ColorGray}

// API response keys.
const (
	KeyFolderUID = "folder_uid"
	KeyName      = "name"
	KeyPath      = "path"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryCreateFailed   = "Create Folder Failed"
	ErrSummaryReadFailed     = "Read Folder Failed"
	ErrSummaryUpdateFailed   = "Update Folder Failed"
	ErrSummaryDeleteFailed   = "Delete Folder Failed"
	ErrSummarySyncDownFailed = "Sync Down Failed"
	ErrSummaryInvalidConfig  = "Invalid Folder Configuration"
)

// Error operation messages (second argument to ExecuteCommand).
const (
	ErrOpCreateFolder = "Unable to create folder"
	ErrOpRenameFolder = "Unable to rename folder"
	ErrOpMoveFolder   = "Unable to move folder"
	ErrOpDeleteFolder = "Unable to delete folder"
	ErrOpGetFolder    = "Unable to get folder"
	ErrOpLinkRecord   = "Unable to link record to folder"
	ErrOpUnlinkRecord = "Unable to remove record from folder"
)

// Schema and validator descriptions.
const (
	DescResource       = "Manages a vault folder. Use this resource to create and manage a non-shared folder."
	DescId             = "The UID of the folder."
	DescName           = "Folder name (leaf name, without parent path)."
	DescFolderLocation = "Parent folder path where the folder will be created. Leave empty for vault root."
	DescColor          = "Folder color. Valid values: none, red, green, blue, orange, yellow, gray."
	DescRecords        = "Set of record UIDs or titles to link into this folder."
)
