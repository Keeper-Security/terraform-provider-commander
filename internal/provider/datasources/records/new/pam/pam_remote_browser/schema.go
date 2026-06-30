// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamremotebrowser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamRemoteBrowserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a new (NSF) PAM remote browser record by UID.",
		MarkdownDescription: "Use this data source to look up a **new (NSF) PAM remote browser** record by **UID**.",
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"remote_browser": dschema.StringAttribute{
					Required:            true,
					Description:         "PAM remote browser record UID to read.",
					MarkdownDescription: "PAM remote browser record **UID** to read.",
				},
			},
			commonpamremotebrowser.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
