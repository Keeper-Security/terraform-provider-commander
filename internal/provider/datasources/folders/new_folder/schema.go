// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *NewFolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         DescDataSource,
		MarkdownDescription: DescDataSourceMD,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				AttrNewFolder: dschema.StringAttribute{
					Required:            true,
					Description:         DescDataSourceNewFolder,
					MarkdownDescription: DescDataSourceNewFolderMD,
				},
			},
			folderutils.DataSourceCommonFolderAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
