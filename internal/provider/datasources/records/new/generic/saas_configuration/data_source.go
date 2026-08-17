// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &SaasConfigurationDataSource{}
var _ datasource.DataSourceWithConfigure = &SaasConfigurationDataSource{}

type SaasConfigurationDataSource struct {
	utils.BaseDataSource
}

func (d *SaasConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_saas_configuration"
}

func (d *SaasConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewSaasConfigurationDataSource() datasource.DataSource {
	return &SaasConfigurationDataSource{}
}
