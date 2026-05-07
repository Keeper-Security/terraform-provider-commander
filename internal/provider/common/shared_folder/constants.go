// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package shared_folder holds types, API read/mapping helpers, and attribute keys shared by
// the commander_shared_folder resource and data source.
package shared_folder

// Object attribute keys (match resource/data source nested schemas and API maps).
const (
	AttrCanShare      = "can_share"
	AttrCanEdit       = "can_edit"
	AttrManageUsers   = "manage_users"
	AttrManageRecords = "manage_records"
	AttrExpiration    = "expiration"
)

const (
	cmdGet     = "get"
	flagFormat = "--format"
	formatJSON = "json"
	errOpGet   = "Unable to get shared folder"
)
