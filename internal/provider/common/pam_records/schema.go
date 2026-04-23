// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// CommonPamSettingsTunnelSchema returns the reusable schema attributes for
// the tunnel (portForward) block inside pam_settings.
// Enabled is required; the remaining fields are optional and only meaningful
// when enabled is true – enforced by TunnelFieldsRequireEnabledValidator.
func CommonPamSettingsTunnelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable": schema.BoolAttribute{
			Required:            true,
			Description:         PamSettingsTunnelEnabledDescription,
			MarkdownDescription: PamSettingsTunnelEnabledMarkdownDescription,
		},
		"remote_target_port": schema.Int32Attribute{
			Optional:            true,
			Description:         PamSettingsTunnelRemoteTargetPortDescription,
			MarkdownDescription: PamSettingsTunnelRemoteTargetPortMarkdownDescription,
			Validators: []validator.Int32{
				utils.Int32NonNegativeValidator("Remote Target Port", true),
			},
		},
		"re_use_port": schema.BoolAttribute{
			Optional:            true,
			Description:         PamSettingsTunnelReUsePortDescription,
			MarkdownDescription: PamSettingsTunnelReUsePortMarkdownDescription,
		},
		"use_specified_local_port": schema.BoolAttribute{
			Optional:            true,
			Description:         PamSettingsTunnelUseSpecifiedLocalPortDescription,
			MarkdownDescription: PamSettingsTunnelUseSpecifiedLocalPortMarkdownDescription,
		},
		"local_port": schema.Int32Attribute{
			Optional:            true,
			Description:         PamSettingsTunnelLocalPortDescription,
			MarkdownDescription: PamSettingsTunnelLocalPortMarkdownDescription,
			Validators: []validator.Int32{
				utils.Int32NonNegativeValidator("Local Port", true),
			},
		},
	}
}

// CommonPamSettingsConnectionSchema returns the reusable schema attributes for
// the connection block inside pam_settings.
// Protocol is required; each protocol-specific block is optional and only the
// one matching the selected protocol may be set – enforced by
// ConnectionProtocolBlockValidator.
func CommonPamSettingsConnectionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable": schema.BoolAttribute{
			Required:            true,
			Description:         PamSettingsConnectionEnableDescription,
			MarkdownDescription: PamSettingsConnectionEnableMarkdownDescription,
		},
		"protocol": schema.StringAttribute{
			Required:            true,
			Description:         PamSettingsConnectionProtocolDescription,
			MarkdownDescription: PamSettingsConnectionProtocolMarkdownDescription,
			Validators: []validator.String{
				ConnectionProtocolValidator(),
			},
		},
		"connection_port": schema.Int32Attribute{
			Optional:            true,
			Description:         PamSettingsConnectionConnectionPortDescription,
			MarkdownDescription: PamSettingsConnectionConnectionPortMarkdownDescription,
			Validators: []validator.Int32{
				utils.Int32NonNegativeValidator("Connection Port", true),
			},
		},
		"launch_credential": schema.StringAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionLaunchCredentialDescription,
			MarkdownDescription: PamSettingsConnectionLaunchCredentialMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Launch Credential", 1, true),
			},
		},
		"kubernetes": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionKubernetesDescription,
			MarkdownDescription: PamSettingsConnectionKubernetesMarkdownDescription,
			Attributes:          ConnectionKubernetesSchema(),
		},
		"mysql": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionMysqlDescription,
			MarkdownDescription: PamSettingsConnectionMysqlMarkdownDescription,
			Attributes:          ConnectionDatabaseSchema(),
		},
		"postgresql": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionPostgreSqlDescription,
			MarkdownDescription: PamSettingsConnectionPostgreSqlMarkdownDescription,
			Attributes:          ConnectionDatabaseSchema(),
		},
		"rdp": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionRdpDescription,
			MarkdownDescription: PamSettingsConnectionRdpMarkdownDescription,
			Attributes:          ConnectionRdpSchema(),
		},
		"sql_server": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionSqlServerDescription,
			MarkdownDescription: PamSettingsConnectionSqlServerMarkdownDescription,
			Attributes:          ConnectionDatabaseSchema(),
		},
		"ssh": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionSshDescription,
			MarkdownDescription: PamSettingsConnectionSshMarkdownDescription,
			Attributes:          ConnectionSshSchema(),
		},
		"telnet": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionTelnetDescription,
			MarkdownDescription: PamSettingsConnectionTelnetMarkdownDescription,
			Attributes:          map[string]schema.Attribute{},
		},
		"vnc": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionVncDescription,
			MarkdownDescription: PamSettingsConnectionVncMarkdownDescription,
			Attributes:          map[string]schema.Attribute{},
		},
	}
}

