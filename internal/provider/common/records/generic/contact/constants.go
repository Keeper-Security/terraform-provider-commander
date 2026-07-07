// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

// Commander CLI record field keys for contact.
const (
	FlagName        = "f.name"
	FlagTextCompany = "f.text.company"
	FlagEmail       = "f.email"
	FlagAddressRef  = "f.addressRef"
	// FlagPhonePrefix — per-type phone slots: phone.Mobile, phone.Home, ...
	FlagPhonePrefix = "phone."
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

	NameDescription         = "Person name (first, middle, last)."
	NameMarkdownDescription = "Person name (`name` field): first, middle, last."

	FirstNameDescription         = "First name."
	FirstNameMarkdownDescription = "First name."

	MiddleNameDescription         = "Middle name."
	MiddleNameMarkdownDescription = "Middle name."

	LastNameDescription         = "Last name."
	LastNameMarkdownDescription = "Last name."

	CompanyDescription         = "Company name."
	CompanyMarkdownDescription = "Company name."

	EmailDescription         = "Email address."
	EmailMarkdownDescription = "Email address."

	PhoneDescription         = "Manage phone numbers for the record."
	PhoneMarkdownDescription = "Manage phone numbers for the record."

	PhoneRegionDescription         = "Region or country code (e.g. US, +1)."
	PhoneRegionMarkdownDescription = "Region or country code."

	PhoneNumberDescription         = "Phone number."
	PhoneNumberMarkdownDescription = "Phone number."

	PhoneExtDescription         = "Extension."
	PhoneExtMarkdownDescription = "Extension."

	PhoneTypeDescription         = "Phone type: Mobile, Home, or Work."
	PhoneTypeMarkdownDescription = "Phone type: `Mobile`, `Home`, or `Work`."

	AddressRefDescription         = "Linked Address record UID."
	AddressRefMarkdownDescription = "UID of an `address` record linked via `addressRef`."

	CustomDescription         = "Manage custom fields for the record."
	CustomMarkdownDescription = "Manage custom fields for the record."

	DSNotesDescription         = "Notes on the record."
	DSNotesMarkdownDescription = "Notes on the record."

	DSFolderDescription         = "Folder path or UID where the record is stored."
	DSFolderMarkdownDescription = "Folder path or UID where the record is stored."

	DSPhoneDescription         = "Phone numbers."
	DSPhoneMarkdownDescription = "Phone numbers."

	DSCustomDescription         = "Custom fields on the record."
	DSCustomMarkdownDescription = "Custom fields on the record."

	DSCustomTypeDescription         = "Keeper field type."
	DSCustomTypeMarkdownDescription = "Keeper field type."

	DSCustomValueDescription         = "Field value."
	DSCustomValueMarkdownDescription = "Field value."

	DSCustomSensitiveDescription         = "Whether the value is sensitive."
	DSCustomSensitiveMarkdownDescription = "Whether the value is sensitive."
)
