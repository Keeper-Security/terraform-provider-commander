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
			Computed:            true,
			Optional:            true,
			Description:         PamSettingsTunnelEnabledDescription,
			MarkdownDescription: PamSettingsTunnelEnabledMarkdownDescription,
			Default:             booldefault.StaticBool(false),
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

// connectionScalarAttributes returns only the scalar attributes for the
// connection block (enable, protocol, connection_port, launch_credential).
// Protocol sub-blocks are defined separately in connectionProtocolBlocks().
func connectionScalarAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable": schema.BoolAttribute{
			Computed:            true,
			Optional:            true,
			Description:         PamSettingsConnectionEnableDescription,
			MarkdownDescription: PamSettingsConnectionEnableMarkdownDescription,
			Default:             booldefault.StaticBool(false),
		},
		"protocol": schema.StringAttribute{
			Optional:            true,
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
	}
}

// connectionProtocolBlocks returns the per-protocol SingleNestedBlock entries.
func connectionProtocolBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"kubernetes": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionKubernetesDescription,
			MarkdownDescription: PamSettingsConnectionKubernetesMarkdownDescription,
			Attributes:          ConnectionKubernetesSchema(),
		},
		"mysql": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionMysqlDescription,
			MarkdownDescription: PamSettingsConnectionMysqlMarkdownDescription,
			Attributes:          ConnectionDatabaseSchema(),
		},
		"postgresql": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionPostgreSqlDescription,
			MarkdownDescription: PamSettingsConnectionPostgreSqlMarkdownDescription,
			Attributes:          ConnectionDatabaseSchema(),
		},
		"rdp": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionRdpDescription,
			MarkdownDescription: PamSettingsConnectionRdpMarkdownDescription,
			Attributes:          ConnectionRdpSchema(),
			Blocks: map[string]schema.Block{
				"sftp": schema.SingleNestedBlock{
					Description:         RdpSftpDescription,
					MarkdownDescription: RdpSftpMarkdownDescription,
					Attributes:          ConnectionSftpSchema(),
					Validators:          []validator.Object{SftpUserUidRequiredValidator()},
				},
			},
		},
		"sql_server": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionSqlServerDescription,
			MarkdownDescription: PamSettingsConnectionSqlServerMarkdownDescription,
			Attributes:          ConnectionDatabaseSchema(),
		},
		"ssh": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionSshDescription,
			MarkdownDescription: PamSettingsConnectionSshMarkdownDescription,
			Attributes:          ConnectionSshSchema(),
			Blocks: map[string]schema.Block{
				"sftp": schema.SingleNestedBlock{
					Description:         SshSftpDescription,
					MarkdownDescription: SshSftpMarkdownDescription,
					Attributes:          ConnectionSshSftpSchema(),
				},
			},
		},
		"telnet": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionTelnetDescription,
			MarkdownDescription: PamSettingsConnectionTelnetMarkdownDescription,
			Attributes:          ConnectionTelnetSchema(),
		},
		"vnc": schema.SingleNestedBlock{
			Description:         PamSettingsConnectionVncDescription,
			MarkdownDescription: PamSettingsConnectionVncMarkdownDescription,
			Attributes:          ConnectionVncSchema(),
			Blocks: map[string]schema.Block{
				"sftp": schema.SingleNestedBlock{
					Description:         VncSftpDescription,
					MarkdownDescription: VncSftpMarkdownDescription,
					Attributes:          ConnectionSftpSchema(),
					Validators:          []validator.Object{SftpUserUidRequiredValidator()},
				},
			},
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

// ConnectionCommonSchema returns the 4 shared attributes used by
// K, D, SSH, Telnet (NOT RDP, VNC — they don't support typescript_recording).
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
		"typescript_recording": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionTypescriptRecordingDescription,
			MarkdownDescription: ConnectionTypescriptRecordingMarkdownDescription,
		},
	}
}

// ConnectionTerminalSchema returns the 5 terminal-related attributes
// shared by Kubernetes, Database, SSH, and Telnet protocols.
func ConnectionTerminalSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"read_only": schema.BoolAttribute{
			Optional:            true,
			Description:         ConnectionReadOnlyDescription,
			MarkdownDescription: ConnectionReadOnlyMarkdownDescription,
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
			"backspace": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("127"),
				Description:         KubernetesBackspaceDescription,
				MarkdownDescription: KubernetesBackspaceMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Backspace", []string{"127", "8"}, true),
				},
			},
		},
	)
}

