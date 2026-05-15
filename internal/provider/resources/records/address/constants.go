// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

// Commander CLI record field keys for address.
const (
	FlagAddress = "f.address"
)

const (
	SchemaDescription         = "Creates and manages a Keeper Address record in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper `address` record in the vault."

	ErrSummaryCreateFailed = "Address Record Create Failed"
	ErrSummaryReadFailed   = "Address Record Read Failed"
	ErrSummaryUpdateFailed = "Address Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the Address record."
	ErrDetailReadFailed   = "Something went wrong when reading the Address record."
	ErrDetailUpdateFailed = "Something went wrong when updating the Address record."
)
