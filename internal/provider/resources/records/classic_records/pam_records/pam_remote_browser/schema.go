// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"

	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamRemoteBrowserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages PAM remote browser record in your Keeper vault.\n\n" +
			"A PAM Remote Browser is a type of KeeperPAM resource that represents a remote browser isolation target, such as a protected internal application or cloud-based web app.\n\n" +
			"For more information, see https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-remote-browser.",
		MarkdownDescription: "Manages **PAM remote browser** record in your Keeper vault.\n\n" +
			"A PAM Remote Browser is a type of KeeperPAM resource that represents a remote browser isolation target, such as a protected internal application or cloud-based web app.\n\n" +
			"For more information, see [Keeper PAM Remote Browser documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-remote-browser).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         commonpamremotebrowser.IDDescription,
				MarkdownDescription: commonpamremotebrowser.IDMarkdownDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         commonpamremotebrowser.TitleDescription,
				MarkdownDescription: commonpamremotebrowser.TitleMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"url": schema.StringAttribute{
				Required:            true,
				Description:         commonpamremotebrowser.URLDescription,
				MarkdownDescription: commonpamremotebrowser.URLMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("URL", 1, false),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamremotebrowser.NotesDescription,
				MarkdownDescription: commonpamremotebrowser.NotesMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Notes", 0, true),
				},
			},
			"folder": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamremotebrowser.FolderDescription,
				MarkdownDescription: commonpamremotebrowser.FolderMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Folder", 1, true),
				},
			},
			"pam_remote_browser_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         commonpamremotebrowser.PamRemoteBrowserSettingsDescription,
				MarkdownDescription: commonpamremotebrowser.PamRemoteBrowserSettingsMarkdownDescription,
				Attributes: map[string]schema.Attribute{
					"configuration": schema.StringAttribute{
						Required:            true,
						Description:         commonpamremotebrowser.SettingsConfigurationDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsConfigurationMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("PAM Configuration UID", 1, false),
						},
					},
					"remote_browser_isolation": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsRemoteBrowserIsolationDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsRemoteBrowserIsolationMarkdownDescription,
					},
					"connections_recording": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsConnectionsRecordingDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsConnectionsRecordingMarkdownDescription,
					},
					"key_events": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsKeyEventsDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsKeyEventsMarkdownDescription,
					},
					"allow_url_navigation": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAllowURLNavigationDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAllowURLNavigationMarkdownDescription,
					},
					"ignore_server_cert": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsIgnoreServerCertDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsIgnoreServerCertMarkdownDescription,
					},
					"allowed_urls": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         commonpamremotebrowser.SettingsAllowedURLsDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAllowedURLsMarkdownDescription,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Allowed URLs"),
							utils.SetNoEmptyStringsValidator("Allowed URLs"),
						},
					},
					"allowed_resource_urls": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         commonpamremotebrowser.SettingsAllowedResourceURLsDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAllowedResourceURLsMarkdownDescription,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Allowed Resource URLs"),
							utils.SetNoEmptyStringsValidator("Allowed Resource URLs"),
						},
					},
					"auto_fill_targets": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         commonpamremotebrowser.SettingsAutoFillTargetsDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAutoFillTargetsMarkdownDescription,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Auto-fill Targets"),
							utils.SetNoEmptyStringsValidator("Auto-fill Targets"),
						},
					},
					"auto_fill_credentials": schema.StringAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAutoFillCredentialsDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAutoFillCredentialsMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Auto-fill Credentials", 1, true),
						},
					},
					"allow_copy": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAllowCopyDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAllowCopyMarkdownDescription,
					},
					"allow_paste": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAllowPasteDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAllowPasteMarkdownDescription,
					},
					"disable_audio": schema.BoolAttribute{
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsDisableAudioDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsDisableAudioMarkdownDescription,
					},
					"audio_channels": schema.Int32Attribute{
						Computed:            true,
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAudioChannelsDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAudioChannelsMarkdownDescription,
						Validators: []validator.Int32{
							AudioChannelsValidator{},
						},
						Default: int32default.StaticInt32(2),
					},
					"audio_bit_depth": schema.Int64Attribute{
						Computed:            true,
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAudioBitDepthDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAudioBitDepthMarkdownDescription,
						Validators: []validator.Int64{
							AudioBitDepthValidator{},
						},
						Default: int64default.StaticInt64(16),
					},
					"audio_sample_rate": schema.Int64Attribute{
						Computed:            true,
						Optional:            true,
						Description:         commonpamremotebrowser.SettingsAudioSampleRateDescription,
						MarkdownDescription: commonpamremotebrowser.SettingsAudioSampleRateMarkdownDescription,
						Default:             int64default.StaticInt64(44100),
					},
				},
			},
		},
	}
}
