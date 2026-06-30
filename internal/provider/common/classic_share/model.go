// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share

import "github.com/hashicorp/terraform-plugin-framework/types"

// ShareModel is an embeddable Terraform model fragment that exposes the
// `share` map attribute. Compose it into any classic record resource/data
// source model that needs per-user share permissions.
//
// Example:
//
//	type Model struct {
//	    commonpamdatabase.PamDatabaseResourceModel
//	    classic_share.ShareModel
//	}
type ShareModel struct {
	Share types.Map `tfsdk:"share"`
}

// ShareRecordPermissionsModel is the nested element of the `share` map: the
// pair of booleans Commander accepts via --share / --write on share-record.
type ShareRecordPermissionsModel struct {
	CanShare types.Bool `tfsdk:"can_share"`
	CanEdit  types.Bool `tfsdk:"can_edit"`
}
