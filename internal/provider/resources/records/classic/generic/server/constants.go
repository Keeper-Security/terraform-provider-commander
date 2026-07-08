// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

const (
	SchemaDescription         = "Creates and manages a Keeper server credentials record (`serverCredentials`) in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper **server credentials** record (`serverCredentials`) in the vault."

	ErrSummaryCreateFailed = "Server Record Create Failed"
	ErrSummaryReadFailed   = "Server Record Read Failed"
	ErrSummaryUpdateFailed = "Server Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the server credentials record."
	ErrDetailReadFailed   = "Something went wrong when reading the server credentials record."
	ErrDetailUpdateFailed = "Something went wrong when updating the server credentials record."
)
