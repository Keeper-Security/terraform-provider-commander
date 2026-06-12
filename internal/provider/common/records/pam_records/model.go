// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HostnameOrIPModel maps the nested hostname_or_ip object.
type HostnameOrIPModel struct {
	HostName           types.String `tfsdk:"hostname"`
	AdministrativePort types.Int32  `tfsdk:"administrative_port"`
}

type CommonPamSettingsConnectionResourceModel struct {
	Enable           types.Bool   `tfsdk:"enable"`
	Protocol         types.String `tfsdk:"protocol"`
	ConnectionPort   types.Int32  `tfsdk:"connection_port"`
	LaunchCredential types.String `tfsdk:"launch_credential"`

	Kubernetes *ConnectionKubernetesModel `tfsdk:"kubernetes"`
	Mysql      *ConnectionDatabaseModel   `tfsdk:"mysql"`
	PostgreSql *ConnectionDatabaseModel   `tfsdk:"postgresql"`
	Rdp        *ConnectionRdpModel        `tfsdk:"rdp"`
	SqlServer  *ConnectionDatabaseModel   `tfsdk:"sql_server"`
	Ssh        *ConnectionSshModel        `tfsdk:"ssh"`
	Telnet     *ConnectionTelnetModel     `tfsdk:"telnet"`
	Vnc        *ConnectionVncModel        `tfsdk:"vnc"`
}

// Shared field groups embedded by per-protocol models.
// ConnectionCommonFields is used by K, D, SSH, Telnet (NOT RDP, VNC).
type ConnectionCommonFields struct {
	SessionRecording     types.Bool `tfsdk:"session_recording"`
	RecordingIncludeKeys types.Bool `tfsdk:"recording_include_keys"`
	AllowSupplyUser      types.Bool `tfsdk:"allow_supply_user"`
	TypescriptRecording  types.Bool `tfsdk:"typescript_recording"`
}

type ConnectionTerminalFields struct {
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	ColorScheme types.String `tfsdk:"color_scheme"`
	FontName    types.String `tfsdk:"font_name"`
	FontSize    types.Int32  `tfsdk:"font_size"`
	Scrollback  types.Int32  `tfsdk:"scrollback"`
}

type ConnectionClipboardFields struct {
	DisableCopy  types.Bool `tfsdk:"disable_copy"`
	DisablePaste types.Bool `tfsdk:"disable_paste"`
}

// Per-protocol attribute models.
type ConnectionKubernetesModel struct {
	ConnectionCommonFields
	ConnectionTerminalFields
	RotateOnTermination types.Bool   `tfsdk:"rotate_on_termination"`
	UseSSL              types.Bool   `tfsdk:"use_ssl"`
	IgnoreCert          types.Bool   `tfsdk:"ignore_cert"`
	CaCert              types.String `tfsdk:"ca_cert"`
	ClientCert          types.String `tfsdk:"client_cert"`
	ClientKey           types.String `tfsdk:"client_key"`
	Namespace           types.String `tfsdk:"namespace"`
	Pod                 types.String `tfsdk:"pod"`
	Container           types.String `tfsdk:"container"`
	Command             types.String `tfsdk:"command"`
	Backspace           types.String `tfsdk:"backspace"`
}

// ConnectionDatabaseModel is shared by mysql, postgresql, and sql_server protocols.
type ConnectionDatabaseModel struct {
	ConnectionCommonFields
	ConnectionTerminalFields
	ConnectionClipboardFields
	DisableCsvExport types.Bool   `tfsdk:"disable_csv_export"`
	DisableCsvImport types.Bool   `tfsdk:"disable_csv_import"`
	Database         types.String `tfsdk:"database"`
}

type ConnectionRdpModel struct {
	SessionRecording     types.Bool `tfsdk:"session_recording"`
	RecordingIncludeKeys types.Bool `tfsdk:"recording_include_keys"`
	AllowSupplyUser      types.Bool `tfsdk:"allow_supply_user"`
	ReadOnly             types.Bool `tfsdk:"read_only"`
	ConnectionClipboardFields
	IgnoreCert               types.Bool           `tfsdk:"ignore_cert"`
	EnableFullWindowDrag     types.Bool           `tfsdk:"enable_full_window_drag"`
	EnableWallpaper          types.Bool           `tfsdk:"enable_wallpaper"`
	EnableTheming            types.Bool           `tfsdk:"enable_theming"`
	EnableFontSmoothing      types.Bool           `tfsdk:"enable_font_smoothing"`
	EnableDesktopComposition types.Bool           `tfsdk:"enable_desktop_composition"`
	EnableMenuAnimations     types.Bool           `tfsdk:"enable_menu_animations"`
	DisableBitmapCaching     types.Bool           `tfsdk:"disable_bitmap_caching"`
	DisableOffscreenCaching  types.Bool           `tfsdk:"disable_offscreen_caching"`
	DisableGlyphCaching      types.Bool           `tfsdk:"disable_glyph_caching"`
	ConsoleAudio             types.Bool           `tfsdk:"console_audio"`
	DisableAudio             types.Bool           `tfsdk:"disable_audio"`
	EnableAudioInput         types.Bool           `tfsdk:"enable_audio_input"`
	EnablePrinting           types.Bool           `tfsdk:"enable_printing"`
	ForceLossless            types.Bool           `tfsdk:"force_lossless"`
	EnableTouch              types.Bool           `tfsdk:"enable_touch"`
	Console                  types.Bool           `tfsdk:"console"`
	DisableAuth              types.Bool           `tfsdk:"disable_auth"`
	NormalizeClipboard       types.String         `tfsdk:"normalize_clipboard"`
	Security                 types.String         `tfsdk:"security"`
	LoadBalanceInfo          types.String         `tfsdk:"load_balance_info"`
	PreconnectionId          types.String         `tfsdk:"preconnection_id"`
	PreconnectionBlob        types.String         `tfsdk:"preconnection_blob"`
	RedirectedPrinterName    types.String         `tfsdk:"redirected_printer_name"`
	RemoteApp                types.String         `tfsdk:"remote_app"`
	RemoteAppDir             types.String         `tfsdk:"remote_app_dir"`
	RemoteAppArgs            types.String         `tfsdk:"remote_app_args"`
	Timezone                 types.String         `tfsdk:"timezone"`
	ClientName               types.String         `tfsdk:"client_name"`
	InitialProgram           types.String         `tfsdk:"initial_program"`
	ResizeMethod             types.String         `tfsdk:"resize_method"`
	ColorDepth               types.Int32          `tfsdk:"color_depth"`
	ServerLayout             types.String         `tfsdk:"server_layout"`
	Dpi                      types.Int32          `tfsdk:"dpi"`
	Height                   types.Int32          `tfsdk:"height"`
	Width                    types.Int32          `tfsdk:"width"`
	Sftp                     *ConnectionSftpModel `tfsdk:"sftp"`
	DriveRedirectionMode     types.String         `tfsdk:"drive_redirection_mode"`
}

