// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(map[string]dschema.Attribute{
			"pam_user": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM user record UID to read.",
				MarkdownDescription: "PAM user record **UID** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         "Same as record_uid from the vault.",
				MarkdownDescription: "Same as **record_uid** from the vault.",
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         "Title of the PAM User record.",
				MarkdownDescription: "**Title** of the PAM User record.",
			},
			"login": dschema.StringAttribute{
				Computed:            true,
				Description:         "Login (username) for the PAM User.",
				MarkdownDescription: "**Login** (username) for the PAM User.",
			},
			"password": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "Password for the PAM User.",
				MarkdownDescription: "**Password** for the PAM User.",
			},
			"folder_location": dschema.StringAttribute{
				Computed:            true,
				Description:         "Folder UID or path for the record.",
				MarkdownDescription: "**Folder** UID or path for the record.",
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         "Notes on the record, if any.",
				MarkdownDescription: "**Notes** on the record, if any.",
			},
			"distinguished_name": dschema.StringAttribute{
				Computed:            true,
				Description:         "LDAP distinguished name of the PAM User.",
				MarkdownDescription: "**LDAP distinguished name** of the PAM User.",
			},
			"private_pem_key": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "Private PEM key associated with the PAM User.",
				MarkdownDescription: "**Private PEM key** associated with the PAM User.",
			},
			"public_key": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "Public key associated with the PAM User.",
				MarkdownDescription: "**Public key** associated with the PAM User.",
			},
			"private_key_passphrase": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "Passphrase for the private key associated with the PAM User.",
				MarkdownDescription: "Passphrase for the private key associated with the PAM User.",
			},
			"connect_database": dschema.StringAttribute{
				Computed:            true,
				Description:         "Database name the PAM User connects to.",
				MarkdownDescription: "**Database name** the PAM User connects to.",
			},
			"managed": dschema.BoolAttribute{
				Computed:            true,
				Description:         "Whether this PAM User account is managed by Keeper.",
				MarkdownDescription: "Whether this PAM User account is **managed** by Keeper.",
			},
			"rotation_settings": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Rotation settings for the PAM User record, if configured.",
				MarkdownDescription: "**Rotation settings** for the PAM User record, if configured.",
				Attributes: map[string]dschema.Attribute{
					"rotation_profile": dschema.StringAttribute{
						Computed:            true,
						Description:         "Rotation profile type: general, iam_user, or scripts_only.",
						MarkdownDescription: "Rotation profile type: `general`, `iam_user`, or `scripts_only`.",
					},
					"configuration": dschema.StringAttribute{
						Computed:            true,
						Description:         "PAM Configuration UID.",
						MarkdownDescription: "**PAM Configuration UID**.",
					},
					"iam_aad_config": dschema.StringAttribute{
						Computed:            true,
						Description:         "IAM/Azure AD PAM Configuration UID.",
						MarkdownDescription: "**IAM/Azure AD** PAM Configuration UID.",
					},
					"resource": dschema.StringAttribute{
						Computed:            true,
						Description:         "PAM resource record UID (machine or database).",
						MarkdownDescription: "**PAM resource** record UID (machine or database).",
					},
					"admin_user": dschema.StringAttribute{
						Computed:            true,
						Description:         "Admin PAM User UID used for rotation.",
						MarkdownDescription: "**Admin PAM User** UID used for rotation.",
					},
					"enabled": dschema.BoolAttribute{
						Computed:            true,
						Description:         "Whether rotation is enabled.",
						MarkdownDescription: "Whether rotation is **enabled**.",
					},
					"schedule_cron": dschema.StringAttribute{
						Computed:            true,
						Description:         "Cron schedule for rotation.",
						MarkdownDescription: "**Cron schedule** for rotation.",
					},
					"schedule_json": dschema.StringAttribute{
						Computed:            true,
						Description:         "JSON schedule for rotation.",
						MarkdownDescription: "**JSON schedule** for rotation.",
					},
					"on_demand": dschema.BoolAttribute{
						Computed:            true,
						Description:         "Whether rotation is on-demand (manual).",
						MarkdownDescription: "Whether rotation is **on-demand** (manual).",
					},
					"schedule_config": dschema.BoolAttribute{
						Computed:            true,
						Description:         "Whether schedule is inherited from PAM Configuration.",
						MarkdownDescription: "Whether schedule is inherited from **PAM Configuration**.",
					},
					"complexity": dschema.StringAttribute{
						Computed:            true,
						Description:         "Password complexity for rotation.",
						MarkdownDescription: "**Password complexity** for rotation.",
					},
				},
			},
		}, classic_share.DataSourceShareAttribute()),
	}
}
