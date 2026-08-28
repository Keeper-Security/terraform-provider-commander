// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordmembership "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/membership"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *MembershipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"membership": dschema.StringAttribute{
					Required:            true,
					Description:         "New (NSF) Membership record UID to look up.",
					MarkdownDescription: "New (NSF) Membership record **UID** to look up.",
				},
			},
			commonrecordmembership.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
