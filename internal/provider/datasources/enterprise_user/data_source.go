// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseUserDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseUserDataSource{}

type EnterpriseUserDataSource struct {
	utils.BaseDataSource
}

func (d *EnterpriseUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_user"
}

func (d *EnterpriseUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewEnterpriseUserDataSource() datasource.DataSource {
	return &EnterpriseUserDataSource{}
}
