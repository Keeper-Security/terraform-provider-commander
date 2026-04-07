// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &SecretsManagerDataSource{}
var _ datasource.DataSourceWithConfigure = &SecretsManagerDataSource{}

type SecretsManagerDataSource struct {
	utils.BaseDataSource
}

func (d *SecretsManagerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secrets_manager"
}

func (d *SecretsManagerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewSecretsManagerDataSource() datasource.DataSource {
	return &SecretsManagerDataSource{}
}
