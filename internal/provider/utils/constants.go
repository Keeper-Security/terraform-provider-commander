// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

// Common Commander CLI commands
const (
	CmdRecordAdd    = "record-add"
	CmdRecordUpdate = "record-update"
	CmdRecordDelete = "rm"
	CmdGetRecord    = "get"
)

// Commander CLI command flags
const (
	FlagFormatJSON = "--format json"
	FlagForce      = "--force"
	FlagQuiet      = "-q"
	FlagTitle      = "--title"
	FlagNotes      = "--notes"
	FlagRbiUrl     = "rbiUrl"
	FlagFolder     = "--folder"
	FlagRecordType = "--record-type"
	FlagRecord     = "--record"
)

// Commander CLI command flag values
const (
	ValueOn  = "on"
	ValueOff = "off"
)

// Vault Record Types
const (
	RecordTypePamDatabase      = "pamDatabase"
	RecordTypePamDirectory     = "pamDirectory"
	RecordTypePamMachine       = "pamMachine"
	RecordTypePamUser          = "pamUser"
	RecordTypePamRemoteBrowser = "pamRemoteBrowser"
)

// Common schema attribute descriptions.
const (
	EnterpriseManagedCompanySchemaAttributeDescription         = "Only applies to MSP accounts. Name or ID of the managed company to scope this resource or data source to. Omit to use the logged-in account context."
	EnterpriseManagedCompanySchemaAttributeMarkdownDescription = "Only applies to **MSP accounts**. **Name** or **ID** of the managed company to scope this resource or data source to. Omit to use the logged-in account context."
)

// Common Error summaries (first argument to AddError).
const (
	ERR_MSG_PROVIDER_CONFIGURATION_ERROR    = "Provider Configuration Error"
	ERR_MSG_INVALID_IMPORT_ID               = "Invalid Import ID"
	ErrOpListVaultRecords                   = "Unable to list vault records for validation"
	ErrSummaryManagedCompanyCannotBeUpdated = "Managed Company Cannot Be Updated"
	ErrSummarySyncDownFailed                = "Sync Down Failed"
	ErrSummaryRecordDeleteFailed            = "Record Delete Failed"
)

// Error details operation messages (second argument to ExecuteCommand and AddError; short description for logs).
const (
	ErrDetailManagedCompany     = "Cannot update the managed_company field. Once an EPM policy is created, the managed company cannot be changed. Remove and recreate the resource to use a different managed company."
	ErrDetailRecordDeleteFailed = "Something went wrong when deleting the record. Check the record UID or title and try again."
)
