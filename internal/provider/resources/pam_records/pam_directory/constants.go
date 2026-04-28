// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

const (
	FlagUseSSL         = "f.checkbox.useSSL"
	FlagDomainName     = "f.text.domainName"
	FlagAlternativeIPs = "f.multiline.alternativeIPs"
	FlagDirectoryId    = "f.text.directoryId"
	FlagDirectoryType  = "directoryType"
	FlagUserMatch      = "f.text.userMatch"
	FlagProviderGroup  = "f.text.providerGroup"
	FlagProviderRegion = "f.text.providerRegion"
)

const (
	ErrSummaryAddPamDirectoryRecordFailed    = "Failed to add PAM directory record"
	ErrSummaryPamDirectoryRecordUpdateFailed = "Failed to update PAM directory record"
	ErrSummaryPamDirectoryReadFailed         = "Failed to read PAM directory record"
)

const (
	ErrDetailAddPamDirectoryRecordFailed    = "Unable to add PAM directory record"
	ErrDetailPamDirectoryRecordUpdateFailed = "Unable to update PAM directory record"
	ErrDetailPamDirectoryReadFailed         = "Unable to read PAM directory record"
)
