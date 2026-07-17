// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

const (
	SchemaDescription         = "Creates and manages a Keeper Bank Account record in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper `bankAccount` (Bank Account) record in the vault."

	ErrSummaryCreateFailed = "Bank Account Record Create Failed"
	ErrSummaryReadFailed   = "Bank Account Record Read Failed"
	ErrSummaryUpdateFailed = "Bank Account Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the Bank Account record."
	ErrDetailReadFailed   = "Something went wrong when reading the Bank Account record."
	ErrDetailUpdateFailed = "Something went wrong when updating the Bank Account record."
)
