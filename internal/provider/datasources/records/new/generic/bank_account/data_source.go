// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &BankAccountDataSource{}
var _ datasource.DataSourceWithConfigure = &BankAccountDataSource{}

type BankAccountDataSource struct {
	utils.BaseDataSource
}

func (d *BankAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_bank_account"
}

func (d *BankAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewBankAccountDataSource() datasource.DataSource {
	return &BankAccountDataSource{}
}
