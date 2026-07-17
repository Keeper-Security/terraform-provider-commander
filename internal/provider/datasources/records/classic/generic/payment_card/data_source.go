// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PaymentCardDataSource{}
var _ datasource.DataSourceWithConfigure = &PaymentCardDataSource{}

type PaymentCardDataSource struct {
	utils.BaseDataSource
}

func (d *PaymentCardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_payment_card"
}

func (d *PaymentCardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPaymentCardDataSource() datasource.DataSource {
	return &PaymentCardDataSource{}
}