// ConnectionSftpModel is the shared SFTP nested block used by RDP and VNC.
type ConnectionSftpModel struct {
	EnableSftp              types.Bool   `tfsdk:"enable_sftp"`
	SftpResourceUid         types.String `tfsdk:"sftp_resource_uid"`
	SftpUserUid             types.String `tfsdk:"sftp_user_uid"`
	SftpDirectory           types.String `tfsdk:"sftp_directory"`
	SftpServerAliveInterval types.Int32  `tfsdk:"sftp_server_alive_interval"`
}
type ConnectionSshModel struct {
	ConnectionCommonFields
	ConnectionTerminalFields
	ConnectionClipboardFields
	HostKey             types.String            `tfsdk:"host_key"`
	Command             types.String            `tfsdk:"command"`
	Locale              types.String            `tfsdk:"locale"`
	Timezone            types.String            `tfsdk:"timezone"`
	ServerAliveInterval types.Int32             `tfsdk:"server_alive_interval"`
	Backspace           types.String            `tfsdk:"backspace"`
	TerminalType        types.String            `tfsdk:"terminal_type"`
	Sftp                *ConnectionSshSftpModel `tfsdk:"sftp"`
}

type ConnectionSshSftpModel struct {
	EnableSftp types.Bool `tfsdk:"enable_sftp"`
}

type ConnectionTelnetModel struct {
	ConnectionCommonFields
	ConnectionTerminalFields
	ConnectionClipboardFields
	UsernameRegex     types.String `tfsdk:"username_regex"`
	PasswordRegex     types.String `tfsdk:"password_regex"`
	LoginSuccessRegex types.String `tfsdk:"login_success_regex"`
	LoginFailureRegex types.String `tfsdk:"login_failure_regex"`
	Backspace         types.String `tfsdk:"backspace"`
	TerminalType      types.String `tfsdk:"terminal_type"`
}

type ConnectionVncModel struct {
	SessionRecording     types.Bool `tfsdk:"session_recording"`
	AllowSupplyUser      types.Bool `tfsdk:"allow_supply_user"`
	RecordingIncludeKeys types.Bool `tfsdk:"recording_include_keys"`
	ReadOnly             types.Bool `tfsdk:"read_only"`
	ConnectionClipboardFields
	SwapRedBlue       types.Bool           `tfsdk:"swap_red_blue"`
	ForceLossless     types.Bool           `tfsdk:"force_lossless"`
	EnableAudio       types.Bool           `tfsdk:"enable_audio"`
	AudioServername   types.String         `tfsdk:"audio_servername"`
	DestHost          types.String         `tfsdk:"dest_host"`
	DestPort          types.Int32          `tfsdk:"dest_port"`
	ClipboardEncoding types.String         `tfsdk:"clipboard_encoding"`
	Cursor            types.String         `tfsdk:"cursor"`
	ColorDepth        types.Int32          `tfsdk:"color_depth"`
	Sftp              *ConnectionSftpModel `tfsdk:"sftp"`
}

// This is structure of "portForward" that we get from the API.
type CommonPamSettingsTunnelResourceModel struct {
	Enable                types.Bool  `tfsdk:"enable"`             // use  --enable-tunneling /  --disable-tunneling to enable / disable tunneling
	RemoteTargetPort      types.Int32 `tfsdk:"remote_target_port"` // this is remote target port, --tunneling-override-port / --remove-tunneling-override-port to override / remove the override
	ReUsePort             types.Bool  `tfsdk:"re_use_port"`
	UseSpecifiedLocalPort types.Bool  `tfsdk:"use_specified_local_port"`
	LocalPort             types.Int32 `tfsdk:"local_port"`
}
type CommonPamSettingsFieldResourceModel struct {
	AllowSupplyHost           types.Bool                                `tfsdk:"allow_supply_host"`
	Connection                *CommonPamSettingsConnectionResourceModel `tfsdk:"connection"`
	Tunnel                    *CommonPamSettingsTunnelResourceModel     `tfsdk:"tunnel"`
	Configuration             types.String                              `tfsdk:"configuration"`
	AdministrativeCredentials types.String                              `tfsdk:"administrative_credentials"`
}
