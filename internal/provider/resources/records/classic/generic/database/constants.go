// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

const (
	SchemaDescription         = "Creates and manages a Database record in the vault."
	SchemaMarkdownDescription = "Creates and manages a **Database** record in the vault."

	ErrSummaryCreateFailed = "Database Record Create Failed"
	ErrSummaryReadFailed   = "Database Record Read Failed"
	ErrSummaryUpdateFailed = "Database Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the Database record."
	ErrDetailReadFailed   = "Something went wrong when reading the Database record."
	ErrDetailUpdateFailed = "Something went wrong when updating the Database record."
)
