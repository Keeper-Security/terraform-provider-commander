// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

// Terraform schema descriptions (plain and Markdown) for registry / docs.
const (
	IDDescription         = "The PAM machine record UID assigned by Keeper after create."
	IDMarkdownDescription = "The PAM machine record **UID** assigned by Keeper after create."

	TitleDescription         = "Title of the PAM machine record."
	TitleMarkdownDescription = "**Title** of the PAM machine record."

	FolderDescription         = "Folder UID or path to store PAM machine record in your Keeper vault. If not provided, the record will be stored in the root path of vault."
	FolderMarkdownDescription = "Folder **UID** or path to store PAM machine record in your Keeper vault. If not provided, the record will be stored in the root path of vault."

	HostnameOrIPDescription         = "Hostname or IP address with an optional port for the PAM machine."
	HostnameOrIPMarkdownDescription = "**Hostname or IP address** with an optional port for the PAM machine."

	HostNameDescription         = "Hostname or IP address of the PAM machine."
	HostNameMarkdownDescription = "**Hostname or IP address** of the PAM machine."

	PortDescription         = "Administrative port number for the PAM machine connection."
	PortMarkdownDescription = "**Administrative port number** for the PAM machine connection."

	OperatingSystemDescription         = "Operating system of the PAM machine."
	OperatingSystemMarkdownDescription = "**Operating system** of the PAM machine."

	InstanceNameDescription         = "Instance name of the PAM machine."
	InstanceNameMarkdownDescription = "**Instance name** of the PAM machine."

	InstanceIdDescription         = "Instance ID of the PAM machine."
	InstanceIdMarkdownDescription = "**Instance ID** of the PAM machine."

	ProviderGroupDescription         = "Provider group of the PAM machine."
	ProviderGroupMarkdownDescription = "**Provider group** of the PAM machine."

	ProviderRegionDescription         = "Provider region of the PAM machine."
	ProviderRegionMarkdownDescription = "**Provider region** of the PAM machine."

	NotesDescription         = "Optional notes for this PAM machine record."
	NotesMarkdownDescription = "Optional **notes** for this PAM machine record."

	PamSettingsDescription         = "PAM settings for the PAM machine record."
	PamSettingsMarkdownDescription = "PAM **settings** for the PAM machine record."
)
