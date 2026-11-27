// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamcloud

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PamCloudDataSource{}
var _ datasource.DataSourceWithConfigure = &PamCloudDataSource{}

type PamCloudDataSource struct {
	utils.BaseDataSource
}

func (d *PamCloudDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_cloud"
}

func (d *PamCloudDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPamCloudDataSource() datasource.DataSource {
	return &PamCloudDataSource{}
}
