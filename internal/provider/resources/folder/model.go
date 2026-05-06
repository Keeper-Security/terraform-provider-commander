// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import "github.com/hashicorp/terraform-plugin-framework/types"

// FolderResourceModel is the Terraform state model for the commander_folder resource.
type FolderResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	FolderLocation types.String `tfsdk:"folder_location"`
	Color          types.String `tfsdk:"color"`
	Records        types.Set    `tfsdk:"records"`
}
