// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamremotebrowser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PamRemoteBrowserDataSource{}
var _ datasource.DataSourceWithConfigure = &PamRemoteBrowserDataSource{}

type PamRemoteBrowserDataSource struct {
	utils.BaseDataSource
}

func (d *PamRemoteBrowserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_remote_browser"
}

func (d *PamRemoteBrowserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPamRemoteBrowserDataSource() datasource.DataSource {
	return &PamRemoteBrowserDataSource{}
}
