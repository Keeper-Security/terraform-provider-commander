// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *ShareFolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         DescResource,
		MarkdownDescription: DescResource,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         DescId,
				MarkdownDescription: DescId,
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         DescName,
				MarkdownDescription: DescName,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Shared Folder Name", 1, false),
				},
			},
			"folder_location": schema.StringAttribute{
				Optional:            true,
				Description:         DescFolderLocation,
				MarkdownDescription: DescFolderLocation,
			},
			"user_permissions": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					userPermissionsDefaultPlanModifier{},
				},
				Description:         DescUserPermissions,
				MarkdownDescription: DescUserPermissionsMD,
				Attributes: map[string]schema.Attribute{
					"manage_users": schema.BoolAttribute{
						Optional:            true,
						Description:         DescUserPermissionsManage,
						MarkdownDescription: DescUserPermissionsManage,
					},
					"manage_records": schema.BoolAttribute{
						Optional:            true,
						Description:         DescUserPermissionsRecords,
						MarkdownDescription: DescUserPermissionsRecords,
					},
				},
			},
			"record_permissions": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					recordPermissionsDefaultPlanModifier{},
				},
				Description:         DescRecordPermissions,
				MarkdownDescription: DescRecordPermissionsMD,
				Attributes: map[string]schema.Attribute{
					"can_share": schema.BoolAttribute{
						Optional:            true,
						Description:         DescRecordPermissionsShare,
						MarkdownDescription: DescRecordPermissionsShare,
					},
					"can_edit": schema.BoolAttribute{
						Optional:            true,
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
						},
						"can_edit": schema.BoolAttribute{
							Optional:            true,
							Description:         DescRecordEdit,
							MarkdownDescription: DescRecordEdit,
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
						},
						"manage_records": schema.BoolAttribute{
							Optional:            true,
							Description:         DescUserManageRecords,
							MarkdownDescription: DescUserManageRecords,
						},
						"expiration": schema.StringAttribute{
							Optional:            true,
							Description:         DescExpiration,
							MarkdownDescription: DescExpiration,
							Validators: []validator.String{
								ExpirationValidator(),
							},
						},
					},
				},
			},
		},
	}
}
