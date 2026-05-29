// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// CommonFolderModel holds id, name and folder_location attributes shared by folder
// resources and data sources. folder_location holds the parent vault path (null
// or empty at vault root).
type CommonFolderModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	FolderLocation types.String `tfsdk:"folder_location"`
}