// mergeSchemaAttributes combines multiple attribute maps into one.
// Later maps take precedence if keys overlap.
func mergeSchemaAttributes(maps ...map[string]schema.Attribute) map[string]schema.Attribute {
	result := map[string]schema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// ConnectionCommonSchema returns the 4 shared attributes used by all
// protocol models (Kubernetes, Database, RDP).
func ConnectionCommonSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"session_recording": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionSessionRecordingDescription,
			MarkdownDescription: ConnectionSessionRecordingMarkdownDescription,
		},
		"recording_include_keys": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionRecordingIncludeKeysDescription,
			MarkdownDescription: ConnectionRecordingIncludeKeysMarkdownDescription,
		},
		"allow_supply_user": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionAllowSupplyUserDescription,
			MarkdownDescription: ConnectionAllowSupplyUserMarkdownDescription,
		},
		"read_only": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionReadOnlyDescription,
			MarkdownDescription: ConnectionReadOnlyMarkdownDescription,
		},
	}
}

// ConnectionTerminalSchema returns the 5 terminal-related attributes
// shared by Kubernetes and Database protocols.
func ConnectionTerminalSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"typescript_recording": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionTypescriptRecordingDescription,
			MarkdownDescription: ConnectionTypescriptRecordingMarkdownDescription,
		},
		"color_scheme": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("black-white"),
			Description:         ConnectionColorSchemeDescription,
			MarkdownDescription: ConnectionColorSchemeMarkdownDescription,
			Validators: []validator.String{
				ColorSchemeValidator(),
			},
		},
		"font_name": schema.StringAttribute{
			Optional:            true,
			Description:         ConnectionFontNameDescription,
			MarkdownDescription: ConnectionFontNameMarkdownDescription,
		},
		"font_size": schema.Int32Attribute{
			Optional:            true,
			Description:         ConnectionFontSizeDescription,
			MarkdownDescription: ConnectionFontSizeMarkdownDescription,
			Validators: []validator.Int32{
				utils.Int32OneOfValidator("Font Size", []int32{8, 9, 10, 11, 12, 14, 18, 24, 30, 36, 48, 60, 72, 96}, true),
			},
		},
		"scrollback": schema.Int32Attribute{
			Optional:            true,
			Description:         ConnectionScrollbackDescription,
			MarkdownDescription: ConnectionScrollbackMarkdownDescription,
			Validators: []validator.Int32{
				utils.Int32NonNegativeValidator("Scrollback", true),
			},
		},
	}
}

// ConnectionClipboardSchema returns the 2 clipboard attributes shared by
// Database and RDP protocols (with default false).
func ConnectionClipboardSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"disable_copy": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			Description:         ConnectionDisableCopyDescription,
			MarkdownDescription: ConnectionDisableCopyMarkdownDescription,
		},
		"disable_paste": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			Description:         ConnectionDisablePasteDescription,
			MarkdownDescription: ConnectionDisablePasteMarkdownDescription,
		},
	}
}

// ConnectionKubernetesSchema returns the schema attributes for
// the kubernetes protocol-specific connection block.
func ConnectionKubernetesSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionCommonSchema(),
		ConnectionTerminalSchema(),
		map[string]schema.Attribute{
			"rotate_on_termination": schema.BoolAttribute{
				Optional:            true,
				Description:         KubernetesRotateOnTerminationDescription,
				MarkdownDescription: KubernetesRotateOnTerminationMarkdownDescription,
			},
			"use_ssl": schema.BoolAttribute{
				Optional:            true,
				Description:         KubernetesUseSSLDescription,
				MarkdownDescription: KubernetesUseSSLMarkdownDescription,
			},
			"ignore_cert": schema.BoolAttribute{
				Optional:            true,
				Description:         KubernetesIgnoreCertDescription,
				MarkdownDescription: KubernetesIgnoreCertMarkdownDescription,
			},
			"ca_cert": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesCaCertDescription,
				MarkdownDescription: KubernetesCaCertMarkdownDescription,
			},
			"client_cert": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesClientCertDescription,
				MarkdownDescription: KubernetesClientCertMarkdownDescription,
			},
			"client_key": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesClientKeyDescription,
				MarkdownDescription: KubernetesClientKeyMarkdownDescription,
			},
			"namespace": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesNamespaceDescription,
				MarkdownDescription: KubernetesNamespaceMarkdownDescription,
			},
			"pod": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesPodDescription,
				MarkdownDescription: KubernetesPodMarkdownDescription,
			},
			"container": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesContainerDescription,
				MarkdownDescription: KubernetesContainerMarkdownDescription,
			},
			"command": schema.StringAttribute{
				Optional:            true,
				Description:         KubernetesCommandDescription,
				MarkdownDescription: KubernetesCommandMarkdownDescription,
			},
		},
	)
}

