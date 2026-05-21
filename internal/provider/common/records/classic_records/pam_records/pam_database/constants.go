// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

const (
	IDDescription         = "The PAM database record UID assigned by Keeper after create."
	IDMarkdownDescription = "The PAM database record **UID** assigned by Keeper after create."

	TitleDescription         = "Title of the PAM database record."
	TitleMarkdownDescription = "**Title** of the PAM database record."

	FolderDescription         = "Folder UID or path to store PAM database record in your Keeper vault. If not provided, the record will be stored in the root path of vault."
	FolderMarkdownDescription = "Folder **UID** or path to store PAM database record in your Keeper vault. If not provided, the record will be stored in the root path of vault."

	HostnameOrIPDescription         = "Hostname or IP address with an optional administrative port for the PAM database."
	HostnameOrIPMarkdownDescription = "**Hostname or IP address** with an optional administrative port for the PAM database."

	HostNameDescription         = "Address of the database resource."
	HostNameMarkdownDescription = "**Address of the database resource**."

	PortDescription         = "Administrative port number for the PAM database connection to connect on."
	PortMarkdownDescription = "**Administrative port number** for the PAM database connection to connect on."

	UseSSLDescription         = "Whether to use SSL while connecting to the database resource."
	UseSSLMarkdownDescription = "Whether to use **SSL** while connecting to the database resource."

	DatabaseIdDescription         = "Azure or AWS Resource ID"
	DatabaseIdMarkdownDescription = "**Azure or AWS Resource ID**"

	DatabaseTypeDescription         = "Database type of the PAM database. Must be one of: postgresql, postgresql-flexible, mysql, mysql-flexible, mariadb, mariadb-flexible, mssql, oracle, mongodb."
	DatabaseTypeMarkdownDescription = "**Database type** of the PAM database. Must be one of: `postgresql`, `postgresql-flexible`, `mysql`, `mysql-flexible`, `mariadb`, `mariadb-flexible`, `mssql`, `oracle`, `mongodb`."

	ProviderGroupDescription         = "Azure or AWS Provider Group."
	ProviderGroupMarkdownDescription = "**Azure or AWS Provider Group**."

	ProviderRegionDescription         = "Azure or AWS Provider Region."
	ProviderRegionMarkdownDescription = "**Azure or AWS Provider Region**."

	NotesDescription         = "Notes for this PAM database record."
	NotesMarkdownDescription = "**Notes** for this PAM database record."
)
