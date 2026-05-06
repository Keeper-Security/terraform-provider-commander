// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

const (
	SchemaDescription         = "Retrieves a Keeper PAM User record from the vault."
	SchemaMarkdownDescription = "Retrieves a Keeper **PAM User** record (`pamUser`) from the vault.\n\nA PAM User record stores privileged credentials (login/password) that can be associated with PAM Machines for rotation, connections, and tunneling."

	ErrSummaryReadFailed        = "Failed to read PAM User data source"
	ErrDetailReadFailed         = "Unable to read PAM User vault record"
	ErrDetailRotationInfoFailed = "Unable to read PAM User rotation info"
)
