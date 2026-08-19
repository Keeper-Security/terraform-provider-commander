// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &HealthInsuranceDataSource{}
var _ datasource.DataSourceWithConfigure = &HealthInsuranceDataSource{}

type HealthInsuranceDataSource struct {
	utils.BaseDataSource
}

func (d *HealthInsuranceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_health_insurance"
}

func (d *HealthInsuranceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewHealthInsuranceDataSource() datasource.DataSource {
	return &HealthInsuranceDataSource{}
}