// ConnectionRdpSchema returns the schema attributes for the RDP protocol connection block.
// RDP does not support typescript_recording, so we inline the base recording fields
// instead of using ConnectionCommonSchema().
func ConnectionRdpSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionClipboardSchema(),
		map[string]schema.Attribute{
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

// ConnectionSftpSchema returns the shared SFTP nested block attributes used by RDP and VNC.
func ConnectionSftpSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable_sftp": schema.BoolAttribute{
			Optional:            true,
			Description:         SftpEnableDescription,
			MarkdownDescription: SftpEnableMarkdownDescription,
		},
		"sftp_resource_uid": schema.StringAttribute{
			Optional:            true,
			Description:         SftpResourceUidDescription,
			MarkdownDescription: SftpResourceUidMarkdownDescription,
		},
		"sftp_user_uid": schema.StringAttribute{
			Optional:            true,
			Description:         SftpUserUidDescription,
			MarkdownDescription: SftpUserUidMarkdownDescription,
		},
		"sftp_directory": schema.StringAttribute{
			Optional:            true,
			Description:         SftpDirectoryDescription,
			MarkdownDescription: SftpDirectoryMarkdownDescription,
		},
		"sftp_server_alive_interval": schema.Int32Attribute{
			Optional:            true,
			Description:         SftpServerAliveIntervalDescription,
			MarkdownDescription: SftpServerAliveIntervalMarkdownDescription,
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

// ConnectionTelnetSchema returns the schema attributes for the Telnet protocol connection block.
func ConnectionTelnetSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionCommonSchema(),
		ConnectionTerminalSchema(),
		ConnectionClipboardSchema(),
		map[string]schema.Attribute{
			"username_regex": schema.StringAttribute{
				Optional:            true,
				Description:         TelnetUsernameRegexDescription,
				MarkdownDescription: TelnetUsernameRegexMarkdownDescription,
			},
			"password_regex": schema.StringAttribute{
				Optional:            true,
				Description:         TelnetPasswordRegexDescription,
				MarkdownDescription: TelnetPasswordRegexMarkdownDescription,
			},
			"login_success_regex": schema.StringAttribute{
				Optional:            true,
				Description:         TelnetLoginSuccessRegexDescription,
				MarkdownDescription: TelnetLoginSuccessRegexMarkdownDescription,
			},
			"login_failure_regex": schema.StringAttribute{
				Optional:            true,
				Description:         TelnetLoginFailureRegexDescription,
				MarkdownDescription: TelnetLoginFailureRegexMarkdownDescription,
			},
			"backspace": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("127"),
				Description:         TelnetBackspaceDescription,
				MarkdownDescription: TelnetBackspaceMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Backspace", []string{"127", "8"}, true),
				},
			},
			"terminal_type": schema.StringAttribute{
				Optional:            true,
				Description:         TelnetTerminalTypeDescription,
				MarkdownDescription: TelnetTerminalTypeMarkdownDescription,
			},
		},
	)
}

