// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	commonnewfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/new_folder"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NewFolderDataSourceModel struct {
	NewFolder types.String `tfsdk:"new_folder"`
	commonnewfolder.Model
}
