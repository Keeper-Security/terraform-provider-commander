// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

const (
	IDDescription         = "The PAM directory record UID assigned by Keeper after create."
	IDMarkdownDescription = "The PAM directory record **UID** assigned by Keeper after create."

	TitleDescription         = "Title of the PAM directory record."
	TitleMarkdownDescription = "**Title** of the PAM directory record."

	FolderDescription         = "Folder UID or path to store PAM directory record in your Keeper vault. If not provided, the record will be stored in the root path of vault."
	FolderMarkdownDescription = "Folder **UID** or path to store PAM directory record in your Keeper vault. If not provided, the record will be stored in the root path of vault."

	HostnameOrIPDescription         = "Hostname or IP address with an optional administrative port for the PAM directory."
	HostnameOrIPMarkdownDescription = "**Hostname or IP address** with an optional administrative port for the PAM directory."

	HostNameDescription         = "Address of the directory resource."
	HostNameMarkdownDescription = "**Address of the directory resource**."

	PortDescription         = "Administrative port number for the PAM directory connection to connect on."
	PortMarkdownDescription = "**Administrative port number** for the PAM directory connection to connect on."

	UseSSLDescription         = "Whether to use SSL while connecting to the directory resource."
	UseSSLMarkdownDescription = "Whether to use **SSL** while connecting to the directory resource."

	DomainNameDescription         = "Domain managed by the directory."
	DomainNameMarkdownDescription = "**Domain managed by the directory**."

	AlternativeIPsDescription         = "List of failover IPs for the directory, used for Discovery. Provide one IP per entry as an array of strings."
	AlternativeIPsMarkdownDescription = "**List of failover IPs** for the directory, used for Discovery. Provide one IP per entry as an array of strings."

	DirectoryIdDescription         = "Instance ID for AD resource in Azure and AWS hosted environments"
	DirectoryIdMarkdownDescription = "**Instance ID** for AD resource in Azure and AWS hosted environments"

	DirectoryTypeDescription         = "Directory type, used for formatting of messaging. Must be one of: active_directory, openldap."
	DirectoryTypeMarkdownDescription = "**Directory type**, used for formatting of messaging. Must be one of: `active_directory`, `openldap`."

	UserMatchDescription         = "Match on OU to filter found users during Discovery. Either match the right side of the DN or surround with slashes for a regular expression. Example: OU=Users,DC=company,DC=com  or /OU=Users/"
	UserMatchMarkdownDescription = "**Match on OU** to filter found users during Discovery. Either match the right side of the DN or surround with slashes for a regular expression. Example: `OU=Users,DC=company,DC=com` or `/OU=Users/`"

	ProviderGroupDescription         = "Provider Group for directories hosted in Azure."
	ProviderGroupMarkdownDescription = "**Provider Group** for directories hosted in Azure."

	ProviderRegionDescription         = "AWS region of hosted directory."
	ProviderRegionMarkdownDescription = "**AWS region** of hosted directory."

	NotesDescription         = "Notes for this PAM directory record."
	NotesMarkdownDescription = "**Notes** for this PAM directory record."
)

// Commander CLI record field flags for pamDirectory record type.
const (
	FlagUseSSL         = "f.checkbox.useSSL"
	FlagDomainName     = "f.text.domainName"
	FlagAlternativeIPs = "f.multiline.alternativeIPs"
	FlagDirectoryId    = "f.text.directoryId"
	FlagDirectoryType  = "directoryType"
	FlagUserMatch      = "f.text.userMatch"
	FlagProviderGroup  = "f.text.providerGroup"
	FlagProviderRegion = "f.text.providerRegion"
)
