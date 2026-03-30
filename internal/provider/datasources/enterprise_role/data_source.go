// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseRoleDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseRoleDataSource{}

type EnterpriseRoleDataSource struct {
	utils.BaseDataSource
}

func (d *EnterpriseRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_role"
}

func (d *EnterpriseRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewEnterpriseRoleDataSource() datasource.DataSource {
	return &EnterpriseRoleDataSource{}
}
