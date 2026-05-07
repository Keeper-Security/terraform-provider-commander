// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EpmPolicyDataSource{}
var _ datasource.DataSourceWithConfigure = &EpmPolicyDataSource{}

// EpmPolicyDataSource reads an existing EPM policy by ID.
type EpmPolicyDataSource struct {
	utils.BaseDataSource
}

func (d *EpmPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_epm_policy"
}

func (d *EpmPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

// NewEpmPolicyDataSource returns a new EPM policy data source.
func NewEpmPolicyDataSource() datasource.DataSource {
	return &EpmPolicyDataSource{}
}
