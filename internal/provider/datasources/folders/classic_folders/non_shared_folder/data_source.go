// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &NonSharedFolderDataSource{}
var _ datasource.DataSourceWithConfigure = &NonSharedFolderDataSource{}

type NonSharedFolderDataSource struct {
	utils.BaseDataSource
}

func (d *NonSharedFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_non_shared_folder"
}

func (d *NonSharedFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewNonSharedFolderDataSource() datasource.DataSource {
	return &NonSharedFolderDataSource{}
}
