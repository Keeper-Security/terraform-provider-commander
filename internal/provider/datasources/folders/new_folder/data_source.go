// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &NewFolderDataSource{}
var _ datasource.DataSourceWithConfigure = &NewFolderDataSource{}

type NewFolderDataSource struct {
	utils.BaseDataSource
}

func (d *NewFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_folder"
}

func (d *NewFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

// NewNewFolderDataSource is the factory used by the provider to register the
// data source.
func NewNewFolderDataSource() datasource.DataSource {
	return &NewFolderDataSource{}
}
