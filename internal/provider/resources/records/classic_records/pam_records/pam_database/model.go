// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_database"
)

// PamDatabaseResourceModel is the classic PAM Database resource state model:
// the shared PAM Database fields plus the `share` attribute reconciled via
// the classic_share package and the `share-record` Commander CLI.
type PamDatabaseResourceModel struct {
	commonpamdatabase.PamDatabaseResourceModel
	classic_share.ShareModel
}
