// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseScimDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseScimDataSource{}

type EnterpriseScimDataSource struct {
	utils.BaseDataSource
}

func (d *EnterpriseScimDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_scim"
}

func (d *EnterpriseScimDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewEnterpriseScimDataSource() datasource.DataSource {
	return &EnterpriseScimDataSource{}
}
