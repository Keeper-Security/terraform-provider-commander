// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *NonSharedFolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         DescResource,
		MarkdownDescription: DescResource,
		Attributes: folderutils.MergeResourceAttributes(
			folderutils.ResourceCommonFolderAttributes(),
			map[string]schema.Attribute{
				"folder_location": schema.StringAttribute{
					Optional:            true,
					Description:         DescFolderLocation,
					MarkdownDescription: DescFolderLocation,
				},
				"records": schema.SetAttribute{
					Optional:            true,
					ElementType:         types.StringType,
					Description:         DescRecords,
					MarkdownDescription: DescRecords,
				},
			},
		),
	}
}
