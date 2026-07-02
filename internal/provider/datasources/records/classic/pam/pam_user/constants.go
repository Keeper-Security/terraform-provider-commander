// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

const (
	SchemaDescription         = "Use this data source to look up a classic PAM user record by UID."
	SchemaMarkdownDescription = "Use this data source to look up a **classic PAM user** record by **UID**."

	ErrSummaryReadFailed        = "Failed to read PAM User data source"
	ErrDetailReadFailed         = "Unable to read PAM User vault record"
	ErrDetailRotationInfoFailed = "Unable to read PAM User rotation info"
)
