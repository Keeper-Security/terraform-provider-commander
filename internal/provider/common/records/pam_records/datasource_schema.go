// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// CommonPamSettingsDataSourceAttribute returns the reusable pam_settings
// SingleNestedAttribute for data sources across pamMachine, pamDirectory, pamDatabase, etc.
// allowedProtocols restricts which connection protocol sub-attributes are exposed
// (MachineDirectoryProtocols vs DatabaseProtocols).
func CommonPamSettingsDataSourceAttribute(allowedProtocols []string) dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         PamSettingsDescription,
		MarkdownDescription: PamSettingsMarkdownDescription,
		Attributes:          PamSettingsDataSourceAttributes(allowedProtocols),
	}
}

func PamSettingsDataSourceAttributes(allowedProtocols []string) map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"allow_supply_host":          dschema.BoolAttribute{Computed: true},
		"configuration":              dschema.StringAttribute{Computed: true},
		"administrative_credentials": dschema.StringAttribute{Computed: true},
		"connection": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: ConnectionDataSourceAttributes(allowedProtocols),
		},
		"tunnel": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: TunnelDataSourceAttributes(),
		},
	}
}

func TunnelDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"enable":                   dschema.BoolAttribute{Computed: true},
		"remote_target_port":       dschema.Int32Attribute{Computed: true},
		"re_use_port":              dschema.BoolAttribute{Computed: true},
		"use_specified_local_port": dschema.BoolAttribute{Computed: true},
		"local_port":               dschema.Int32Attribute{Computed: true},
	}
}

func ConnectionDataSourceAttributes(allowedProtocols []string) map[string]dschema.Attribute {
	attrs := map[string]dschema.Attribute{
		"enable":            dschema.BoolAttribute{Computed: true},
		"protocol":          dschema.StringAttribute{Computed: true},
		"connection_port":   dschema.Int32Attribute{Computed: true},
		"launch_credential": dschema.StringAttribute{Computed: true},
	}
	for key, attr := range connectionProtocolDataSourceAttributes(allowedProtocols) {
		attrs[key] = attr
	}
	return attrs
}

func allConnectionProtocolDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"kubernetes": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: KubernetesDataSourceAttributes(),
		},
		"mysql": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: DatabaseDataSourceAttributes(),
		},
		"postgresql": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: DatabaseDataSourceAttributes(),
		},
		"rdp": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: RdpDataSourceAttributes(),
		},
		"sql_server": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: DatabaseDataSourceAttributes(),
		},
		"ssh": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: SshDataSourceAttributes(),
		},
		"telnet": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: TelnetDataSourceAttributes(),
		},
		"vnc": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: VncDataSourceAttributes(),
		},
	}
}

func connectionProtocolDataSourceAttributes(allowedProtocols []string) map[string]dschema.Attribute {
	all := allConnectionProtocolDataSourceAttributes()
	filtered := make(map[string]dschema.Attribute, len(allowedProtocols))
	for _, protocol := range allowedProtocols {
		key, ok := protocolToAttributeKey[protocol]
		if !ok {
			continue
		}
		if attr, found := all[key]; found {
			filtered[key] = attr
		}
	}
	return filtered
}

func commonFieldsDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"session_recording":      dschema.BoolAttribute{Computed: true},
		"recording_include_keys": dschema.BoolAttribute{Computed: true},
		"allow_supply_user":      dschema.BoolAttribute{Computed: true},
		"typescript_recording":   dschema.BoolAttribute{Computed: true},
	}
}

func terminalFieldsDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"read_only":    dschema.BoolAttribute{Computed: true},
		"color_scheme": dschema.StringAttribute{Computed: true},
		"font_name":    dschema.StringAttribute{Computed: true},
		"font_size":    dschema.Int32Attribute{Computed: true},
		"scrollback":   dschema.Int32Attribute{Computed: true},
	}
}

func clipboardFieldsDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"disable_copy":  dschema.BoolAttribute{Computed: true},
		"disable_paste": dschema.BoolAttribute{Computed: true},
	}
}

func sftpDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"enable_sftp":                dschema.BoolAttribute{Computed: true},
		"sftp_resource_uid":          dschema.StringAttribute{Computed: true},
		"sftp_user_uid":              dschema.StringAttribute{Computed: true},
		"sftp_directory":             dschema.StringAttribute{Computed: true},
		"sftp_server_alive_interval": dschema.Int32Attribute{Computed: true},
	}
}

func sshSftpDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"enable_sftp": dschema.BoolAttribute{Computed: true},
	}
}

func mergeDataSourceAttributes(maps ...map[string]dschema.Attribute) map[string]dschema.Attribute {
	result := map[string]dschema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func KubernetesDataSourceAttributes() map[string]dschema.Attribute {
	return mergeDataSourceAttributes(
		commonFieldsDataSourceAttributes(),
		terminalFieldsDataSourceAttributes(),
		map[string]dschema.Attribute{
			"rotate_on_termination": dschema.BoolAttribute{Computed: true},
			"use_ssl":               dschema.BoolAttribute{Computed: true},
			"ignore_cert":           dschema.BoolAttribute{Computed: true},
			"ca_cert":               dschema.StringAttribute{Computed: true},
			"client_cert":           dschema.StringAttribute{Computed: true},
			"client_key":            dschema.StringAttribute{Computed: true},
			"namespace":             dschema.StringAttribute{Computed: true},
			"pod":                   dschema.StringAttribute{Computed: true},
			"container":             dschema.StringAttribute{Computed: true},
			"command":               dschema.StringAttribute{Computed: true},
			"backspace":             dschema.StringAttribute{Computed: true},
		},
	)
}

func DatabaseDataSourceAttributes() map[string]dschema.Attribute {
	return mergeDataSourceAttributes(
		commonFieldsDataSourceAttributes(),
		terminalFieldsDataSourceAttributes(),
		clipboardFieldsDataSourceAttributes(),
		map[string]dschema.Attribute{
			"disable_csv_export": dschema.BoolAttribute{Computed: true},
			"disable_csv_import": dschema.BoolAttribute{Computed: true},
			"database":           dschema.StringAttribute{Computed: true},
		},
	)
}

func RdpDataSourceAttributes() map[string]dschema.Attribute {
	return mergeDataSourceAttributes(
		clipboardFieldsDataSourceAttributes(),
		map[string]dschema.Attribute{
			"session_recording":          dschema.BoolAttribute{Computed: true},
			"recording_include_keys":     dschema.BoolAttribute{Computed: true},
			"allow_supply_user":          dschema.BoolAttribute{Computed: true},
			"read_only":                  dschema.BoolAttribute{Computed: true},
			"ignore_cert":                dschema.BoolAttribute{Computed: true},
			"enable_full_window_drag":    dschema.BoolAttribute{Computed: true},
			"enable_wallpaper":           dschema.BoolAttribute{Computed: true},
			"enable_theming":             dschema.BoolAttribute{Computed: true},
			"enable_font_smoothing":      dschema.BoolAttribute{Computed: true},
			"enable_desktop_composition": dschema.BoolAttribute{Computed: true},
			"enable_menu_animations":     dschema.BoolAttribute{Computed: true},
			"disable_bitmap_caching":     dschema.BoolAttribute{Computed: true},
			"disable_offscreen_caching":  dschema.BoolAttribute{Computed: true},
			"disable_glyph_caching":      dschema.BoolAttribute{Computed: true},
			"normalize_clipboard":        dschema.StringAttribute{Computed: true},
			"security":                   dschema.StringAttribute{Computed: true},
			"load_balance_info":          dschema.StringAttribute{Computed: true},
			"preconnection_id":           dschema.StringAttribute{Computed: true},
			"preconnection_blob":         dschema.StringAttribute{Computed: true},
			"console_audio":              dschema.BoolAttribute{Computed: true},
			"disable_audio":              dschema.BoolAttribute{Computed: true},
			"enable_audio_input":         dschema.BoolAttribute{Computed: true},
			"enable_printing":            dschema.BoolAttribute{Computed: true},
			"redirected_printer_name":    dschema.StringAttribute{Computed: true},
			"remote_app":                 dschema.StringAttribute{Computed: true},
			"remote_app_dir":             dschema.StringAttribute{Computed: true},
			"remote_app_args":            dschema.StringAttribute{Computed: true},
			"force_lossless":             dschema.BoolAttribute{Computed: true},
			"dpi":                        dschema.Int32Attribute{Computed: true},
			"height":                     dschema.Int32Attribute{Computed: true},
			"width":                      dschema.Int32Attribute{Computed: true},
			"enable_touch":               dschema.BoolAttribute{Computed: true},
			"console":                    dschema.BoolAttribute{Computed: true},
			"timezone":                   dschema.StringAttribute{Computed: true},
			"client_name":                dschema.StringAttribute{Computed: true},
			"initial_program":            dschema.StringAttribute{Computed: true},
			"disable_auth":               dschema.BoolAttribute{Computed: true},
			"resize_method":              dschema.StringAttribute{Computed: true},
			"color_depth":                dschema.Int32Attribute{Computed: true},
			"server_layout":              dschema.StringAttribute{Computed: true},
			"drive_redirection_mode":     dschema.StringAttribute{Computed: true},
			"sftp": dschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: sftpDataSourceAttributes(),
			},
		},
	)
}

