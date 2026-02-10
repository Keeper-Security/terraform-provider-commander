// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &ManageCompanyDataSource{}
var _ datasource.DataSourceWithConfigure = &ManageCompanyDataSource{}

type ManageCompanyDataSource struct {
	utils.BaseDataSource
}

func (d *ManageCompanyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manage_company"
}

func (d *ManageCompanyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.BaseDataSource.ConfigureDataSource(ctx, req, resp)
}

func NewManageCompanyDataSource() datasource.DataSource {
	return &ManageCompanyDataSource{}
}
