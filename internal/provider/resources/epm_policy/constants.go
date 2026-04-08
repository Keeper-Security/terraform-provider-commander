// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

// Error summaries (first argument to AddError).
const (
	ErrSummaryCreateFailed    = "Create EPM Policy Failed"
	ErrSummaryUpdateFailed    = "Update EPM Policy Failed"
	ErrSummaryDeleteFailed    = "Delete EPM Policy Failed"
	ErrSummaryInvalidImportID = "Invalid Import ID"
	ErrSummaryReadFailed      = "Read EPM Policy Failed"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpCreateEpmPolicy = "Unable to create EPM policy"
	ErrOpUpdateEpmPolicy = "Unable to update EPM policy"
	ErrOpDeleteEpmPolicy = "Unable to delete EPM policy"
	ErrOpReadEpmPolicy   = "Unable to read EPM policy"
)

// Import resource type name (used by ParseManagedCompanyImportID).
const (
	ImportResourceType = "epm_policy"
)
