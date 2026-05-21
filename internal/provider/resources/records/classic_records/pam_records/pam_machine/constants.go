// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

// Commander CLI record field flags for pamMachine record type.
const (
	FlagOperatingSystem = "f.text.operatingSystem"
	FlagInstanceName    = "f.text.instanceName"
	FlagInstanceId      = "f.text.instanceId"
	FlagProviderGroup   = "f.text.providerGroup"
	FlagProviderRegion  = "f.text.providerRegion"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryAddPamMachineRecordFailed    = "Failed to add PAM machine record"
	ErrSummaryPamMachineRecordUpdateFailed = "Failed to update PAM machine record"
	ErrSummaryPamMachineReadFailed         = "Failed to read PAM machine record"
)

// Error details operation messages (second argument to ExecuteCommand and AddError; short description for logs).
const (
	ErrDetailAddPamMachineRecordFailed    = "Unable to add PAM machine record"
	ErrDetailPamMachineRecordUpdateFailed = "Unable to update PAM machine record"
	ErrDetailPamMachineReadFailed         = "Unable to read PAM machine record"
)
