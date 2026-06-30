// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Defaults when user_permissions / record_permissions are omitted (whole block null).
var (
	userPermissionsSchemaDefault = types.ObjectValueMust(
		map[string]attr.Type{
			AttrManageUsers:   types.BoolType,
			AttrManageRecords: types.BoolType,
		},
		map[string]attr.Value{
			AttrManageUsers:   types.BoolValue(false),
			AttrManageRecords: types.BoolValue(false),
		},
	)
	recordPermissionsSchemaDefault = types.ObjectValueMust(
		map[string]attr.Type{
			AttrCanShare: types.BoolType,
			AttrCanEdit:  types.BoolType,
		},
		map[string]attr.Value{
			AttrCanShare: types.BoolValue(false),
			AttrCanEdit:  types.BoolValue(false),
		},
	)
)

func (r *ClassicSharedFolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         DescResource,
		MarkdownDescription: DescResource,
		Attributes: utils.MergeResourceAttributes(
			folderutils.ResourceCommonFolderAttributes(),
			map[string]schema.Attribute{
				"user_permissions": schema.SingleNestedAttribute{
					Optional:            true,
					Computed:            true,
					Default:             objectdefault.StaticValue(userPermissionsSchemaDefault),
					Description:         DescUserPermissions,
					MarkdownDescription: DescUserPermissionsMD,
					Attributes: map[string]schema.Attribute{
						"manage_users": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							Description:         DescUserPermissionsManage,
							MarkdownDescription: DescUserPermissionsManage,
						},
						"manage_records": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							Description:         DescUserPermissionsRecords,
							MarkdownDescription: DescUserPermissionsRecords,
						},
					},
				},
				"record_permissions": schema.SingleNestedAttribute{
					Optional:            true,
					Computed:            true,
					Default:             objectdefault.StaticValue(recordPermissionsSchemaDefault),
					Description:         DescRecordPermissions,
					MarkdownDescription: DescRecordPermissionsMD,
					Attributes: map[string]schema.Attribute{
						"can_share": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							Description:         DescRecordPermissionsShare,
							MarkdownDescription: DescRecordPermissionsShare,
						},
						"can_edit": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							Description:         DescRecordPermissionsEdit,
							MarkdownDescription: DescRecordPermissionsEdit,
						},
					},
				},
				"records": schema.MapNestedAttribute{
					Optional:            true,
					Description:         DescRecords,
					MarkdownDescription: DescRecordsMD,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"can_share": schema.BoolAttribute{
								Optional:            true,
								Description:         DescRecordShare,
								MarkdownDescription: DescRecordShare,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"can_edit": schema.BoolAttribute{
								Optional:            true,
								Description:         DescRecordEdit,
								MarkdownDescription: DescRecordEdit,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
						},
					},
				},
				"users": schema.MapNestedAttribute{
					Optional:            true,
					Description:         DescUsers,
					MarkdownDescription: DescUsersMD,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"manage_users": schema.BoolAttribute{
								Optional:            true,
								Description:         DescUserManageUsers,
								MarkdownDescription: DescUserManageUsers,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"manage_records": schema.BoolAttribute{
								Optional:            true,
								Description:         DescUserManageRecords,
								MarkdownDescription: DescUserManageRecords,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
						},
					},
				},
			},
		),
	}
}
