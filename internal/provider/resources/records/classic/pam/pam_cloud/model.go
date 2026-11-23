// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_cloud"
)

// PamCloudResourceModel is the classic PAM Cloud resource state model:
// the shared PAM Cloud fields plus the `share` attribute reconciled via
// the classic_share package and the `share-record` Commander CLI.
type PamCloudResourceModel struct {
	commonpamcloud.PamCloudResourceModel
	classic_share.ShareModel
}
