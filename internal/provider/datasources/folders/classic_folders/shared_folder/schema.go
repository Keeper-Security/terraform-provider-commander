// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"

	sfres "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *ClassicSharedFolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         sfres.DescDataSource,
		MarkdownDescription: sfres.DescDataSourceMD,
		Attributes: map[string]dschema.Attribute{
			"shared_folder": dschema.StringAttribute{
				Required:            true,
				Description:         sfres.DescDataSourceSharedFolder,
				MarkdownDescription: sfres.DescDataSourceSharedFolderMD,
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         sfres.DescDataSourceId,
				MarkdownDescription: sfres.DescDataSourceIdMD,
			},
			"name": dschema.StringAttribute{
				Computed:            true,
				Description:         sfres.DescName,
				MarkdownDescription: sfres.DescName,
			},
			"user_permissions": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"manage_users": dschema.BoolAttribute{
						Computed:            true,
						Description:         sfres.DescUserPermissionsManage,
						MarkdownDescription: sfres.DescUserPermissionsManage,
					},
					"manage_records": dschema.BoolAttribute{
						Computed:            true,
						Description:         sfres.DescUserPermissionsRecords,
						MarkdownDescription: sfres.DescUserPermissionsRecords,
					},
				},
			},
			"record_permissions": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"can_share": dschema.BoolAttribute{
						Computed:            true,
						Description:         sfres.DescRecordPermissionsShare,
						MarkdownDescription: sfres.DescRecordPermissionsShare,
					},
					"can_edit": dschema.BoolAttribute{
						Computed:            true,
						Description:         sfres.DescRecordPermissionsEdit,
						MarkdownDescription: sfres.DescRecordPermissionsEdit,
					},
				},
			},
			"records": dschema.MapNestedAttribute{
				Computed:            true,
				Description:         sfres.DescRecords,
				MarkdownDescription: sfres.DescRecordsMD,
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"can_share": dschema.BoolAttribute{
							Computed:            true,
							Description:         sfres.DescRecordShare,
							MarkdownDescription: sfres.DescRecordShare,
						},
						"can_edit": dschema.BoolAttribute{
							Computed:            true,
							Description:         sfres.DescRecordEdit,
							MarkdownDescription: sfres.DescRecordEdit,
						},
					},
				},
			},
			"users": dschema.MapNestedAttribute{
				Computed:            true,
				Description:         sfres.DescUsers,
				MarkdownDescription: sfres.DescUsersMD,
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"manage_users": dschema.BoolAttribute{
							Computed:            true,
							Description:         sfres.DescUserManageUsers,
							MarkdownDescription: sfres.DescUserManageUsers,
						},
						"manage_records": dschema.BoolAttribute{
							Computed:            true,
							Description:         sfres.DescUserManageRecords,
							MarkdownDescription: sfres.DescUserManageRecords,
						},
						"expiration": dschema.StringAttribute{
							Computed:            true,
							Description:         sfres.DescExpiration,
							MarkdownDescription: sfres.DescExpiration,
						},
					},
				},
			},
		},
	}
}
