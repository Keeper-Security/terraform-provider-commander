// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PamConfigurationDataSource{}
var _ datasource.DataSourceWithConfigure = &PamConfigurationDataSource{}

type PamConfigurationDataSource struct {
	utils.BaseDataSource
}

func (d *PamConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_configuration"
}

func (d *PamConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPamConfigurationDataSource() datasource.DataSource {
	return &PamConfigurationDataSource{}
}
