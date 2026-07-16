// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_folder

import (
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model is the terraform-plugin-framework model for new folder resource attributes.
type Model struct {
	folderutils.CommonFolderModel
	new_share.ShareModel
	Records types.Set `tfsdk:"records"`
}
