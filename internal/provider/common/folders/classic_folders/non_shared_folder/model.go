// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package non_shared_folder

import (
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	folderutils.CommonFolderModel
	Records types.Set `tfsdk:"records"`
}
