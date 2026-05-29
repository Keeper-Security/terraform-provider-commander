// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *NonSharedFolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         DescDataSource,
		MarkdownDescription: DescDataSourceMD,
		Attributes: folderutils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"folder": dschema.StringAttribute{
					Required:            true,
					Description:         DescDataSourceFolder,
					MarkdownDescription: DescDataSourceFolderMD,
				},
			},
			folderutils.DataSourceCommonFolderAttributes(),
			map[string]dschema.Attribute{
				"records": dschema.SetAttribute{
					Computed:            true,
					ElementType:         types.StringType,
					Description:         DescRecords,
					MarkdownDescription: DescRecords,
				},
			},
		),
	}
}
