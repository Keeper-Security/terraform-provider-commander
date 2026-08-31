// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

const (
	SchemaDescription         = "Creates and manages a Classic SSH keys record in the vault."
	SchemaMarkdownDescription = "Creates and manages a Classic `sshKeys` record in the vault."

	ErrSummaryCreateFailed = "Classic SSH Keys Record Create Failed"
	ErrSummaryReadFailed   = "Classic SSH Keys Record Read Failed"
	ErrSummaryUpdateFailed = "Classic SSH Keys Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the Classic SSH keys record."
	ErrDetailReadFailed   = "Something went wrong when reading the Classic SSH keys record."
	ErrDetailUpdateFailed = "Something went wrong when updating the Classic SSH keys record."
)
