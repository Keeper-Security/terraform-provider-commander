// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

const (
	SchemaDescription         = "Creates and manages a Keeper software license record (`softwareLicense`) in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper **software license** record (`softwareLicense`) in the vault."

	ErrSummaryCreateFailed = "Software License Record Create Failed"
	ErrSummaryReadFailed   = "Software License Record Read Failed"
	ErrSummaryUpdateFailed = "Software License Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the software license record."
	ErrDetailReadFailed   = "Something went wrong when reading the software license record."
	ErrDetailUpdateFailed = "Something went wrong when updating the software license record."
)
