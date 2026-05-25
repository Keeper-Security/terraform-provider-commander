// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_folder

import (
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
)

// Model is the terraform-plugin-framework model for new folder resource attributes.
type Model struct {
	folderutils.IdentityModel
	new_share.ShareModel
}
