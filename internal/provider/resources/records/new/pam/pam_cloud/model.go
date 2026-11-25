// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamcloud

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_cloud"
)

// PamCloudResourceModel is the new (nested-shared) PAM Cloud resource
// state model: the shared PAM Cloud fields plus the `share` attribute.
type PamCloudResourceModel struct {
	commonpamcloud.PamCloudResourceModel
	new_share.ShareModel
}
