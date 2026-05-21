// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SharedFolderDataSourceModel matches commander_shared_folder data source attributes.
type SharedFolderDataSourceModel struct {
	SharedFolder types.String `tfsdk:"shared_folder"`
	commonsharedfolder.Model
}
