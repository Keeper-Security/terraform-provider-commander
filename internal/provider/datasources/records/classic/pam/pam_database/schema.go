// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamDatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a classic PAM database record by UID.",
		MarkdownDescription: "Use this data source to look up a **classic PAM database** record by **UID**.",
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"pam_database": dschema.StringAttribute{
					Required:            true,
					Description:         "PAM database record UID to read.",
					MarkdownDescription: "PAM database record **UID** to read.",
				},
			},
			commonpamdatabase.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
