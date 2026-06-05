// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamremotebrowser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamRemoteBrowserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a new (nested-shared) PAM remote browser record by UID or name and read its per-user share permissions. You can use this data source to reference a PAM remote browser record from other resources.",
		MarkdownDescription: "Use this data source to look up a **new (nested-shared) PAM remote browser** record by **UID** or **name** and read its **per-user share permissions**. You can use this data source to reference a PAM remote browser record from other resources.",
		Attributes: utils.MergeDataSourceAttributes(map[string]dschema.Attribute{
			"remote_browser": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM remote browser record UID or name to read.",
				MarkdownDescription: "PAM remote browser record **UID** or **name** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamremotebrowser.IDDescription,
				MarkdownDescription: commonpamremotebrowser.IDMarkdownDescription,
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamremotebrowser.TitleDescription,
				MarkdownDescription: commonpamremotebrowser.TitleMarkdownDescription,
			},
			"url": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamremotebrowser.URLDescription,
				MarkdownDescription: commonpamremotebrowser.URLMarkdownDescription,
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamremotebrowser.NotesDescription,
				MarkdownDescription: commonpamremotebrowser.NotesMarkdownDescription,
			},
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpamremotebrowser.FolderDescription,
				MarkdownDescription: commonpamremotebrowser.FolderMarkdownDescription,
			},
			"pam_remote_browser_settings": dschema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				Description:         commonpamremotebrowser.PamRemoteBrowserSettingsDescription,
				MarkdownDescription: commonpamremotebrowser.PamRemoteBrowserSettingsMarkdownDescription,
				Attributes:          pamRemoteBrowserRBISettingsDataSourceAttributes(),
			},
		}, new_share.DataSourceShareAttribute()),
	}
}

func pamRemoteBrowserRBISettingsDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"configuration": dschema.StringAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsConfigurationDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsConfigurationMarkdownDescription,
		},
		"remote_browser_isolation": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsRemoteBrowserIsolationDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsRemoteBrowserIsolationMarkdownDescription,
		},
		"connections_recording": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsConnectionsRecordingDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsConnectionsRecordingMarkdownDescription,
		},
		"key_events": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsKeyEventsDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsKeyEventsMarkdownDescription,
		},
		"allow_url_navigation": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAllowURLNavigationDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAllowURLNavigationMarkdownDescription,
		},
		"ignore_server_cert": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsIgnoreServerCertDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsIgnoreServerCertMarkdownDescription,
		},
		"allowed_urls": dschema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         commonpamremotebrowser.SettingsAllowedURLsDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAllowedURLsMarkdownDescription,
		},
		"allowed_resource_urls": dschema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         commonpamremotebrowser.SettingsAllowedResourceURLsDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAllowedResourceURLsMarkdownDescription,
		},
		"auto_fill_targets": dschema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         commonpamremotebrowser.SettingsAutoFillTargetsDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAutoFillTargetsMarkdownDescription,
		},
		"auto_fill_credentials": dschema.StringAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAutoFillCredentialsDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAutoFillCredentialsMarkdownDescription,
		},
		"allow_copy": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAllowCopyDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAllowCopyMarkdownDescription,
		},
		"allow_paste": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAllowPasteDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAllowPasteMarkdownDescription,
		},
		"disable_audio": dschema.BoolAttribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsDisableAudioDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsDisableAudioMarkdownDescription,
		},
		"audio_channels": dschema.Int32Attribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAudioChannelsDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAudioChannelsMarkdownDescription,
		},
		"audio_bit_depth": dschema.Int64Attribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAudioBitDepthDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAudioBitDepthMarkdownDescription,
		},
		"audio_sample_rate": dschema.Int64Attribute{
			Computed:            true,
			Description:         commonpamremotebrowser.SettingsAudioSampleRateDescription,
			MarkdownDescription: commonpamremotebrowser.SettingsAudioSampleRateMarkdownDescription,
		},
	}
}
