// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

const (
	FlagSSID         = "text.SSID"
	FlagPassword     = "password"
	FlagEncryption   = "wifiEncryption"
	FlagIsSSIDHidden = "isSSIDHidden"
)

const (
	IDDescription         = "The unique identifier (UID) of the vault record."
	IDMarkdownDescription = "The unique identifier (**UID**) of the vault record."

	TitleDescription         = "Record title."
	TitleMarkdownDescription = "Record title."

	NotesDescription         = "Manage note for the record."
	NotesMarkdownDescription = "Manage note for the record."

	FolderDescription         = "Folder path or UID where the record is to be stored."
	FolderMarkdownDescription = "Folder `path` or `UID` where the record is to be stored."

	SSIDDescription         = "WiFi network SSID (network name)."
	SSIDMarkdownDescription = "WiFi network SSID (network name). Maps to the record's `text` field with label `SSID`."

	PasswordDescription         = "Password for the WiFi network."
	PasswordMarkdownDescription = "Password for the WiFi network. Maps to the record's `password` field."

	EncryptionDescription         = "Encryption type. One of: wep, wpa, noEncryption."
	EncryptionMarkdownDescription = "Encryption type. One of: `wep`, `wpa`, `noEncryption`. Maps to the record's `wifiEncryption` field."

	IsSSIDHiddenDescription         = "Whether the SSID is hidden (not broadcast)."
	IsSSIDHiddenMarkdownDescription = "Whether the SSID is hidden (not broadcast). Maps to the record's `isSSIDHidden` field."

	DSNotesDescription         = "Notes on the record, if any."
	DSNotesMarkdownDescription = "**Notes** on the record, if any."

	DSFolderDescription         = "Folder path where the record is stored."
	DSFolderMarkdownDescription = "**Folder path** where the record is stored."
)

// AllowedEncryptions lists the supported wifiEncryption values accepted by Keeper.
var AllowedEncryptions = []string{"wep", "wpa", "noEncryption"}
