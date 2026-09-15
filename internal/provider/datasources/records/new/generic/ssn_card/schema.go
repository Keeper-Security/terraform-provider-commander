// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordssncard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssn_card"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *SsnCardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"ssn_card": dschema.StringAttribute{
					Required:            true,
					Description:         "SSN Card record title or UID to look up.",
					MarkdownDescription: "SSN Card record **title** or **UID** to look up.",
				},
			},
			commonrecordssncard.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