// ConnectionRdpSchema returns the schema attributes for the RDP protocol connection block.
func ConnectionRdpSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionCommonSchema(),
		ConnectionClipboardSchema(),
		map[string]schema.Attribute{
			"ignore_cert": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpIgnoreCertDescription,
				MarkdownDescription: RdpIgnoreCertMarkdownDescription,
			},
			"enable_full_window_drag": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableFullWindowDragDescription,
				MarkdownDescription: RdpEnableFullWindowDragMarkdownDescription,
			},
			"enable_wallpaper": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableWallpaperDescription,
				MarkdownDescription: RdpEnableWallpaperMarkdownDescription,
			},
			"enable_theming": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableThemingDescription,
				MarkdownDescription: RdpEnableThemingMarkdownDescription,
			},
			"enable_font_smoothing": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableFontSmoothingDescription,
				MarkdownDescription: RdpEnableFontSmoothingMarkdownDescription,
			},
			"enable_desktop_composition": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableDesktopCompositionDescription,
				MarkdownDescription: RdpEnableDesktopCompositionMarkdownDescription,
			},
			"enable_menu_animations": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableMenuAnimationsDescription,
				MarkdownDescription: RdpEnableMenuAnimationsMarkdownDescription,
			},
			"disable_bitmap_caching": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpDisableBitmapCachingDescription,
				MarkdownDescription: RdpDisableBitmapCachingMarkdownDescription,
			},
			"disable_offscreen_caching": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpDisableOffscreenCachingDescription,
				MarkdownDescription: RdpDisableOffscreenCachingMarkdownDescription,
			},
			"disable_glyph_caching": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpDisableGlyphCachingDescription,
				MarkdownDescription: RdpDisableGlyphCachingMarkdownDescription,
			},
			"normalize_clipboard": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("preserve"),
				Description:         RdpNormalizeClipboardDescription,
				MarkdownDescription: RdpNormalizeClipboardMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Normalize Clipboard", []string{"preserve", "unix", "windows"}, true),
				},
			},
			"security": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("any"),
				Description:         RdpSecurityDescription,
				MarkdownDescription: RdpSecurityMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Security", []string{"any", "nla", "tls", "vmconnect", "rdp"}, true),
				},
			},
			"load_balance_info": schema.StringAttribute{
				Optional:            true,
				Description:         RdpLoadBalanceInfoDescription,
				MarkdownDescription: RdpLoadBalanceInfoMarkdownDescription,
			},
			"preconnection_id": schema.StringAttribute{
				Optional:            true,
				Description:         RdpPreconnectionIdDescription,
				MarkdownDescription: RdpPreconnectionIdMarkdownDescription,
			},
			"preconnection_blob": schema.StringAttribute{
				Optional:            true,
				Description:         RdpPreconnectionBlobDescription,
				MarkdownDescription: RdpPreconnectionBlobMarkdownDescription,
			},
			"sftp": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         RdpSftpDescription,
				MarkdownDescription: RdpSftpMarkdownDescription,
				Attributes:          ConnectionRdpSftpSchema(),
				Validators:          []validator.Object{SftpUserUidRequiredValidator()},
			},
			"console_audio": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpConsoleAudioDescription,
				MarkdownDescription: RdpConsoleAudioMarkdownDescription,
			},
			"disable_audio": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpDisableAudioDescription,
				MarkdownDescription: RdpDisableAudioMarkdownDescription,
			},
			"enable_audio_input": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableAudioInputDescription,
				MarkdownDescription: RdpEnableAudioInputMarkdownDescription,
			},
			"enable_printing": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnablePrintingDescription,
				MarkdownDescription: RdpEnablePrintingMarkdownDescription,
			},
			"redirected_printer_name": schema.StringAttribute{
				Optional:            true,
				Description:         RdpRedirectedPrinterNameDescription,
				MarkdownDescription: RdpRedirectedPrinterNameMarkdownDescription,
			},
			"remote_app": schema.StringAttribute{
				Optional:            true,
				Description:         RdpRemoteAppDescription,
				MarkdownDescription: RdpRemoteAppMarkdownDescription,
			},
			"remote_app_dir": schema.StringAttribute{
				Optional:            true,
				Description:         RdpRemoteAppDirDescription,
				MarkdownDescription: RdpRemoteAppDirMarkdownDescription,
			},
			"remote_app_args": schema.StringAttribute{
				Optional:            true,
				Description:         RdpRemoteAppArgsDescription,
				MarkdownDescription: RdpRemoteAppArgsMarkdownDescription,
			},
			"force_lossless": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpForceLosslessDescription,
				MarkdownDescription: RdpForceLosslessMarkdownDescription,
			},
			"dpi": schema.Int32Attribute{
				Optional:            true,
				Description:         RdpDpiDescription,
				MarkdownDescription: RdpDpiMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32NonNegativeValidator("DPI", true),
				},
			},
			"height": schema.Int32Attribute{
				Optional:            true,
				Description:         RdpHeightDescription,
				MarkdownDescription: RdpHeightMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32NonNegativeValidator("Height", true),
				},
			},
			"width": schema.Int32Attribute{
				Optional:            true,
				Description:         RdpWidthDescription,
				MarkdownDescription: RdpWidthMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32NonNegativeValidator("Width", true),
				},
			},
			"enable_touch": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpEnableTouchDescription,
				MarkdownDescription: RdpEnableTouchMarkdownDescription,
			},
			"console": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpConsoleDescription,
				MarkdownDescription: RdpConsoleMarkdownDescription,
			},
			"timezone": schema.StringAttribute{
				Optional:            true,
				Description:         RdpTimezoneDescription,
				MarkdownDescription: RdpTimezoneMarkdownDescription,
			},
			"client_name": schema.StringAttribute{
				Optional:            true,
				Description:         RdpClientNameDescription,
				MarkdownDescription: RdpClientNameMarkdownDescription,
			},
			"initial_program": schema.StringAttribute{
				Optional:            true,
				Description:         RdpInitialProgramDescription,
				MarkdownDescription: RdpInitialProgramMarkdownDescription,
			},
			"disable_auth": schema.BoolAttribute{
				Optional:            true,
				Description:         RdpDisableAuthDescription,
				MarkdownDescription: RdpDisableAuthMarkdownDescription,
			},
			"resize_method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("display-update"),
				Description:         RdpResizeMethodDescription,
				MarkdownDescription: RdpResizeMethodMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Resize Method", []string{"display-update", "reconnect"}, true),
				},
			},
			"color_depth": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(8),
				Description:         RdpColorDepthDescription,
				MarkdownDescription: RdpColorDepthMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32OneOfValidator("Color Depth", []int32{8, 16, 24, 32}, true),
				},
			},
			"server_layout": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("en-us-qwerty"),
				Description:         RdpServerLayoutDescription,
				MarkdownDescription: RdpServerLayoutMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Server Layout", []string{
						"en-us-qwerty", "en-gb-qwerty", "de-de-qwertz", "fr-fr-azerty",
						"fr-ch-qwertz", "it-it-qwerty", "ja-jp-qwerty", "pt-br-qwerty",
						"es-es-qwerty", "sv-se-qwerty", "tr-tr-qwerty", "failsafe",
					}, true),
				},
			},
		},
	)
}

