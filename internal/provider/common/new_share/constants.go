// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package new_share provides a reusable Terraform attribute, model, schema, and
// helpers for the `share` block used by Keeper Drive (nsf-) folder and record
// resources. The same package wires both `nsf-share-folder` (folders) and
// `nsf-share-record` (records) since they share identical flags and semantics.
package new_share

// Terraform attribute name for the share block.
const (
	AttrShare = "share"
)

// AttrShareValidatorLabel is the human-readable display name passed to
// MapKeysMinLengthValidator and MapValuesStringOneOfValidator. nsf share
// commands accept either a Keeper user UID or an email address as the
// principal, so the validator only enforces a non-empty trimmed key; format
// validation is delegated to the Commander CLI.
const (
	AttrShareValidatorLabel = "Share User UID or Email"
	AttrShareValueLabel     = "Share Permission"
)

// Commander CLI commands used by share operations. Both commands accept the
// same flags (--email, --action, --role). Callers pick which one to pass based
// on whether they are sharing a folder or a record.
const (
	CmdNsfShareFolder = "nsf-share-folder"
	CmdNsfShareRecord = "nsf-share-record"
)

// Command flags shared by nsf-share-folder and nsf-share-record.
const (
	FlagEmail  = "--email"
	FlagAction = "--action"
	FlagRole   = "--role"
)

// share --action values.
const (
	ActionGrant  = "grant"
	ActionRevoke = "revoke" // Used for removing a user from a nsf record
	ActionRemove = "remove" // Used for removing a user from a nsf folder
)

// Permission role values accepted by the `share` attribute. RoleOwner is only
// returned by the API (never written by the user) and is filtered out when
// mapping the API response into the Terraform state.
const (
	RoleViewer              = "viewer"
	RoleShareManager        = "share-manager"
	RoleContentManager      = "content-manager"
	RoleContentShareManager = "content-share-manager"
	RoleFullManager         = "full-manager"

	RoleOwner = "owner"
)

// AllowedRoles is the set of role values a user may set on `share`. RoleOwner
// is intentionally excluded; the Keeper Drive API rejects user-supplied owner
// entries and the response filter drops them on read.
var AllowedRoles = []string{
	RoleViewer,
	RoleShareManager,
	RoleContentManager,
	RoleContentShareManager,
	RoleFullManager,
}

// API response keys for the share fragment of a get response.
const (
	KeyUserPermissions = "user_permissions"
	KeyAccessor        = "accessor"
	KeyRole            = "role"
)

// Error operation messages (second argument to ExecuteCommand; short
// description for logs).
const (
	ErrOpShareGrant  = "Unable to grant share access"
	ErrOpShareRevoke = "Unable to revoke share access"
)

// Schema descriptions.
const (
	DescShare   = "Mapping of share permissions for this folder or record. For folders, keys can identify either users (UID or email) or teams (UID or name). For records, keys can identify users only (UID or email). Values specify the permission level: viewer, share-manager, content-manager, content-share-manager, or full-manager. The owner is implicit and never present in this map; the API rejects owner entries."
	DescShareMD = "Mapping of share permissions for this folder or record. For folders, keys can identify either users (**UID** or **email**) or teams (**UID** or **name**). For records, keys can identify users only (**UID** or **email**). Values specify the permission level: `viewer`, `share-manager`, `content-manager`, `content-share-manager`, or `full-manager`. The owner is implicit and never present in this map; the API rejects owner entries."
)
