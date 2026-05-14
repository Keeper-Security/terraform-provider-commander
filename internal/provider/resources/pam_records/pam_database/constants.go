// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

const (
	FlagUseSSL         = "f.checkbox.useSSL"
	FlagDatabaseId     = "f.text.databaseId"
	FlagDatabaseType   = "databaseType"
	FlagProviderGroup  = "f.text.providerGroup"
	FlagProviderRegion = "f.text.providerRegion"
)

const (
	ErrSummaryAddPamDatabaseRecordFailed    = "Failed to add PAM database record"
	ErrSummaryPamDatabaseRecordUpdateFailed = "Failed to update PAM database record"
	ErrSummaryPamDatabaseReadFailed         = "Failed to read PAM database record"
)

const (
	ErrDetailAddPamDatabaseRecordFailed    = "Unable to add PAM database record"
	ErrDetailPamDatabaseRecordUpdateFailed = "Unable to update PAM database record"
	ErrDetailPamDatabaseReadFailed         = "Unable to read PAM database record"
)