// ConnectionRdpSftpSchema returns the schema attributes for the SFTP nested block inside RDP.
func ConnectionRdpSftpSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable_sftp": schema.BoolAttribute{
			Optional:            true,
			Description:         RdpSftpEnableDescription,
			MarkdownDescription: RdpSftpEnableMarkdownDescription,
		},
		"sftp_resource_uid": schema.StringAttribute{
			Optional:            true,
			Description:         RdpSftpResourceUidDescription,
			MarkdownDescription: RdpSftpResourceUidMarkdownDescription,
		},
		"sftp_user_uid": schema.StringAttribute{
			Optional:            true,
			Description:         RdpSftpUserUidDescription,
			MarkdownDescription: RdpSftpUserUidMarkdownDescription,
		},
		"sftp_directory": schema.StringAttribute{
			Optional:            true,
			Description:         RdpSftpDirectoryDescription,
			MarkdownDescription: RdpSftpDirectoryMarkdownDescription,
		},
		"sftp_server_alive_interval": schema.Int32Attribute{
			Optional:            true,
			Description:         RdpSftpServerAliveIntervalDescription,
			MarkdownDescription: RdpSftpServerAliveIntervalMarkdownDescription,
			Validators: []validator.Int32{
				utils.Int32NonNegativeValidator("SFTP Server Alive Interval", true),
			},
		},
	}
}