// ConnectionVncSchema returns the schema attributes for the VNC protocol connection block.
// VNC does not support typescript_recording, so we inline the applicable
// fields instead of using ConnectionCommonSchema().
func ConnectionVncSchema() map[string]schema.Attribute {
	return mergeSchemaAttributes(
		ConnectionClipboardSchema(),
		map[string]schema.Attribute{
			"session_recording": schema.BoolAttribute{
				Optional:            true,
				Description:         ConnectionSessionRecordingDescription,
				MarkdownDescription: ConnectionSessionRecordingMarkdownDescription,
			},
			"allow_supply_user": schema.BoolAttribute{
				Optional:            true,
				Description:         ConnectionAllowSupplyUserDescription,
				MarkdownDescription: ConnectionAllowSupplyUserMarkdownDescription,
			},
			"recording_include_keys": schema.BoolAttribute{
				Optional:            true,
				Description:         ConnectionRecordingIncludeKeysDescription,
				MarkdownDescription: ConnectionRecordingIncludeKeysMarkdownDescription,
			},
			"read_only": schema.BoolAttribute{
				Optional:            true,
				Description:         ConnectionReadOnlyDescription,
				MarkdownDescription: ConnectionReadOnlyMarkdownDescription,
			},
			"swap_red_blue": schema.BoolAttribute{
				Optional:            true,
				Description:         VncSwapRedBlueDescription,
				MarkdownDescription: VncSwapRedBlueMarkdownDescription,
			},
			"force_lossless": schema.BoolAttribute{
				Optional:            true,
				Description:         VncForceLosslessDescription,
				MarkdownDescription: VncForceLosslessMarkdownDescription,
			},
			"enable_audio": schema.BoolAttribute{
				Optional:            true,
				Description:         VncEnableAudioDescription,
				MarkdownDescription: VncEnableAudioMarkdownDescription,
			},
			"audio_servername": schema.StringAttribute{
				Optional:            true,
				Description:         VncAudioServernameDescription,
				MarkdownDescription: VncAudioServernameMarkdownDescription,
			},
			"dest_host": schema.StringAttribute{
				Optional:            true,
				Description:         VncDestHostDescription,
				MarkdownDescription: VncDestHostMarkdownDescription,
			},
			"dest_port": schema.Int32Attribute{
				Optional:            true,
				Description:         VncDestPortDescription,
				MarkdownDescription: VncDestPortMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32NonNegativeValidator("Dest Port", true),
				},
			},
			"clipboard_encoding": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("UTF-8"),
				Description:         VncClipboardEncodingDescription,
				MarkdownDescription: VncClipboardEncodingMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Clipboard Encoding", []string{"UTF-8", "UTF-16", "ISO8859-1", "CP1252"}, true),
				},
			},
			"cursor": schema.StringAttribute{
				Optional:            true,
				Description:         VncCursorDescription,
				MarkdownDescription: VncCursorMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Cursor", []string{"local", "remote"}, true),
				},
			},
			"color_depth": schema.Int32Attribute{
				Optional:            true,
				Description:         VncColorDepthDescription,
				MarkdownDescription: VncColorDepthMarkdownDescription,
				Validators: []validator.Int32{
					utils.Int32OneOfValidator("Color Depth", []int32{8, 16, 24, 32}, true),
				},
			},
		},
	)
}

// CommonPamSettingsBlock returns the reusable SingleNestedBlock for
// the pam_settings block used across pamMachine, pamDatabase, pamDirectory, etc.
// Uses blocks (not attributes) for nested structures so Terraform strictly
// rejects unknown attribute names.
func CommonPamSettingsBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description:         PamSettingsDescription,
		MarkdownDescription: PamSettingsMarkdownDescription,
		Validators:          []validator.Object{PamSettingsRequiredFieldsValidator()},
		Attributes: map[string]schema.Attribute{
			"allow_supply_host": schema.BoolAttribute{
				Optional:            true,
				Description:         PamSettingsAllowSupplyHostDescription,
				MarkdownDescription: PamSettingsAllowSupplyHostMarkdownDescription,
			},
			"configuration": schema.StringAttribute{
				Optional:            true,
				Description:         PamSettingsConfigurationDescription,
				MarkdownDescription: PamSettingsConfigurationMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Configuration", 1, true),
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
		},
		Blocks: map[string]schema.Block{
			"connection": schema.SingleNestedBlock{
				Description:         PamSettingsConnectionDescription,
				MarkdownDescription: PamSettingsConnectionMarkdownDescription,
				Attributes:          connectionScalarAttributes(),
				Blocks:              connectionProtocolBlocks(),
				Validators:          []validator.Object{ConnectionRequiredFieldsValidator(), ConnectionProtocolBlockValidator(), ConnectionFieldsRequireEnabledValidator()},
			},
			"tunnel": schema.SingleNestedBlock{
				Description:         PamSettingsTunnelDescription,
				MarkdownDescription: PamSettingsTunnelMarkdownDescription,
				Attributes:          CommonPamSettingsTunnelSchema(),
				Validators:          []validator.Object{TunnelRequiredFieldsValidator(), TunnelFieldsRequireEnabledValidator(), TunnelLocalPortRequiredValidator()},
			},
		},
	}
}
