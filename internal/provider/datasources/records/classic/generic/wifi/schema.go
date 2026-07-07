// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordwifi "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/wifi"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *WifiDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"wifi": dschema.StringAttribute{
					Required:            true,
					Description:         "Vault record UID or title of the WIFI Login record to read.",
					MarkdownDescription: "Vault **record UID** or title of the **WIFI Login** record to read.",
				},
			},
			commonrecordwifi.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
