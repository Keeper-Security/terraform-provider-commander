// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package classic_share provides a reusable Terraform attribute, model,
// schema, and helpers for the `share` block used by classic vault record
// resources (e.g. classic PAM records, classic generic records). It wires
// the `share-record` Commander CLI command and per-user can_share/can_edit
// flags. Folder sharing is handled by the classic_folders/shared_folder
// package; NSF (Keeper Drive) sharing is handled by new_share.
package classic_share

// Terraform attribute name for the share block and nested attribute keys.
const (
	AttrShare    = "share"
	AttrCanShare = "can_share"
	AttrCanEdit  = "can_edit"
)

// AttrShareValidatorLabel is the human-readable display name passed to
// MapKeysEmailValidator.
const (
	AttrShareValidatorLabel = "Share User Email"
)

// Commander CLI command for sharing a classic vault record.
const (
	CmdShareRecord = "share-record"
)

// Command flags accepted by share-record.
const (
	FlagEmail  = "--email"
	FlagShare  = "--share"
	FlagWrite  = "--write"
	FlagAction = "--action"
)

// share-record --action values. Grant is the default when --action is omitted.
const (
	// ActionRevoke either removes the user from the record share entirely
	// (when no permission flags are passed) or strips the listed flags
	// (e.g. --action revoke --share --write keeps view-only access).
	ActionRevoke = "revoke"
)

// API response keys for the share fragment of a classic vault record get
// response.
const (
	KeyUserPermissions = "user_permissions"
	KeyUsername        = "username"
	KeyShareable       = "shareable"
	KeyEditable        = "editable"
)

// Error operation messages (second argument to ExecuteCommand; short
// description for logs).
const (
	ErrOpShareGrant  = "Unable to grant record share access"
	ErrOpShareRevoke = "Unable to revoke record share access"
)

// Schema descriptions.
const (
	DescShare    = "Mapping of share permissions for this record. Each map key is a user email; each value is an object with can_share and can_edit booleans."
	DescShareMD  = "Mapping of share permissions for this record. Each map **key** is a **user email**; each **value** is an object with `can_share` and `can_edit` booleans."
	DescCanShare = "Allow the user to re-share this record with other users. " +
		"Defaults to `false`."
	DescCanEdit = "Allow the user to edit this record. Defaults to `false`."
)
