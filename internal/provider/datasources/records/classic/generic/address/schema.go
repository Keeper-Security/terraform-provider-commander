// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordaddress "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/address"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *AddressDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"location": dschema.StringAttribute{
					Required:            true,
					Description:         "Address record title or UID to look up.",
					MarkdownDescription: "Address record **title** or **UID** to look up.",
				},
			},
			commonrecordaddress.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