// ConnectionDatabaseSchema returns the schema attributes shared by the
// mysql, postgresql, and sql_server protocol connection blocks.
func ConnectionDatabaseSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionCommonSchema(),
		ConnectionTerminalSchema(),
		ConnectionClipboardSchema(),
		map[string]schema.Attribute{
			"disable_csv_export": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DatabaseDisableCsvExportDescription,
				MarkdownDescription: DatabaseDisableCsvExportMarkdownDescription,
			},
			"disable_csv_import": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DatabaseDisableCsvImportDescription,
				MarkdownDescription: DatabaseDisableCsvImportMarkdownDescription,
			},
			"database": schema.StringAttribute{
				Optional:            true,
				Description:         DatabaseDatabaseDescription,
				MarkdownDescription: DatabaseDatabaseMarkdownDescription,
			},
		},
	)
}

// ConnectionSshSchema returns the schema attributes for the SSH protocol connection block.
func ConnectionSshSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionCommonSchema(),
		ConnectionTerminalSchema(),
		ConnectionClipboardSchema(),
		map[string]schema.Attribute{
			"host_key": schema.StringAttribute{
				Optional:            true,
				Description:         SshHostKeyDescription,
				MarkdownDescription: SshHostKeyMarkdownDescription,
			},
			"command": schema.StringAttribute{
				Optional:            true,
				Description:         SshCommandDescription,
				MarkdownDescription: SshCommandMarkdownDescription,
			},
			"locale": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("$LANG"),
				Description:         SshLocaleDescription,
				MarkdownDescription: SshLocaleMarkdownDescription,
				Validators: []validator.String{
					LocaleValidator(),
				},
			},
			"timezone": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("$TZ"),
				Description:         SshTimezoneDescription,
				MarkdownDescription: SshTimezoneMarkdownDescription,
				Validators: []validator.String{
					TimezoneValidator(),
				},
			},
			"server_alive_interval": schema.Int32Attribute{
				Optional:            true,
				Description:         SshServerAliveIntervalDescription,
				MarkdownDescription: SshServerAliveIntervalMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32NonNegativeValidator("Server Alive Interval", true),
				},
			},
			"backspace": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("127"),
				Description:         SshBackspaceDescription,
				MarkdownDescription: SshBackspaceMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Backspace", []string{"127", "8"}, true),
				},
			},
			"terminal_type": schema.StringAttribute{
				Optional:            true,
				Description:         SshTerminalTypeDescription,
				MarkdownDescription: SshTerminalTypeMarkdownDescription,
			},
			"sftp": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         SshSftpDescription,
				MarkdownDescription: SshSftpMarkdownDescription,
				Attributes:          ConnectionSshSftpSchema(),
			},
		},
	)
}

// ConnectionSshSftpSchema returns the schema attributes for the SFTP nested block inside SSH.
func ConnectionSshSftpSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable_sftp": schema.BoolAttribute{
			Optional:            true,
			Description:         SshSftpEnableDescription,
			MarkdownDescription: SshSftpEnableMarkdownDescription,
		},
	}
}

// CommonPamSettingsSchema returns the reusable schema attributes for
// the pam_settings block used across pamMachine, pamDatabase, pamDirectory, etc.
func CommonPamSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"allow_supply_host": schema.BoolAttribute{
			Optional:            true,
			Description:         PamSettingsAllowSupplyHostDescription,
			MarkdownDescription: PamSettingsAllowSupplyHostMarkdownDescription,
		},
		"connection": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsConnectionDescription,
			MarkdownDescription: PamSettingsConnectionMarkdownDescription,
			Attributes:          CommonPamSettingsConnectionSchema(),
			Validators:          []validator.Object{ConnectionProtocolBlockValidator(), ConnectionFieldsRequireEnabledValidator()},
		},
		"tunnel": schema.SingleNestedAttribute{
			Optional:            true,
			Description:         PamSettingsTunnelDescription,
			MarkdownDescription: PamSettingsTunnelMarkdownDescription,
			Attributes:          CommonPamSettingsTunnelSchema(),
			Validators:          []validator.Object{TunnelFieldsRequireEnabledValidator()},
		},
		"configuration": schema.StringAttribute{
			Required:            true,
			Description:         PamSettingsConfigurationDescription,
			MarkdownDescription: PamSettingsConfigurationMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Configuration", 1, false),
			},
		},
		"administrative_credentials": schema.StringAttribute{
			Optional:            true,
			Description:         PamSettingsAdministrativeCredentialsDescription,
			MarkdownDescription: PamSettingsAdministrativeCredentialsMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Administrative Credentials", 1, true),
			},
		},
	}
}
