// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *PamDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages PAM database record with pam settings in your Keeper vault.\n\n" +
			"A PAM Database record is a type of KeeperPAM resource that represents a database server, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database.",
		MarkdownDescription: "Creates and manages **PAM database record with pam settings** in your Keeper vault.\n\n" +
			"A PAM Database record is a type of KeeperPAM resource that represents a database server, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the [PAM Database documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.IDDescription,
				MarkdownDescription: commonpamdatabase.IDMarkdownDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         commonpamdatabase.TitleDescription,
				MarkdownDescription: commonpamdatabase.TitleMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"hostname_or_ip": schema.SingleNestedAttribute{
				Required:            true,
				Description:         commonpamdatabase.HostnameOrIPDescription,
				MarkdownDescription: commonpamdatabase.HostnameOrIPMarkdownDescription,
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						Required:            true,
						Description:         commonpamdatabase.HostNameDescription,
						MarkdownDescription: commonpamdatabase.HostNameMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Host Name", 1, false),
						},
					},
					"administrative_port": schema.Int32Attribute{
						Optional:            true,
						Description:         commonpamdatabase.PortDescription,
						MarkdownDescription: commonpamdatabase.PortMarkdownDescription,
						Validators: []validator.Int32{
							utils.Int32NonNegativeValidator("Administrative Port", true),
						},
					},
				},
			},
			"use_ssl": schema.BoolAttribute{
				Optional:            true,
				Description:         commonpamdatabase.UseSSLDescription,
				MarkdownDescription: commonpamdatabase.UseSSLMarkdownDescription,
			},
			"database_id": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdatabase.DatabaseIdDescription,
				MarkdownDescription: commonpamdatabase.DatabaseIdMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Database ID", 1, true),
				},
			},
			"database_type": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdatabase.DatabaseTypeDescription,
				MarkdownDescription: commonpamdatabase.DatabaseTypeMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Database Type", []string{
						"postgresql", "postgresql-flexible",
						"mysql", "mysql-flexible",
						"mariadb", "mariadb-flexible",
						"mssql", "oracle", "mongodb",
					}, true),
				},
			},
			"provider_group": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdatabase.ProviderGroupDescription,
				MarkdownDescription: commonpamdatabase.ProviderGroupMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Provider Group", 1, true),
				},
			},
			"provider_region": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdatabase.ProviderRegionDescription,
				MarkdownDescription: commonpamdatabase.ProviderRegionMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Provider Region", 1, true),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdatabase.NotesDescription,
				MarkdownDescription: commonpamdatabase.NotesMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Notes", 0, true),
				},
			},
			"folder": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdatabase.FolderDescription,
				MarkdownDescription: commonpamdatabase.FolderMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Folder", 1, true),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(),
		},
	}
}
