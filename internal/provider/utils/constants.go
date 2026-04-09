// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

// Common schema attribute descriptions.
const (
	EnterpriseManagedCompanySchemaAttributeDescription         = "Only applies to MSP accounts. Name or ID of the managed company to scope this resource or data source to. Omit to use the logged-in account context."
	EnterpriseManagedCompanySchemaAttributeMarkdownDescription = "Only applies to **MSP accounts**. **Name** or **ID** of the managed company to scope this resource or data source to. Omit to use the logged-in account context."
)

// Common error messages.
const (
	ERR_MSG_PROVIDER_CONFIGURATION_ERROR = "Provider Configuration Error"
	ERR_MSG_INVALID_IMPORT_ID            = "Invalid Import ID"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryManagedCompanyCannotBeUpdated = "Managed Company Cannot Be Updated"
)

// Error details.
const (
	ErrDetailManagedCompany = "Cannot update the managed_company field. Once an EPM policy is created, the managed company cannot be changed. Remove and recreate the resource to use a different managed company."
)

// Flag value literals (on/off).
const (
	ValueOn  = "on"
	ValueOff = "off"
)
