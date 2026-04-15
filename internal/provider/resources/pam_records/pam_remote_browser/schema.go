// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"

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
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         IDDescription,
				MarkdownDescription: IDMarkdownDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         TitleDescription,
				MarkdownDescription: TitleMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"url": schema.StringAttribute{
				Required:            true,
				Description:         URLDescription,
				MarkdownDescription: URLMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("URL", 1, false),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Description:         NotesDescription,
				MarkdownDescription: NotesMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Notes", 0, true),
				},
			},
			"folder": schema.StringAttribute{
				Optional:            true,
				Description:         FolderDescription,
				MarkdownDescription: FolderMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Folder", 1, true),
				},
			},
			"pam_remote_browser_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         PamRemoteBrowserSettingsDescription,
				MarkdownDescription: PamRemoteBrowserSettingsMarkdownDescription,
				Attributes: map[string]schema.Attribute{
					"configuration": schema.StringAttribute{
						Required:            true,
						Description:         SettingsConfigurationDescription,
						MarkdownDescription: SettingsConfigurationMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("PAM Configuration UID", 1, false),
						},
					},
					"remote_browser_isolation": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsRemoteBrowserIsolationDescription,
						MarkdownDescription: SettingsRemoteBrowserIsolationMarkdownDescription,
					},
					"connections_recording": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsConnectionsRecordingDescription,
						MarkdownDescription: SettingsConnectionsRecordingMarkdownDescription,
					},
					"key_events": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsKeyEventsDescription,
						MarkdownDescription: SettingsKeyEventsMarkdownDescription,
					},
					"allow_url_navigation": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsAllowURLNavigationDescription,
						MarkdownDescription: SettingsAllowURLNavigationMarkdownDescription,
					},
					"ignore_server_cert": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsIgnoreServerCertDescription,
						MarkdownDescription: SettingsIgnoreServerCertMarkdownDescription,
					},
					"allowed_urls": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         SettingsAllowedURLsDescription,
						MarkdownDescription: SettingsAllowedURLsMarkdownDescription,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Allowed URLs"),
							utils.SetNoEmptyStringsValidator("Allowed URLs"),
						},
					},
					"allowed_resource_urls": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         SettingsAllowedResourceURLsDescription,
						MarkdownDescription: SettingsAllowedResourceURLsMarkdownDescription,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Allowed Resource URLs"),
							utils.SetNoEmptyStringsValidator("Allowed Resource URLs"),
						},
					},
					"auto_fill_targets": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         SettingsAutoFillTargetsDescription,
						MarkdownDescription: SettingsAutoFillTargetsMarkdownDescription,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Auto-fill Targets"),
							utils.SetNoEmptyStringsValidator("Auto-fill Targets"),
						},
					},
					"auto_fill_credentials": schema.StringAttribute{
						Optional:            true,
						Description:         SettingsAutoFillCredentialsDescription,
						MarkdownDescription: SettingsAutoFillCredentialsMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Auto-fill Credentials", 1, true),
						},
					},
					"allow_copy": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsAllowCopyDescription,
						MarkdownDescription: SettingsAllowCopyMarkdownDescription,
					},
					"allow_paste": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsAllowPasteDescription,
						MarkdownDescription: SettingsAllowPasteMarkdownDescription,
					},
					"disable_audio": schema.BoolAttribute{
						Optional:            true,
						Description:         SettingsDisableAudioDescription,
						MarkdownDescription: SettingsDisableAudioMarkdownDescription,
					},
					"audio_channels": schema.Int32Attribute{
						Computed:            true,
						Optional:            true,
						Description:         SettingsAudioChannelsDescription,
						MarkdownDescription: SettingsAudioChannelsMarkdownDescription,
						Validators: []validator.Int32{
							AudioChannelsValidator{},
						},
						Default: int32default.StaticInt32(2),
					},
					"audio_bit_depth": schema.Int64Attribute{
						Computed:            true,
						Optional:            true,
						Description:         SettingsAudioBitDepthDescription,
						MarkdownDescription: SettingsAudioBitDepthMarkdownDescription,
						Validators: []validator.Int64{
							AudioBitDepthValidator{},
						},
						Default: int64default.StaticInt64(16),
					},
					"audio_sample_rate": schema.Int64Attribute{
						Computed:            true,
						Optional:            true,
						Description:         SettingsAudioSampleRateDescription,
						MarkdownDescription: SettingsAudioSampleRateMarkdownDescription,
						Default:             int64default.StaticInt64(44100),
					},
				},
			},
		},
	}
}
