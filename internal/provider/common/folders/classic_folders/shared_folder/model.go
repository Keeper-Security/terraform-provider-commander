// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package shared_folder

import "github.com/hashicorp/terraform-plugin-framework/types"

// UserPermissionsModel is default/global user permissions for the shared folder.
type UserPermissionsModel struct {
	ManageUsers   types.Bool `tfsdk:"manage_users"`
	ManageRecords types.Bool `tfsdk:"manage_records"`
}

// RecordPermissionsModel is default/global record permissions for the shared folder.
type RecordPermissionsModel struct {
	CanShare types.Bool `tfsdk:"can_share"`
	CanEdit  types.Bool `tfsdk:"can_edit"`
}

// Model is the terraform-plugin-framework model for shared folder attributes (resource state and data source).
type Model struct {
	Id                types.String            `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
	UserPermissions   *UserPermissionsModel   `tfsdk:"user_permissions"`
	RecordPermissions *RecordPermissionsModel `tfsdk:"record_permissions"`
	Records           types.Map               `tfsdk:"records"`
	Users             types.Map               `tfsdk:"users"`
}
