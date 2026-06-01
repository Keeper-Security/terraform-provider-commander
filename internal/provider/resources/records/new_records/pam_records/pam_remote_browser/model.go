// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamremotebrowser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
)

// PamRemoteBrowserResourceModel is the new (nested-shared) PAM remote browser
// resource state model: the shared PAM remote browser fields plus the `share`
// attribute.
type PamRemoteBrowserResourceModel struct {
	commonpamremotebrowser.PamRemoteBrowserResourceModel
	new_share.ShareModel
}
