// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &SsnCardDataSource{}
var _ datasource.DataSourceWithConfigure = &SsnCardDataSource{}

type SsnCardDataSource struct {
	utils.BaseDataSource
}

func (d *SsnCardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_ssn_card"
}

func (d *SsnCardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewSsnCardDataSource() datasource.DataSource {
	return &SsnCardDataSource{}
}
