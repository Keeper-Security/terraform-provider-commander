// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// BaseVaultRecordModel is the canonical record-level base model embedded by
// every vault record resource and data source, regardless of share variant
// (classic / new) or record family (pam / generic).
//
// Share is intentionally NOT included here: it is composed at the resource
// layer via classic_share.ShareModel or new_share.ShareModel so the base stays
// usable from both variants.
//
// folder_location semantics: vault path of the parent folder where the record
// lives (e.g. "Engineering/Platform"). Null or empty means vault root. Never a
// folder UID.
type BaseVaultRecordModel struct {
	Id             types.String `tfsdk:"id"`
	Title          types.String `tfsdk:"title"`
	Notes          types.String `tfsdk:"notes"`
	FolderLocation types.String `tfsdk:"folder_location"`
}

// CustomFieldModel is one user-defined custom field (maps to API `custom` array).
type CustomFieldModel struct {
	Type      types.String `tfsdk:"type"`
	Label     types.String `tfsdk:"label"`
	Value     types.String `tfsdk:"value"`
	Sensitive types.Bool   `tfsdk:"sensitive"`
}
