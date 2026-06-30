// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordcontact "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic/generic/contact"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *ContactDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"contact": dschema.StringAttribute{
					Required:            true,
					Description:         "Contact record title or UID to look up.",
					MarkdownDescription: "Contact record **title** or **UID** to look up.",
				},
			},
			commonrecordcontact.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
