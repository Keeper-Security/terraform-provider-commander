// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &ClassicSharedFolderDataSource{}
var _ datasource.DataSourceWithConfigure = &ClassicSharedFolderDataSource{}

type ClassicSharedFolderDataSource struct {
	utils.BaseDataSource
}

func (d *ClassicSharedFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_folder"
}

func (d *ClassicSharedFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewClassicSharedFolderDataSource() datasource.DataSource {
	return &ClassicSharedFolderDataSource{}
}
