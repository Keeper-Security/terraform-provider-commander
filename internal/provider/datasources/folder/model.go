// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import "github.com/hashicorp/terraform-plugin-framework/types"

// FolderDataSourceModel maps the data source schema attributes.
type FolderDataSourceModel struct {
	Folder         types.String `tfsdk:"folder"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	FolderLocation types.String `tfsdk:"folder_location"`
	Records        types.Set    `tfsdk:"records"`
}
