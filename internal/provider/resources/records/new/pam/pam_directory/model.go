// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_directory"
)

// PamDirectoryResourceModel is the new (nested-shared) PAM Directory resource
// state model: the shared PAM Directory fields plus the `share` attribute.
type PamDirectoryResourceModel struct {
	commonpamdirectory.PamDirectoryResourceModel
	new_share.ShareModel
}
