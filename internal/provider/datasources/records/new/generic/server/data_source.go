// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &ServerDataSource{}
var _ datasource.DataSourceWithConfigure = &ServerDataSource{}

type ServerDataSource struct {
	utils.BaseDataSource
}

func (d *ServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_server"
}

func (d *ServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewServerDataSource() datasource.DataSource {
	return &ServerDataSource{}
}
