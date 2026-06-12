// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamMachineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a classic PAM machine record by UID or name and read its per-user share permissions.",
		MarkdownDescription: "Use this data source to look up a **classic PAM machine** record by **UID** or **name** and read its **per-user share permissions**.",
		Attributes: utils.MergeDataSourceAttributes(map[string]dschema.Attribute{
			"pam_machine": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM machine record UID or name to read.",
				MarkdownDescription: "PAM machine record **UID** or **name** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.IDDescription,
				MarkdownDescription: commonpammachine.IDMarkdownDescription,
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.TitleDescription,
				MarkdownDescription: commonpammachine.TitleMarkdownDescription,
			},
			"hostname_or_ip": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         commonpammachine.HostnameOrIPDescription,
				MarkdownDescription: commonpammachine.HostnameOrIPMarkdownDescription,
				Attributes: map[string]dschema.Attribute{
					"hostname": dschema.StringAttribute{
						Computed:            true,
						Description:         commonpammachine.HostNameDescription,
						MarkdownDescription: commonpammachine.HostNameMarkdownDescription,
					},
					"administrative_port": dschema.Int32Attribute{
						Computed:            true,
						Description:         commonpammachine.PortDescription,
						MarkdownDescription: commonpammachine.PortMarkdownDescription,
					},
				},
			},
			"operating_system": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.OperatingSystemDescription,
				MarkdownDescription: commonpammachine.OperatingSystemMarkdownDescription,
			},
			"instance_name": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.InstanceNameDescription,
				MarkdownDescription: commonpammachine.InstanceNameMarkdownDescription,
			},
			"instance_id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.InstanceIdDescription,
				MarkdownDescription: commonpammachine.InstanceIdMarkdownDescription,
			},
			"provider_group": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.ProviderGroupDescription,
				MarkdownDescription: commonpammachine.ProviderGroupMarkdownDescription,
			},
			"provider_region": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.ProviderRegionDescription,
				MarkdownDescription: commonpammachine.ProviderRegionMarkdownDescription,
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.NotesDescription,
				MarkdownDescription: commonpammachine.NotesMarkdownDescription,
			},
			"folder_location": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.FolderDescription,
				MarkdownDescription: commonpammachine.FolderMarkdownDescription,
			},
			"pam_settings": commonpamrecords.CommonPamSettingsDataSourceAttribute(commonpamrecords.MachineDirectoryProtocols),
		}, classic_share.DataSourceShareAttribute()),
	}
}
