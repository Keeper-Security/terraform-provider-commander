// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamRemoteBrowserDataSourceModel struct {
	RemoteBrowser types.String `tfsdk:"remote_browser"`
	commonpamremotebrowser.PamRemoteBrowserResourceModel
	classic_share.ShareModel
}
