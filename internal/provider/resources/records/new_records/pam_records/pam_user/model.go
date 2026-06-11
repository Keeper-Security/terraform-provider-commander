// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
)

// PamUserResourceModel is the new (nested-shared) PAM User resource state
// model: the shared PAM User fields plus the `share` attribute.
type PamUserResourceModel struct {
	commonpamuser.PamUserSharedModel
	new_share.ShareModel
}
