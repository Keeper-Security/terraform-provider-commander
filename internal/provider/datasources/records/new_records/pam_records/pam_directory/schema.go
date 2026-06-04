// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamDirectoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a new (nested-shared) PAM directory record by UID or name and read its per-user share permissions.",
		MarkdownDescription: "Use this data source to look up a **new (nested-shared) PAM directory** record by **UID** or **name** and read its **per-user share permissions**.",
		Attributes: folderutils.MergeDataSourceAttributes(map[string]dschema.Attribute{
			"pam_directory": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM directory record UID or name to read.",
				MarkdownDescription: "PAM directory record **UID** or **name** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.IDDescription,
				MarkdownDescription: commonpamdirectory.IDMarkdownDescription,
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.TitleDescription,
				MarkdownDescription: commonpamdirectory.TitleMarkdownDescription,
			},
			"hostname_or_ip": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         commonpamdirectory.HostnameOrIPDescription,
				MarkdownDescription: commonpamdirectory.HostnameOrIPMarkdownDescription,
				Attributes: map[string]dschema.Attribute{
					"hostname": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpamdirectory.HostNameDescription,
						MarkdownDescription: commonpamdirectory.HostNameMarkdownDescription,
					},
					"administrative_port": dschema.Int32Attribute{
						Computed:            true,
						Description:         commonpamdirectory.PortDescription,
						MarkdownDescription: commonpamdirectory.PortMarkdownDescription,
					},
				},
			},
			"use_ssl": dschema.BoolAttribute{
				Computed:            true,
				Description:         commonpamdirectory.UseSSLDescription,
				MarkdownDescription: commonpamdirectory.UseSSLMarkdownDescription,
			},
			"domain_name": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.DomainNameDescription,
				MarkdownDescription: commonpamdirectory.DomainNameMarkdownDescription,
			},
			"alternative_ips": dschema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         commonpamdirectory.AlternativeIPsDescription,
				MarkdownDescription: commonpamdirectory.AlternativeIPsMarkdownDescription,
			},
			"directory_id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.DirectoryIdDescription,
				MarkdownDescription: commonpamdirectory.DirectoryIdMarkdownDescription,
			},
			"directory_type": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.DirectoryTypeDescription,
				MarkdownDescription: commonpamdirectory.DirectoryTypeMarkdownDescription,
			},
			"user_match": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.UserMatchDescription,
				MarkdownDescription: commonpamdirectory.UserMatchMarkdownDescription,
			},
			"provider_group": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.ProviderGroupDescription,
				MarkdownDescription: commonpamdirectory.ProviderGroupMarkdownDescription,
			},
			"provider_region": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.ProviderRegionDescription,
				MarkdownDescription: commonpamdirectory.ProviderRegionMarkdownDescription,
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.NotesDescription,
				MarkdownDescription: commonpamdirectory.NotesMarkdownDescription,
			},
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.FolderDescription,
				MarkdownDescription: commonpamdirectory.FolderMarkdownDescription,
			},
			"pam_settings": commonpamrecords.CommonPamSettingsDataSourceAttribute(),
		}, new_share.DataSourceShareAttribute()),
	}
}
