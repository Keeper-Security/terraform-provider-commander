// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *FolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Look up an existing vault folder by UID or path.",
		MarkdownDescription: "Look up an existing vault folder by UID or path.",
		Attributes: map[string]dschema.Attribute{
			"folder": dschema.StringAttribute{
				Required:            true,
				Description:         "Folder UID or vault path to look up.",
				MarkdownDescription: "Folder UID or vault path to look up.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         "The UID of the folder.",
				MarkdownDescription: "The UID of the folder.",
			},
			"name": dschema.StringAttribute{
				Computed:            true,
				Description:         "Folder name.",
				MarkdownDescription: "Folder name.",
			},
			"type": dschema.StringAttribute{
				Computed:            true,
				Description:         "Folder type (e.g. user_folder).",
				MarkdownDescription: "Folder type (e.g. `user_folder`).",
			},
			"folder_location": dschema.StringAttribute{
				Computed:            true,
				Description:         "Parent folder path where the folder resides. Empty if at vault root.",
				MarkdownDescription: "Parent folder path where the folder resides. Empty if at vault root.",
			},
			"records": dschema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Set of record UIDs linked to this folder.",
				MarkdownDescription: "Set of record UIDs linked to this folder.",
			},
		},
	}
}
