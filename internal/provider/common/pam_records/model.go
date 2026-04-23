package pamrecords

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CommonPamRecordsResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Title  types.String `tfsdk:"title"`
	Notes  types.String `tfsdk:"notes"`
	Folder types.String `tfsdk:"folder"`
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
type ConnectionCommonFields struct {
	SessionRecording     types.Bool `tfsdk:"session_recording"`
	RecordingIncludeKeys types.Bool `tfsdk:"recording_include_keys"`
	AllowSupplyUser      types.Bool `tfsdk:"allow_supply_user"`
	ReadOnly             types.Bool `tfsdk:"read_only"`
}

type ConnectionTerminalFields struct {
	TypescriptRecording types.Bool   `tfsdk:"typescript_recording"`
	ColorScheme         types.String `tfsdk:"color_scheme"`
	FontName            types.String `tfsdk:"font_name"`
	FontSize            types.Int32  `tfsdk:"font_size"`
	Scrollback          types.Int32  `tfsdk:"scrollback"`
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
	ConnectionCommonFields
	ConnectionClipboardFields
	IgnoreCert               types.Bool              `tfsdk:"ignore_cert"`
	EnableFullWindowDrag     types.Bool              `tfsdk:"enable_full_window_drag"`
	EnableWallpaper          types.Bool              `tfsdk:"enable_wallpaper"`
	EnableTheming            types.Bool              `tfsdk:"enable_theming"`
	EnableFontSmoothing      types.Bool              `tfsdk:"enable_font_smoothing"`
	EnableDesktopComposition types.Bool              `tfsdk:"enable_desktop_composition"`
	EnableMenuAnimations     types.Bool              `tfsdk:"enable_menu_animations"`
	DisableBitmapCaching     types.Bool              `tfsdk:"disable_bitmap_caching"`
	DisableOffscreenCaching  types.Bool              `tfsdk:"disable_offscreen_caching"`
	DisableGlyphCaching      types.Bool              `tfsdk:"disable_glyph_caching"`
	ConsoleAudio             types.Bool              `tfsdk:"console_audio"`
	DisableAudio             types.Bool              `tfsdk:"disable_audio"`
	EnableAudioInput         types.Bool              `tfsdk:"enable_audio_input"`
	EnablePrinting           types.Bool              `tfsdk:"enable_printing"`
	ForceLossless            types.Bool              `tfsdk:"force_lossless"`
	EnableTouch              types.Bool              `tfsdk:"enable_touch"`
	Console                  types.Bool              `tfsdk:"console"`
	DisableAuth              types.Bool              `tfsdk:"disable_auth"`
	NormalizeClipboard       types.String            `tfsdk:"normalize_clipboard"`
	Security                 types.String            `tfsdk:"security"`
	LoadBalanceInfo          types.String            `tfsdk:"load_balance_info"`
	PreconnectionId          types.String            `tfsdk:"preconnection_id"`
	PreconnectionBlob        types.String            `tfsdk:"preconnection_blob"`
	RedirectedPrinterName    types.String            `tfsdk:"redirected_printer_name"`
	RemoteApp                types.String            `tfsdk:"remote_app"`
	RemoteAppDir             types.String            `tfsdk:"remote_app_dir"`
	RemoteAppArgs            types.String            `tfsdk:"remote_app_args"`
	Timezone                 types.String            `tfsdk:"timezone"`
	ClientName               types.String            `tfsdk:"client_name"`
	InitialProgram           types.String            `tfsdk:"initial_program"`
	ResizeMethod             types.String            `tfsdk:"resize_method"`
	ColorDepth               types.Int32             `tfsdk:"color_depth"`
	ServerLayout             types.String            `tfsdk:"server_layout"`
	Dpi                      types.Int32             `tfsdk:"dpi"`
	Height                   types.Int32             `tfsdk:"height"`
	Width                    types.Int32             `tfsdk:"width"`
	Sftp                     *ConnectionRdpSftpModel `tfsdk:"sftp"`
}

type ConnectionRdpSftpModel struct {
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
	HostKey              types.String           `tfsdk:"host_key"`
	Command              types.String           `tfsdk:"command"`
	Locale               types.String           `tfsdk:"locale"`
	Timezone             types.String           `tfsdk:"timezone"`
	ServerAliveInterval  types.Int32            `tfsdk:"server_alive_interval"`
	Backspace            types.String           `tfsdk:"backspace"`
	TerminalType         types.String           `tfsdk:"terminal_type"`
	Sftp                 *ConnectionSshSftpModel `tfsdk:"sftp"`
}

type ConnectionSshSftpModel struct {
	EnableSftp types.Bool `tfsdk:"enable_sftp"`
}

type ConnectionTelnetModel struct{}
type ConnectionVncModel struct{}

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
