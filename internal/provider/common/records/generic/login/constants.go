// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

const (
	FlagLogin    = "login"
	FlagPassword = "password"
	FlagURL      = "url"
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

	LoginDescription         = "Username or login identifier."
	LoginMarkdownDescription = "Username or login identifier. Maps to the record's `login` field."

	PasswordDescription         = "Password for the login."
	PasswordMarkdownDescription = "Password for the login. Maps to the record's `password` field."

	WebsiteAddressDescription         = "Website address for the login."
	WebsiteAddressMarkdownDescription = "Website address for the login. Maps to the record's `url` field."

	DSNotesDescription         = "Notes on the record, if any."
	DSNotesMarkdownDescription = "**Notes** on the record, if any."

	DSFolderDescription         = "Folder path where the record is stored."
	DSFolderMarkdownDescription = "**Folder path** where the record is stored."
)
