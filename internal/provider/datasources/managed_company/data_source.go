// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &ManagedCompanyDataSource{}
var _ datasource.DataSourceWithConfigure = &ManagedCompanyDataSource{}

type ManagedCompanyDataSource struct {
	utils.BaseDataSource
}

func (d *ManagedCompanyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_company"
}

func (d *ManagedCompanyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewManagedCompanyDataSource() datasource.DataSource {
	return &ManagedCompanyDataSource{}
}
