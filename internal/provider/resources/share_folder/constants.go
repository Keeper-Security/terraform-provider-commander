// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder

// Commander CLI commands for shared folder operations.
const (
	CmdGet         = "get"
	CmdMkdir       = "mkdir"
	CmdRndir       = "rndir"
	CmdRmdir       = "rmdir"
	CmdShareFolder = "share-folder"
)

// Command flags (shared across mkdir, rndir, rmdir, share-folder).
const (
	FlagSharedFolder  = "--shared-folder"
	FlagAction        = "--action"
	FlagName          = "--name"
	FlagEmail         = "--email"
	FlagRecord        = "--record"
	FlagManageUsers   = "--manage-users"
	FlagManageRecords = "--manage-records"
	FlagCanShare      = "--can-share"
	FlagCanEdit       = "--can-edit"
	FlagExpireAt      = "--expire-at"
	FlagExpireIn      = "--expire-in"
	FlagForce         = "-f"
	FlagQuiet         = "-q"
	FlagFormat        = "--format"
)

// Format values for get command.
const (
	FormatJSON = "json"
)

// share-folder --action values.
const (
	ActionGrant  = "grant"
	ActionRemove = "remove"
)

// Flag value literals (on/off, never, wildcard for default permissions).
const (
	ValueOn     = "on"
	ValueOff    = "off"
	ValueNever  = "never"
	WildcardAll = "*"
)

// API response keys (e.g. from mkdir/create response).
const (
	KeyFolderUID = "folder_uid"
)

// API response keys for get shared folder response.
const (
	KeySharedFolderUID = "shared_folder_uid"
	KeyName            = "name"
	KeyRecords         = "records"
	KeyUsers           = "users"
	KeyRecordUID       = "record_uid"
	KeyUsername        = "username"
)

// Object attribute names (match tfsdk/schema: record and user permission entries).
const (
	AttrCanShare      = "can_share"
	AttrCanEdit       = "can_edit"
	AttrManageUsers   = "manage_users"
	AttrManageRecords = "manage_records"
	AttrExpiration    = "expiration"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryProviderConfig = "Provider Configuration Error"
	ErrSummaryCreateFailed   = "Create Shared Folder Failed"
	ErrSummaryReadFailed     = "Read Shared Folder Failed"
	ErrSummaryUpdateFailed   = "Update Shared Folder Failed"
	ErrSummaryDeleteFailed   = "Delete Shared Folder Failed"
	ErrSummarySyncDownFailed = "Sync Down Failed"
	ErrSummaryInvalidConfig  = "Invalid Shared Folder Configuration"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpCreateSF           = "Unable to create shared folder"
	ErrOpRenameSF           = "Unable to rename shared folder"
	ErrOpUpdateDefaultPerms = "Unable to update shared folder default permissions"
	ErrOpRemoveRecord       = "Unable to remove record from shared folder"
	ErrOpAddUpdateRecord    = "Unable to add/update record in shared folder"
	ErrOpRemoveUser         = "Unable to remove user from shared folder"
	ErrOpAddUpdateUser      = "Unable to add/update user in shared folder"
	ErrOpDeleteSF           = "Unable to delete shared folder"
	ErrOpGetSF              = "Unable to get shared folder"
)

// Time layouts for expiration validation (Go reference time: 2006-01-02 15:04:05).
const (
	TimeLayoutDate          = "2006-01-02"
	TimeLayoutDateTime      = "2006-01-02 15:04:05"
	TimeLayoutDateTimeShort = "2006-01-02 15:04"
)

// Schema and validator descriptions (resource and attribute level).
const (
	DescResource                   = "Manages a shared folder. Use this resource to create and manage shared folders."
	DescId                         = "The ID of the shared folder."
	DescName                       = "Shared folder name."
	DescFolderLocation             = "Folder path or identifier where the shared folder is located."
	DescUserPermissions            = "Default user permissions for the shared folder. When omitted, defaults to manage_users = false, manage_records = false. Allowed keys: manage_users, manage_records."
	DescUserPermissionsMD          = "Default user permissions for the shared folder. When omitted, defaults to `manage_users = false`, `manage_records = false`. Allowed keys: `manage_users`, `manage_records`."
	DescUserPermissionsManage      = "Allow managing users in the shared folder."
	DescUserPermissionsRecords     = "Allow managing records in the shared folder."
	DescRecordPermissions          = "Default record permissions for the shared folder. When omitted, defaults to can_share = false, can_edit = false. Allowed keys: can_share, can_edit."
	DescRecordPermissionsMD        = "Default record permissions for the shared folder. When omitted, defaults to `can_share = false`, `can_edit = false`. Allowed keys: `can_share`, `can_edit`."
	DescRecordPermissionsShare     = "Allow sharing records."
	DescRecordPermissionsEdit      = "Allow editing records."
	DescRecords                    = "Per-record permissions. Map key is record UID or name; value is an object with can_share and can_edit."
	DescRecordsMD                  = "Per-record permissions. Map key is record UID or name; value is an object with `can_share` and `can_edit`."
	DescRecordShare                = "Allow sharing this record."
	DescRecordEdit                 = "Allow editing this record."
	DescUsers                      = "Per-user permissions. Map key is user email or UID; value is an object with manage_users, manage_records, and expiration (\"never\" | ISO date/datetime | relative e.g. 30d, 1y)."
	DescUsersMD                    = "Per-user permissions. Map key is user email or UID; value is an object with `manage_users`, `manage_records`, and `expiration` (\"never\" | ISO date/datetime | relative e.g. 30d, 1y)."
	DescUserManageUsers            = "Allow this user to manage users."
	DescUserManageRecords          = "Allow this user to manage records."
	DescExpiration                 = "Access expiration: \"never\", ISO date/datetime (yyyy-MM-dd or yyyy-MM-dd HH:mm:ss), or relative period (e.g. 30d, 1y, 6mo, 24h, 90days)."
	ExpirationDoc                  = "Must be \"never\", an ISO date/datetime (yyyy-MM-dd or yyyy-MM-dd HH:mm:ss), or a relative period (e.g. 30d, 1y, 6mo, 24h, 90days, 30minutes)."
	DescUserPermissionsDefault     = "When null, defaults to manage_users = false, manage_records = false."
	DescUserPermissionsDefaultMD   = "When null, defaults to `manage_users = false`, `manage_records = false`."
	DescRecordPermissionsDefault   = "When null, defaults to can_share = false, can_edit = false."
	DescRecordPermissionsDefaultMD = "When null, defaults to `can_share = false`, `can_edit = false`."
)

// Expiration validator error messages.
const (
	ErrMsgInvalidExpiration = "Invalid expiration"
	ErrMsgExpirationEmpty   = "expiration cannot be empty. Use \"never\" for no expiration."
)
