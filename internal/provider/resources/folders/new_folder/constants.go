// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

// Commander CLI commands used only by this resource.
const (
	CmdNsfMkdir = "nsf-mkdir"
	CmdNsfRmdir = "nsf-rmdir"
	CmdNsfRndir = "nsf-rndir"
	CmdNsfLn    = "nsf-ln"
	CmdNsfRm    = "nsf-rm"
)

// Schema descriptions for the new folder resource. Generic folder error
// summaries and CRUD operation messages live in
// internal/provider/common/folders/utils/constants.go (folderutils.ErrSummary*
// and folderutils.ErrOp*).
const (
	DescResource = "Manages a Nested Shared Folder."
	DescRecords  = "Set of record UIDs to link into this folder."
)

const (
	ErrSummaryNsfMoveFolderFailed = "Nsf Move Folder Failed"
)
