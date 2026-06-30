// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpammachine

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamMachineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a new (NSF) PAM machine record by UID.",
		MarkdownDescription: "Use this data source to look up a **new (NSF) PAM machine** record by **UID**.",
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"pam_machine": dschema.StringAttribute{
					Required:            true,
					Description:         "PAM machine record UID to read.",
					MarkdownDescription: "PAM machine record **UID** to read.",
				},
			},
			commonpammachine.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