func SshDataSourceAttributes() map[string]dschema.Attribute {
	return mergeDataSourceAttributes(
		commonFieldsDataSourceAttributes(),
		terminalFieldsDataSourceAttributes(),
		clipboardFieldsDataSourceAttributes(),
		map[string]dschema.Attribute{
			"host_key":              dschema.StringAttribute{Computed: true},
			"command":               dschema.StringAttribute{Computed: true},
			"locale":                dschema.StringAttribute{Computed: true},
			"timezone":              dschema.StringAttribute{Computed: true},
			"server_alive_interval": dschema.Int32Attribute{Computed: true},
			"backspace":             dschema.StringAttribute{Computed: true},
			"terminal_type":         dschema.StringAttribute{Computed: true},
			"sftp": dschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: sshSftpDataSourceAttributes(),
			},
		},
	)
}

func TelnetDataSourceAttributes() map[string]dschema.Attribute {
	return mergeDataSourceAttributes(
		commonFieldsDataSourceAttributes(),
		terminalFieldsDataSourceAttributes(),
		clipboardFieldsDataSourceAttributes(),
		map[string]dschema.Attribute{
			"username_regex":      dschema.StringAttribute{Computed: true},
			"password_regex":      dschema.StringAttribute{Computed: true},
			"login_success_regex": dschema.StringAttribute{Computed: true},
			"login_failure_regex": dschema.StringAttribute{Computed: true},
			"backspace":           dschema.StringAttribute{Computed: true},
			"terminal_type":       dschema.StringAttribute{Computed: true},
		},
	)
}

func VncDataSourceAttributes() map[string]dschema.Attribute {
	return mergeDataSourceAttributes(
		clipboardFieldsDataSourceAttributes(),
		map[string]dschema.Attribute{
			"session_recording":      dschema.BoolAttribute{Computed: true},
			"allow_supply_user":      dschema.BoolAttribute{Computed: true},
			"recording_include_keys": dschema.BoolAttribute{Computed: true},
			"read_only":              dschema.BoolAttribute{Computed: true},
			"swap_red_blue":          dschema.BoolAttribute{Computed: true},
			"force_lossless":         dschema.BoolAttribute{Computed: true},
			"enable_audio":           dschema.BoolAttribute{Computed: true},
			"audio_servername":       dschema.StringAttribute{Computed: true},
			"dest_host":              dschema.StringAttribute{Computed: true},
			"dest_port":              dschema.Int32Attribute{Computed: true},
			"clipboard_encoding":     dschema.StringAttribute{Computed: true},
			"cursor":                 dschema.StringAttribute{Computed: true},
			"color_depth":            dschema.Int32Attribute{Computed: true},
			"sftp": dschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: sftpDataSourceAttributes(),
			},
		},
	)
}
