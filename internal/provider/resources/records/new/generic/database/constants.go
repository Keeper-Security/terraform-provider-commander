// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

const (
	SchemaDescription         = "Creates and manages a New (NSF) database credentials record in the vault."
	SchemaMarkdownDescription = "Creates and manages a New (NSF) **database credentials** record in the vault."

	ErrSummaryCreateFailed = "New (NSF) Database Record Create Failed"
	ErrSummaryReadFailed   = "New (NSF) Database Record Read Failed"
	ErrSummaryUpdateFailed = "New (NSF) Database Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the New (NSF) database credentials record."
	ErrDetailReadFailed   = "Something went wrong when reading the New (NSF) database credentials record."
	ErrDetailUpdateFailed = "Something went wrong when updating the New (NSF) database credentials record."
)
