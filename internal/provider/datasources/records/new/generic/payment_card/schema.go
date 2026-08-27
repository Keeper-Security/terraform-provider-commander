// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordpaymentcard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/payment_card"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PaymentCardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"bank_card": dschema.StringAttribute{
					Required:            true,
					Description:         "Payment Card record title or UID to look up.",
					MarkdownDescription: "Payment Card record **title** or **UID** to look up.",
				},
			},
			commonrecordpaymentcard.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
