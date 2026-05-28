// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// CommonFolderModel holds id and name attributes shared by folder resources and data sources.
type CommonFolderModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
