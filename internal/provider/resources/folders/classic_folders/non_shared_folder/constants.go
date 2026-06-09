// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

// Commander CLI commands for non-shared folder operations.
const (
	CmdMkdir = "mkdir"
	CmdRndir = "rndir"
	CmdLn    = "ln"
	CmdRm    = "rm"
)

// Command flags.
const (
	FlagName       = "-n"
	FlagUserFolder = "--user-folder"
)

// Schema descriptions used by the non-shared folder resource. Identity (id, name,
// folder_location) descriptions come from folderutils. Generic folder error
// summaries and CRUD operation messages live in
// internal/provider/common/folders/utils/constants.go (folderutils.ErrSummary*
// and folderutils.ErrOp*).
const (
	DescResource = "Manages a vault folder. Use this resource to create and manage a non-shared folder."
	DescRecords  = "Set of record UIDs to link into this folder."
)
