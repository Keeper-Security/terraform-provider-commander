// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &AddressDataSource{}
var _ datasource.DataSourceWithConfigure = &AddressDataSource{}

type AddressDataSource struct {
	utils.BaseDataSource
}

func (d *AddressDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_address"
}

func (d *AddressDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewAddressDataSource() datasource.DataSource {
	return &AddressDataSource{}
}
