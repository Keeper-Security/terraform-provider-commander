// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

const (
	SchemaDescription         = "Creates and manages a Keeper SaaS configuration record (`saasConfiguration`) in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper **SaaS configuration** record (`saasConfiguration`) in the vault."

	ErrSummaryCreateFailed = "SaaS Configuration Record Create Failed"
	ErrSummaryReadFailed   = "SaaS Configuration Record Read Failed"
	ErrSummaryUpdateFailed = "SaaS Configuration Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the SaaS configuration record."
	ErrDetailReadFailed   = "Something went wrong when reading the SaaS configuration record."
	ErrDetailUpdateFailed = "Something went wrong when updating the SaaS configuration record."
)
