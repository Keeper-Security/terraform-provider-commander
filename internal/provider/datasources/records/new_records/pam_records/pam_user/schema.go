// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a new (nested-shared) PAM user record by UID or name and read its per-user share permissions.",
		MarkdownDescription: "Use this data source to look up a **new (nested-shared) PAM user** record by **UID** or **name** and read its **per-user share permissions**.",
		Attributes: utils.MergeDataSourceAttributes(map[string]dschema.Attribute{
			"pam_user": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM user record UID or name to read.",
				MarkdownDescription: "PAM user record **UID** or **name** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.IDDescription,
				MarkdownDescription: commonpamuser.IDMarkdownDescription,
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.TitleDescription,
				MarkdownDescription: commonpamuser.TitleMarkdownDescription,
			},
			"login": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.LoginDescription,
				MarkdownDescription: commonpamuser.LoginMarkdownDescription,
			},
			"password": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         commonpamuser.PasswordDescription,
				MarkdownDescription: commonpamuser.PasswordMarkdownDescription,
			},
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.FolderDescription,
				MarkdownDescription: commonpamuser.FolderMarkdownDescription,
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.NotesDescription,
				MarkdownDescription: commonpamuser.NotesMarkdownDescription,
			},
			"distinguished_name": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.DistinguishedNameDescription,
				MarkdownDescription: commonpamuser.DistinguishedNameMarkdownDescription,
			},
			"private_pem_key": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         commonpamuser.PrivatePEMKeyDescription,
				MarkdownDescription: commonpamuser.PrivatePEMKeyMarkdownDescription,
			},
			"connect_database": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamuser.ConnectDatabaseDescription,
				MarkdownDescription: commonpamuser.ConnectDatabaseMarkdownDescription,
			},
			"managed": dschema.BoolAttribute{
				Computed:            true,
				Description:         commonpamuser.ManagedDescription,
				MarkdownDescription: commonpamuser.ManagedMarkdownDescription,
			},
			"rotation_settings": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         commonpamuser.RotationSettingsDescription,
				MarkdownDescription: commonpamuser.RotationSettingsMarkdownDescription,
				Attributes: map[string]dschema.Attribute{
					"rotation_profile": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotProfileDescription,
						MarkdownDescription: commonpamuser.RotProfileMarkdownDescription,
					},
					"configuration": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotConfigDescription,
						MarkdownDescription: commonpamuser.RotConfigMarkdownDescription,
					},
					"iam_aad_config": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotIamAadConfigDescription,
						MarkdownDescription: commonpamuser.RotIamAadConfigMarkdownDescription,
					},
					"resource": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotResourceDescription,
						MarkdownDescription: commonpamuser.RotResourceMarkdownDescription,
					},
					"admin_user": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotAdminUserDescription,
						MarkdownDescription: commonpamuser.RotAdminUserMarkdownDescription,
					},
					"enabled": dschema.BoolAttribute{
						Computed:            true,
						Description:         commonpamuser.RotEnabledDescription,
						MarkdownDescription: commonpamuser.RotEnabledMarkdownDescription,
					},
					"schedule_cron": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotScheduleCronDescription,
						MarkdownDescription: commonpamuser.RotScheduleCronMarkdownDescription,
					},
					"schedule_json": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotScheduleJSONDescription,
						MarkdownDescription: commonpamuser.RotScheduleJSONMarkdownDescription,
					},
					"on_demand": dschema.BoolAttribute{
						Computed:            true,
						Description:         commonpamuser.RotOnDemandDescription,
						MarkdownDescription: commonpamuser.RotOnDemandMarkdownDescription,
					},
					"schedule_config": dschema.BoolAttribute{
						Computed:            true,
						Description:         commonpamuser.RotScheduleConfigDescription,
						MarkdownDescription: commonpamuser.RotScheduleConfigMarkdownDescription,
					},
					"complexity": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamuser.RotComplexityDescription,
						MarkdownDescription: commonpamuser.RotComplexityMarkdownDescription,
					},
				},
			},
		}, new_share.DataSourceShareAttribute()),
	}
}
