// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	commonnonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/non_shared_folder"
)

// NonSharedFolderResourceModel is the resource state model (same fields as the data source's embedded model).
type NonSharedFolderResourceModel = commonnonsharedfolder.Model
