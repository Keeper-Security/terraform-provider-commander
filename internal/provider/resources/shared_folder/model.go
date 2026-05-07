// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/shared_folder"
)

// UserPermissionsModel is re-exported for call sites that reference the resource nested type by name.
type UserPermissionsModel = commonsharedfolder.UserPermissionsModel

// RecordPermissionsModel is re-exported for call sites that reference the resource nested type by name.
type RecordPermissionsModel = commonsharedfolder.RecordPermissionsModel

// SharedFolderResourceModel is the resource state model (same fields as the data source).
type SharedFolderResourceModel = commonsharedfolder.Model
