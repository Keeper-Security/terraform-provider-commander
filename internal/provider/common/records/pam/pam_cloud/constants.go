// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

const (
	IDDescription         = "The PAM cloud record UID assigned by Keeper after create."
	IDMarkdownDescription = "The PAM cloud record **UID** assigned by Keeper after create."

	TitleDescription         = "Title of the PAM cloud record."
	TitleMarkdownDescription = "**Title** of the PAM cloud record."

	FolderDescription         = "Folder UID or path to store PAM cloud record in your Keeper vault. If not provided, the record will be stored in the root path of vault."
	FolderMarkdownDescription = "Folder **UID** or path to store PAM cloud record in your Keeper vault. If not provided, the record will be stored in the root path of vault."

	URLDescription         = "URL of the cloud resource."
	URLMarkdownDescription = "**URL** of the cloud resource."

	AccountIdentifierDescription         = "Account identifier of the cloud resource."
	AccountIdentifierMarkdownDescription = "**Account identifier** of the cloud resource."

	NotesDescription         = "Notes for this PAM cloud record."
	NotesMarkdownDescription = "**Notes** for this PAM cloud record."
)

// Vault field types of the pamCloudResource record type, as returned by
// `get <uid> --format json`.
const (
	FieldTypeURL               = "rbiUrl"
	FieldTypeAccountIdentifier = "accountNumber"
)

// Commander CLI record field flags for pamCloudResource record type.
const (
	FlagURL               = "f.rbiUrl"
	FlagAccountIdentifier = "f.accountNumber"
)
