// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecorddatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"database": dschema.StringAttribute{
					Required:            true,
					Description:         "Database record UID to look up.",
					MarkdownDescription: "Database record **UID** to look up.",
				},
			},
			commonrecorddatabase.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
