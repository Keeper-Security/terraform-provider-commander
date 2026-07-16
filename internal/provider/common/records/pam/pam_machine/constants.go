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

	HostNameDescription         = "Address of the machine resource."
	HostNameMarkdownDescription = "**Address of the machine resource**."

	PortDescription         = "Port to connect on. The Gateway uses this to determine connection method."
	PortMarkdownDescription = "**Port to connect on. The Gateway uses this to determine connection method.**"

	OperatingSystemDescription         = "The target's Operating System"
	OperatingSystemMarkdownDescription = "**The target's Operating System**"

	InstanceNameDescription         = "Azure or AWS Instance Name"
	InstanceNameMarkdownDescription = "**Azure or AWS Instance Name**"

	InstanceIdDescription         = "Azure or AWS Instance ID"
	InstanceIdMarkdownDescription = "**Azure or AWS Instance ID**"

	ProviderGroupDescription         = "Provider Group for directories hosted in Azure."
	ProviderGroupMarkdownDescription = "**Provider group** of the PAM machine."

	ProviderRegionDescription         = "AWS region of hosted directory."
	ProviderRegionMarkdownDescription = "**AWS region** of hosted directory."

	NotesDescription         = "Notes for this PAM machine record."
	NotesMarkdownDescription = "**Notes** for this PAM machine record."

	PamSettingsDescription         = "This is where you configure Connection and Tunnel settings for this machine."
	PamSettingsMarkdownDescription = "This is where you configure **Connection and Tunnel settings** for this machine."
)

// Commander CLI record field flags for pamMachine record type.
const (
	FlagOperatingSystem = "f.text.operatingSystem"
	FlagInstanceName    = "f.text.instanceName"
	FlagInstanceId      = "f.text.instanceId"
	FlagProviderGroup   = "f.text.providerGroup"
	FlagProviderRegion  = "f.text.providerRegion"
)
