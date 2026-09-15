// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

const (
	SchemaDescription         = "Creates and manages a SSH keys record in the vault."
	SchemaMarkdownDescription = "Creates and manages a **SSH keys** record in the vault."

	ErrSummaryCreateFailed = "SSH Keys Record Create Failed"
	ErrSummaryReadFailed   = "SSH Keys Record Read Failed"
	ErrSummaryUpdateFailed = "SSH Keys Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the SSH keys record."
	ErrDetailReadFailed   = "Something went wrong when reading the SSH keys record."
	ErrDetailUpdateFailed = "Something went wrong when updating the SSH keys record."
)
