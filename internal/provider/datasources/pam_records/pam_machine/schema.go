// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"

	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_machine"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *PamMachineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a PAM machine record by UID or name.",
		MarkdownDescription: "Use this data source to look up a **PAM machine** record by **UID** or **name**.",
		Attributes: map[string]dschema.Attribute{
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
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.FolderDescription,
				MarkdownDescription: commonpammachine.FolderMarkdownDescription,
			},
			"pam_settings": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         commonpamrecords.PamSettingsDescription,
				MarkdownDescription: commonpamrecords.PamSettingsMarkdownDescription,
				Attributes:          pamSettingsDataSourceAttributes(),
			},
		},
	}
}

func pamSettingsDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"allow_supply_host":          dschema.BoolAttribute{Computed: true},
		"configuration":              dschema.StringAttribute{Computed: true},
		"administrative_credentials": dschema.StringAttribute{Computed: true},
		"connection": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: connectionDataSourceAttributes(),
		},
		"tunnel": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: tunnelDataSourceAttributes(),
		},
	}
}

func tunnelDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"enable":                   dschema.BoolAttribute{Computed: true},
		"remote_target_port":       dschema.Int32Attribute{Computed: true},
		"re_use_port":              dschema.BoolAttribute{Computed: true},
		"use_specified_local_port": dschema.BoolAttribute{Computed: true},
		"local_port":               dschema.Int32Attribute{Computed: true},
	}
}

func connectionDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"enable":            dschema.BoolAttribute{Computed: true},
		"protocol":          dschema.StringAttribute{Computed: true},
		"connection_port":   dschema.Int32Attribute{Computed: true},
		"launch_credential": dschema.StringAttribute{Computed: true},
		"kubernetes": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: kubernetesDataSourceAttributes(),
		},
		"mysql": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: databaseDataSourceAttributes(),
		},
		"postgresql": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: databaseDataSourceAttributes(),
		},
		"rdp": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: rdpDataSourceAttributes(),
		},
		"sql_server": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: databaseDataSourceAttributes(),
		},
		"ssh": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: sshDataSourceAttributes(),
		},
		"telnet": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: telnetDataSourceAttributes(),
		},
		"vnc": dschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: vncDataSourceAttributes(),
		},
	}
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

func kubernetesDataSourceAttributes() map[string]dschema.Attribute {
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

func databaseDataSourceAttributes() map[string]dschema.Attribute {
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

func rdpDataSourceAttributes() map[string]dschema.Attribute {
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
			"sftp": dschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: sftpDataSourceAttributes(),
			},
		},
	)
}

func sshDataSourceAttributes() map[string]dschema.Attribute {
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

func telnetDataSourceAttributes() map[string]dschema.Attribute {
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

func vncDataSourceAttributes() map[string]dschema.Attribute {
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
