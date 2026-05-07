// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

// Commander CLI commands for enterprise SCIM operations.
const (
	CmdScimCreate = "scim create"
	CmdScimEdit   = "scim edit"
	CmdScimDelete = "scim delete"
)

// Command flags for scim create/edit/delete.
const (
	FlagNode         = "--node"
	FlagPrefix       = "--prefix"
	FlagUniqueGroups = "--unique-groups"
	FlagForce        = "--force"
)

// Flag value literals for unique-groups (on/off).
const (
	ValueOn  = "on"
	ValueOff = "off"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryCreateFailed                  = "Create Enterprise SCIM Failed"
	ErrSummaryReadFailed                    = "Read Enterprise SCIM Failed"
	ErrSummaryUpdateFailed                  = "Update Enterprise SCIM Failed"
	ErrSummaryDeleteFailed                  = "Delete Enterprise SCIM Failed"
	ErrSummaryInvalidImportID               = "Invalid Import ID"
	ErrSummaryManagedCompanyCannotBeUpdated = "Managed Company Cannot Be Updated"
	ErrSummaryNodeCannotBeUpdated           = "Node Cannot Be Updated"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpCreateScim = "Unable to create enterprise SCIM"
	ErrOpDeleteScim = "Unable to delete enterprise SCIM"
	ErrOpUpdateScim = "Unable to update enterprise SCIM"
)

// Import and validation error details.
const (
	ErrMsgImportIDEmpty     = "Import ID cannot be empty. Use the SCIM ID (e.g. 1169425105420640) or \"managed_company_name_or_id,scim_id\"."
	ErrDetailManagedCompany = "Cannot update the managed_company field. Once an enterprise SCIM is created, the managed company cannot be changed. Remove and recreate the resource to use a different managed company."
	ErrDetailNode           = "Cannot update the node field. Once an enterprise SCIM is created."
)

// Import resource type name (used by ParseManagedCompanyImportID).
const (
	ImportResourceType = "scim"
)
