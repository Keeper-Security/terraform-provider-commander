// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdatabase

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_database"
)

// PamDatabaseResourceModel is the new (nested-shared) PAM Database resource
// state model: the shared PAM Database fields plus the `share` attribute.
type PamDatabaseResourceModel struct {
	commonpamdatabase.PamDatabaseResourceModel
	new_share.ShareModel
}
