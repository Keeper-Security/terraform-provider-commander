// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PamDirectoryDataSource{}
var _ datasource.DataSourceWithConfigure = &PamDirectoryDataSource{}

type PamDirectoryDataSource struct {
	utils.BaseDataSource
}

func (d *PamDirectoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_directory"
}

func (d *PamDirectoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPamDirectoryDataSource() datasource.DataSource {
	return &PamDirectoryDataSource{}
}
