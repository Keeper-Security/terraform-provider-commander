// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	commonnonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/non_shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FolderDataSourceModel maps the data source schema attributes.
// Embeds commonnonsharedfolder.Model for id/name/folder_location/records.
type FolderDataSourceModel struct {
	Folder types.String `tfsdk:"folder"`
	commonnonsharedfolder.Model
}
