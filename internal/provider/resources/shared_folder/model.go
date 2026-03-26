// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import "github.com/hashicorp/terraform-plugin-framework/types"

// UserPermissionsModel is the default/global user permissions for the shared folder.
// Allowed keys: manage_users, manage_records. Default when omitted: both false.
type UserPermissionsModel struct {
	ManageUsers   types.Bool `tfsdk:"manage_users"`
	ManageRecords types.Bool `tfsdk:"manage_records"`
}

// RecordPermissionsModel is the default/global record permissions for the shared folder.
// Allowed keys: can_share, can_edit. Default when omitted: both false.
type RecordPermissionsModel struct {
	CanShare types.Bool `tfsdk:"can_share"`
	CanEdit  types.Bool `tfsdk:"can_edit"`
}

type SharedFolderResourceModel struct {
	Id                types.String            `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
	UserPermissions   *UserPermissionsModel   `tfsdk:"user_permissions"`
	RecordPermissions *RecordPermissionsModel `tfsdk:"record_permissions"`
	Records           types.Map               `tfsdk:"records"`
	Users             types.Map               `tfsdk:"users"`
}
