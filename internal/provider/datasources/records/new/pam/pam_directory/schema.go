// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamDirectoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a new (NSF) PAM directory record by UID.",
		MarkdownDescription: "Use this data source to look up a **new (NSF) PAM directory** record by **UID**.",
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"pam_directory": dschema.StringAttribute{
					Required:            true,
					Description:         "PAM directory record UID to read.",
					MarkdownDescription: "PAM directory record **UID** to read.",
				},
			},
			commonpamdirectory.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
