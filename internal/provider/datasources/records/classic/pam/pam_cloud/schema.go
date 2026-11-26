// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_cloud"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamCloudDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a classic PAM cloud record by UID.",
		MarkdownDescription: "Use this data source to look up a **classic PAM cloud** record by **UID**.",
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"pam_cloud": dschema.StringAttribute{
					Required:            true,
					Description:         "PAM cloud record UID to read.",
					MarkdownDescription: "PAM cloud record **UID** to read.",
				},
			},
			commonpamcloud.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
