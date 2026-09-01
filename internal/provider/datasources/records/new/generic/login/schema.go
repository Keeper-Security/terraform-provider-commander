// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordlogin "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/login"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *LoginDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		// Required lookup must be merged last so it is not overwritten by the
		// computed `login` username field from SharedDataSourceAttributes.
		Attributes: utils.MergeDataSourceAttributes(
			commonrecordlogin.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
			map[string]dschema.Attribute{
				"login_record": dschema.StringAttribute{
					Required:            true,
					Description:         "New (NSF) Login record title or UID to look up.",
					MarkdownDescription: "New (NSF) Login record **title** or **UID** to look up.",
				},
			},
		),
	}
}
