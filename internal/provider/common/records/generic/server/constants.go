// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

const (
	FlagLogin    = "login"
	FlagPassword = "password"
	FlagHost     = "f.host"
)

const (
	LoginDescription         = "Server login or username."
	LoginMarkdownDescription = "Server login or username. Maps to the record's `login` field."

	PasswordDescription         = "Server password."
	PasswordMarkdownDescription = "Server password. Maps to the record's `password` field."

	HostnameDescription         = "Server host name or IP address."
	HostnameMarkdownDescription = "Server host name or IP address. Maps to `host.hostName` in the record's `host` field."

	PortDescription         = "Server port."
	PortMarkdownDescription = "Server port. Maps to `host.port` in the record's `host` field."

	DSNotesDescription         = "Notes on the record, if any."
	DSNotesMarkdownDescription = "**Notes** on the record, if any."

	DSFolderDescription         = "Folder path where the record is stored."
	DSFolderMarkdownDescription = "**Folder path** where the record is stored."
)
