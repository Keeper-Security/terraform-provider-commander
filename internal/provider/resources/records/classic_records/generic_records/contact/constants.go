// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

// Commander CLI record field keys for contact (pam_records-style: explicit f.text.* / f.* paths).
const (
	FlagName        = "f.name"
	FlagTextCompany = "f.text.company"
	FlagEmail       = "f.email"
	FlagAddressRef  = "f.addressRef"
	// FlagPhonePrefix — per-type phone slots: phone.Mobile, phone.Home, ...
	FlagPhonePrefix = "phone."
)

const (
	SchemaDescription         = "Creates and manages a Keeper Contact record in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper `contact` record in the vault."

	ErrSummaryCreateFailed = "Contact Record Create Failed"
	ErrSummaryReadFailed   = "Contact Record Read Failed"
	ErrSummaryUpdateFailed = "Contact Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the Contact record."
	ErrDetailReadFailed   = "Something went wrong when reading the Contact record."
	ErrDetailUpdateFailed = "Something went wrong when updating the Contact record."
)
