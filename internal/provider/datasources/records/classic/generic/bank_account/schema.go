// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordbankaccount "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/bank_account"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *BankAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"account": dschema.StringAttribute{
					Required:            true,
					Description:         "Bank Account record title or UID to look up.",
					MarkdownDescription: "Bank Account record **title** or **UID** to look up.",
				},
			},
			commonrecordbankaccount.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
