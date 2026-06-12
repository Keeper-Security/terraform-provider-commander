// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamDatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a classic PAM database record by UID or name and read its per-user share permissions.",
		MarkdownDescription: "Use this data source to look up a **classic PAM database** record by **UID** or **name** and read its **per-user share permissions**.",
		Attributes: utils.MergeDataSourceAttributes(map[string]dschema.Attribute{
			"pam_database": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM database record UID or name to read.",
				MarkdownDescription: "PAM database record **UID** or **name** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.IDDescription,
				MarkdownDescription: commonpamdatabase.IDMarkdownDescription,
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.TitleDescription,
				MarkdownDescription: commonpamdatabase.TitleMarkdownDescription,
			},
			"hostname_or_ip": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         commonpamdatabase.HostnameOrIPDescription,
				MarkdownDescription: commonpamdatabase.HostnameOrIPMarkdownDescription,
				Attributes: map[string]dschema.Attribute{
					"hostname": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamdatabase.HostNameDescription,
						MarkdownDescription: commonpamdatabase.HostNameMarkdownDescription,
					},
					"administrative_port": dschema.Int32Attribute{
						Computed:            true,
						Description:         commonpamdatabase.PortDescription,
						MarkdownDescription: commonpamdatabase.PortMarkdownDescription,
					},
				},
			},
			"use_ssl": dschema.BoolAttribute{
				Computed:            true,
				Description:         commonpamdatabase.UseSSLDescription,
				MarkdownDescription: commonpamdatabase.UseSSLMarkdownDescription,
			},
			"database_id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.DatabaseIdDescription,
				MarkdownDescription: commonpamdatabase.DatabaseIdMarkdownDescription,
			},
			"database_type": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.DatabaseTypeDescription,
				MarkdownDescription: commonpamdatabase.DatabaseTypeMarkdownDescription,
			},
			"provider_group": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.ProviderGroupDescription,
				MarkdownDescription: commonpamdatabase.ProviderGroupMarkdownDescription,
			},
			"provider_region": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.ProviderRegionDescription,
				MarkdownDescription: commonpamdatabase.ProviderRegionMarkdownDescription,
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.NotesDescription,
				MarkdownDescription: commonpamdatabase.NotesMarkdownDescription,
			},
			"folder_location": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdatabase.FolderDescription,
				MarkdownDescription: commonpamdatabase.FolderMarkdownDescription,
			},
			"pam_settings": commonpamrecords.CommonPamSettingsDataSourceAttribute(commonpamrecords.DatabaseProtocols),
		}, classic_share.DataSourceShareAttribute()),
	}
}
