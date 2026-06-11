// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
)

// PamDirectoryResourceModel is the classic PAM Directory resource state
// model: the shared PAM Directory fields plus the `share` attribute
// reconciled via the classic_share package and the `share-record`
// Commander CLI.
type PamDirectoryResourceModel struct {
	commonpamdirectory.PamDirectoryResourceModel
	classic_share.ShareModel
}
