// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package passport

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordpassport "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/passport"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PassportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"passport": dschema.StringAttribute{
					Required:            true,
					Description:         "New (NSF) Passport record title or UID to look up.",
					MarkdownDescription: "New (NSF) Passport record **title** or **UID** to look up.",
				},
			},
			commonrecordpassport.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
