// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

const (
	SchemaDescription         = "Retrieves a Keeper WiFi credentials record from the vault."
	SchemaMarkdownDescription = "Retrieves a Keeper **WiFi credentials** record (`wifiCredentials`) from the vault.\n\nA WiFi credentials record stores network credentials (SSID, password, encryption type) used to connect to a wireless network."

	ErrSummaryReadFailed = "Failed to read WiFi data source"
	ErrDetailReadFailed  = "Unable to read WiFi vault record"
)
