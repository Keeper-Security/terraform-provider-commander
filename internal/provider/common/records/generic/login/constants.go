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

	CustomDescription                  = "Manage custom fields for the record."
	CustomMarkdownDescription          = "Manage custom fields for the record."
	CustomTypeDescription              = "Keeper field type (e.g. text, email, secret, phone, name, date)."
	CustomTypeMarkdownDescription      = "Keeper field type (e.g. `text`, `email`, `secret`, `phone`, `name`, `date`)."
	CustomLabelDescription             = "Field label."
	CustomLabelMarkdownDescription     = "Field label."
	CustomValueDescription             = "Field value; for complex types use jsonencode(JSON) matching the Keeper field schema."
	CustomValueMarkdownDescription     = "Field value; for complex types use `jsonencode(JSON)` matching the [Keeper field schema](https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types)."
	CustomSensitiveDescription         = "Whether to mark the value as sensitive in Terraform state display."
	CustomSensitiveMarkdownDescription = "Whether to mark the value as sensitive in Terraform state display."

	DSNotesDescription         = "Notes on the record, if any."
	DSNotesMarkdownDescription = "**Notes** on the record, if any."

	DSFolderDescription         = "Folder path where the record is stored."
	DSFolderMarkdownDescription = "**Folder path** where the record is stored."

	DSCustomDescription                  = "Custom fields stored in the record's `custom` array."
	DSCustomMarkdownDescription          = "Custom fields stored in the record's `custom` array."
	DSCustomTypeDescription              = "Keeper field type (e.g. text, email, secret)."
	DSCustomTypeMarkdownDescription      = "Keeper field **type** (e.g. `text`, `email`, `secret`)."
	DSCustomValueDescription             = "Field value (JSON-encoded for complex types)."
	DSCustomValueMarkdownDescription     = "Field **value** (JSON-encoded for complex types)."
	DSCustomSensitiveDescription         = "Whether the value should be treated as sensitive."
	DSCustomSensitiveMarkdownDescription = "Whether the value should be treated as **sensitive**."
)
