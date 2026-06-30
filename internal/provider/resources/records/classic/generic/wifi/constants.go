// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

const (
	SchemaDescription         = "Creates and manages a Keeper WiFi credentials record (`wifiCredentials`) in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper **WiFi credentials** record (`wifiCredentials`) in the vault."

	ErrSummaryCreateFailed = "WiFi Record Create Failed"
	ErrSummaryReadFailed   = "WiFi Record Read Failed"
	ErrSummaryUpdateFailed = "WiFi Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the WiFi credentials record."
	ErrDetailReadFailed   = "Something went wrong when reading the WiFi credentials record."
	ErrDetailUpdateFailed = "Something went wrong when updating the WiFi credentials record."
)
