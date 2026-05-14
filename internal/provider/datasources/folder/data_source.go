// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &FolderDataSource{}
var _ datasource.DataSourceWithConfigure = &FolderDataSource{}

type FolderDataSource struct {
	utils.BaseDataSource
}

func (d *FolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (d *FolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewFolderDataSource() datasource.DataSource {
	return &FolderDataSource{}
}
