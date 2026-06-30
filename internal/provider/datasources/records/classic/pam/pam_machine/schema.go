// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamMachineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a classic PAM machine record by UID.",
		MarkdownDescription: "Use this data source to look up a **classic PAM machine** record by **UID**.",
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"pam_machine": dschema.StringAttribute{
					Required:            true,
					Description:         "PAM machine record UID to read.",
					MarkdownDescription: "PAM machine record **UID** to read.",
				},
			},
			commonpammachine.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
