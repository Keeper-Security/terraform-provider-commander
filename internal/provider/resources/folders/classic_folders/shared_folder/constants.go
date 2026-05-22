// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/shared_folder"

// Re-export map/Terraform attribute keys and element types for the resource package (shared with common).
const (
	AttrCanShare      = commonsharedfolder.AttrCanShare
	AttrCanEdit       = commonsharedfolder.AttrCanEdit
	AttrManageUsers   = commonsharedfolder.AttrManageUsers
	AttrManageRecords = commonsharedfolder.AttrManageRecords
	AttrExpiration    = commonsharedfolder.AttrExpiration
)

// RecordEntryMapElemType and UserEntryMapElemType match users/records nested map schemas.
var (
	RecordEntryMapElemType = commonsharedfolder.RecordEntryMapElemType
	UserEntryMapElemType   = commonsharedfolder.UserEntryMapElemType
)

// Commander CLI commands for classic shared folder operations.
const (
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

// API response keys for get classic shared folder response.
const (
	KeySharedFolderUID = "shared_folder_uid"
	KeyName            = "name"
	KeyPath            = "path"
	KeyRecords         = "records"
	KeyUsers           = "users"
	KeyRecordUID       = "record_uid"
	KeyUsername        = "username"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryCreateFailed   = "Create Shared Folder Failed"
	ErrSummaryReadFailed     = "Read Shared Folder Failed"
	ErrSummaryUpdateFailed   = "Update Shared Folder Failed"
	ErrSummaryDeleteFailed   = "Delete Shared Folder Failed"
	ErrSummarySyncDownFailed = "Sync Down Failed"
	ErrSummaryInvalidConfig  = "Invalid Shared Folder Configuration"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpCreateSF           = "Unable to create classic shared folder"
	ErrOpRenameSF           = "Unable to rename classic shared folder"
	ErrOpMoveSF             = "Unable to move classic shared folder"
	ErrOpUpdateDefaultPerms = "Unable to update classic shared folder default permissions"
	ErrOpRemoveRecord       = "Unable to remove record from classic shared folder"
	ErrOpAddUpdateRecord    = "Unable to add/update record in classic shared folder"
	ErrOpRemoveUser         = "Unable to remove user from classic shared folder"
	ErrOpAddUpdateUser      = "Unable to add/update user in classic shared folder"
	ErrOpDeleteSF           = "Unable to delete classic shared folder"
	ErrOpGetSF              = "Unable to get classic shared folder"
)

// TimeLayoutExpiration is the expiration datetime format for Terraform config and share-folder --expire-at (yyyy-MM-ddTHH:mm:ss).
const TimeLayoutExpiration = "2006-01-02T15:04:05"

// Schema and validator descriptions (resource and attribute level).
const (
	DescResource                 = "Manages a classic shared folder. Use this resource to create and manage classic shared folder. \n\nClassic shared folder uses classic permission model, Limits sharing to basic access levels. Recommended only for compatibility with older workflows. "
	DescDataSource               = "Look up an existing classic shared folder by UID or name. \n\nClassic shared folder uses classic permission model"
	DescDataSourceMD             = "Look up an existing classic shared folder by **UID** or **name**."
	DescDataSourceSharedFolder   = "Shared folder UID or name to look up."
	DescDataSourceSharedFolderMD = "Shared folder **UID** or **name** to look up."
	DescDataSourceId             = "ID of the found classic shared folder."
	DescDataSourceIdMD           = "**ID** of the found classic shared folder."
	DescId                       = "The ID of the classic shared folder."
	DescName                     = "Shared folder name."
	DescFolderLocation           = "Folder path or identifier where the classic shared folder is located."
	DescUserPermissions          = "Default user permissions for the classic shared folder. When omitted, defaults to manage_users = false, manage_records = false. Allowed keys: manage_users, manage_records."
	DescUserPermissionsMD        = "Default user permissions for the classic shared folder. When omitted, defaults to `manage_users = false`, `manage_records = false`. Allowed keys: `manage_users`, `manage_records`."
	DescUserPermissionsManage    = "Allow managing users in the classic shared folder."
	DescUserPermissionsRecords   = "Allow managing records in the classic shared folder."
	DescRecordPermissions        = "Default record permissions for the classic shared folder. When omitted, defaults to can_share = false, can_edit = false. Allowed keys: can_share, can_edit."
	DescRecordPermissionsMD      = "Default record permissions for the classic shared folder. When omitted, defaults to `can_share = false`, `can_edit = false`. Allowed keys: `can_share`, `can_edit`."
	DescRecordPermissionsShare   = "Allow sharing records."
	DescRecordPermissionsEdit    = "Allow editing records."
	DescRecords                  = "Per-record permissions. Map key is record UID or name; value is an object with can_share and can_edit."
	DescRecordsMD                = "Per-record permissions. Map key is record UID or name; value is an object with `can_share` and `can_edit`."
	DescRecordShare              = "Allow sharing this record."
	DescRecordEdit               = "Allow editing this record."
	DescUsers                    = "Per-user permissions. Map key is user email or UID; value is an object with manage_users, manage_records, and expiration (\"never\" or yyyy-MM-ddTHH:mm:ss)."
	DescUsersMD                  = "Per-user permissions. Map key is user email or UID; value is an object with `manage_users`, `manage_records`, and `expiration` (`\"never\"` or `yyyy-MM-ddTHH:mm:ss`)."
	DescUserManageUsers          = "Allow this user to manage users. Defaults to `false` if not set."
	DescUserManageRecords        = "Allow this user to manage records. Defaults to `false` if not set."
	DescExpiration               = "Access expiration: \"never\" or absolute datetime as yyyy-MM-ddTHH:mm:ss (e.g. 2026-04-02T11:11:00). Defaults to `never` if not set."
	ExpirationDoc                = "Must be \"never\" or yyyy-MM-ddTHH:mm:ss (e.g. 2026-04-02T11:11:00). Defaults to `never` if not set."
)

// Expiration validator error messages.
const (
	ErrMsgInvalidExpiration = "Invalid expiration"
	ErrMsgExpirationEmpty   = "expiration cannot be empty. Use \"never\" for no expiration."
)

// User map entry: manage_users vs expiration.
const (
	ErrMsgManageUsersWithTimeLimitedExpiration = "The ability to manage users is restricted for users with time-limited access. manage_users cannot be true when expiration is a datetime; use \"never\" for users who manage other users."
	ErrMsgInvalidUserPermissionsCombo          = "Invalid user permissions"
)
