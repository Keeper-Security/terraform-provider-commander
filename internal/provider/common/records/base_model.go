// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

import "github.com/hashicorp/terraform-plugin-framework/types"

// BaseVaultRecordModel holds attributes shared by all standard vault record resources.
type BaseVaultRecordModel struct {
	Id     types.String       `tfsdk:"id"`
	Title  types.String       `tfsdk:"title"`
	Notes  types.String       `tfsdk:"notes"`
	Folder types.String       `tfsdk:"folder"`
	Custom []CustomFieldModel `tfsdk:"custom"`
	Share  types.Map          `tfsdk:"share"`
}

// ShareRecordPermissionsModel is the model for the record permissions.
type ShareRecordPermissionsModel struct {
	CanShare types.Bool `tfsdk:"can_share"`
	CanEdit  types.Bool `tfsdk:"can_edit"`
}
