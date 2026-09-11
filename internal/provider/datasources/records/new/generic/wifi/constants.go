// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

const (
	SchemaDescription         = "Retrieves a New (NSF) WiFi credentials record from the vault."
	SchemaMarkdownDescription = "Retrieves a New (NSF) **WiFi credentials** record from the vault.\n\nA WiFi credentials record stores network credentials (SSID, password, encryption type) used to connect to a wireless network."

	ErrSummaryReadFailed = "Failed to read New (NSF) WiFi data source"
	ErrDetailReadFailed  = "Unable to read New (NSF) WiFi vault record"
)
