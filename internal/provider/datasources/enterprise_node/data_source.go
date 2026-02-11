// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseNodesDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseNodesDataSource{}

type EnterpriseNodesDataSource struct {
	utils.BaseDataSource
}

func (d *EnterpriseNodesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_node"
}

func (d *EnterpriseNodesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewEnterpriseNodesDataSource() datasource.DataSource {
	return &EnterpriseNodesDataSource{}
}
