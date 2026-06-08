// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pamremotebrowser holds shared Terraform schema fragments for PAM remote browser resources.
package pamremotebrowser

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PamRemoteBrowserSettingsModel maps PAM remote browser session/settings options.
type PamRemoteBrowserSettingsModel struct {
	Configuration          types.String `tfsdk:"configuration"`
	RemoteBrowserIsolation types.Bool   `tfsdk:"remote_browser_isolation"`
	ConnectionsRecording   types.Bool   `tfsdk:"connections_recording"`
	KeyEvents              types.Bool   `tfsdk:"key_events"`
	AllowUrlNavigation     types.Bool   `tfsdk:"allow_url_navigation"`
	IgnoreServerCert       types.Bool   `tfsdk:"ignore_server_cert"`
	AllowedUrls            types.Set    `tfsdk:"allowed_urls"`
	AllowedResourceUrls    types.Set    `tfsdk:"allowed_resource_urls"`
	AutoFillTargets        types.Set    `tfsdk:"auto_fill_targets"`
	AutoFillCredentials    types.String `tfsdk:"auto_fill_credentials"`
	AllowCopy              types.Bool   `tfsdk:"allow_copy"`
	AllowPaste             types.Bool   `tfsdk:"allow_paste"`
	DisableAudio           types.Bool   `tfsdk:"disable_audio"`
	AudioChannels          types.Int32  `tfsdk:"audio_channels"`
	AudioBitDepth          types.Int64  `tfsdk:"audio_bit_depth"`
	AudioSampleRate        types.Int64  `tfsdk:"audio_sample_rate"`
}

// PamRemoteBrowserResourceModel is the Terraform state shared by
// commander_classic_pam_remote_browser and commander_new_pam_remote_browser.
type PamRemoteBrowserResourceModel struct {
	commonrecordsutils.BaseVaultRecordModel

	Url types.String `tfsdk:"url"`

	PamRemoteBrowserSettings *PamRemoteBrowserSettingsModel `tfsdk:"pam_remote_browser_settings"`
}
