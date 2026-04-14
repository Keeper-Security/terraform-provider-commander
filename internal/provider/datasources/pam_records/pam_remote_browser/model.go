// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PamRemoteBrowserDataSourceModel is the Terraform model for commander_pam_remote_browser data source.
type PamRemoteBrowserDataSourceModel struct {
	RemoteBrowser types.String `tfsdk:"remote_browser"`
	commonpamremotebrowser.PamRemoteBrowserResourceModel
}
