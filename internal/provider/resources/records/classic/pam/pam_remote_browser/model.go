// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_remote_browser"
)

// PamRemoteBrowserResourceModel is the classic PAM Remote Browser resource
// state model: the shared PAM Remote Browser fields plus the `share`
// attribute reconciled via the classic_share package and the
// `share-record` Commander CLI.
type PamRemoteBrowserResourceModel struct {
	commonpamremotebrowser.PamRemoteBrowserResourceModel
	classic_share.ShareModel
}
