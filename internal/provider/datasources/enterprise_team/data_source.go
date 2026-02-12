// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseTeamDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseTeamDataSource{}

type EnterpriseTeamDataSource struct {
	utils.BaseDataSource
}

func (d *EnterpriseTeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_team"
}

func (d *EnterpriseTeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewEnterpriseTeamDataSource() datasource.DataSource {
	return &EnterpriseTeamDataSource{}
}
