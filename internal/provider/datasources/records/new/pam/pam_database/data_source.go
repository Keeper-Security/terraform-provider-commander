// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdatabase

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PamDatabaseDataSource{}
var _ datasource.DataSourceWithConfigure = &PamDatabaseDataSource{}

type PamDatabaseDataSource struct {
	utils.BaseDataSource
}

func (d *PamDatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_database"
}

func (d *PamDatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPamDatabaseDataSource() datasource.DataSource {
	return &PamDatabaseDataSource{}
}
