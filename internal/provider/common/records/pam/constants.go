// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

const (
	FlagPamHostname = "pamHostname"
)

// PAM Settings Connection Protocol.
const (
	ConnectionProtocolKubernetes = "kubernetes"
	ConnectionProtocolMysql      = "mysql"
	ConnectionProtocolPostgreSql = "postgresql"
	ConnectionProtocolRdp        = "rdp"
	ConnectionProtocolSqlServer  = "sql-server"
	ConnectionProtocolSsh        = "ssh"
	ConnectionProtocolTelnet     = "telnet"
	ConnectionProtocolVnc        = "vnc"
	ConnectionProtocolMariaDb    = "mariadb"
	ConnectionProtocolOracle     = "oracle"
)

const (
	// --- PAM Settings ---.
	PamSettingsDescription         = "PAM settings for the record, including connection, tunnel, and administrative options."
	PamSettingsMarkdownDescription = "**PAM settings** for the record, including connection, tunnel, and administrative options."

	PamSettingsAllowSupplyHostDescription         = "Whether the PAM record allows supplying a host at connection time. When true, hostname_or_ip must not be set. When false or unset, hostname_or_ip.hostname is required."
	PamSettingsAllowSupplyHostMarkdownDescription = "Whether the PAM record allows supplying a host at connection time. When **true**, `hostname_or_ip` must not be set. When **false** or unset, `hostname_or_ip.hostname` is required."

	PamSettingsConnectionDescription         = "Connection settings for the PAM record."
	PamSettingsConnectionMarkdownDescription = "**Connection** settings for the PAM record."

	PamSettingsConnectionEnableDescription         = "Whether the connection is enabled for this PAM record."
	PamSettingsConnectionEnableMarkdownDescription = "Whether the **connection** is enabled for this PAM record."

	PamSettingsConnectionConnectionPortDescription         = "Connection port. Only applicable when connection is enabled."
	PamSettingsConnectionConnectionPortMarkdownDescription = "**Connection port**. Only applicable when connection is enabled."

	PamSettingsConnectionLaunchCredentialDescription         = "Launch credential for the connection. Only applicable when connection is enabled."
	PamSettingsConnectionLaunchCredentialMarkdownDescription = "**Launch credential**. Only applicable when connection is enabled."

	PamSettingsConnectionKubernetesDescription         = "Kubernetes protocol-specific connection attributes."
	PamSettingsConnectionKubernetesMarkdownDescription = "**Kubernetes** protocol-specific connection attributes."

	PamSettingsConnectionMysqlDescription         = "MySQL protocol-specific connection attributes."
	PamSettingsConnectionMysqlMarkdownDescription = "**MySQL** protocol-specific connection attributes."

	PamSettingsConnectionPostgreSqlDescription         = "PostgreSQL protocol-specific connection attributes."
	PamSettingsConnectionPostgreSqlMarkdownDescription = "**PostgreSQL** protocol-specific connection attributes."

	PamSettingsConnectionRdpDescription         = "RDP protocol-specific connection attributes."
	PamSettingsConnectionRdpMarkdownDescription = "**RDP** protocol-specific connection attributes."

	PamSettingsConnectionSqlServerDescription         = "SQL Server protocol-specific connection attributes."
	PamSettingsConnectionSqlServerMarkdownDescription = "**SQL Server** protocol-specific connection attributes."

	PamSettingsConnectionSshDescription         = "SSH protocol-specific connection attributes."
	PamSettingsConnectionSshMarkdownDescription = "**SSH** protocol-specific connection attributes."

	PamSettingsConnectionTelnetDescription         = "Telnet protocol-specific connection attributes."
	PamSettingsConnectionTelnetMarkdownDescription = "**Telnet** protocol-specific connection attributes."

	PamSettingsConnectionVncDescription         = "VNC protocol-specific connection attributes."
	PamSettingsConnectionVncMarkdownDescription = "**VNC** protocol-specific connection attributes."

	PamSettingsConnectionMariaDbDescription         = "MariaDB protocol-specific connection attributes."
	PamSettingsConnectionMariaDbMarkdownDescription = "**MariaDB** protocol-specific connection attributes."

	PamSettingsConnectionOracleDescription         = "Oracle protocol-specific connection attributes."
	PamSettingsConnectionOracleMarkdownDescription = "**Oracle** protocol-specific connection attributes."

	// --- Shared Connection Attributes (ConnectionCommonFields) ---.
	ConnectionSessionRecordingDescription         = "Whether session recording is enabled."
	ConnectionSessionRecordingMarkdownDescription = "Whether **session recording** is enabled."

	ConnectionRecordingIncludeKeysDescription         = "Whether to include keystrokes in the recording."
	ConnectionRecordingIncludeKeysMarkdownDescription = "Whether to include **keystrokes** in the recording."

	ConnectionAllowSupplyUserDescription         = "Allow users to select credentials from their vault."
	ConnectionAllowSupplyUserMarkdownDescription = "Allow users to **select credentials** from their vault."

	ConnectionReadOnlyDescription         = "Whether the connection is read-only."
	ConnectionReadOnlyMarkdownDescription = "Whether the connection is **read-only**."

	// --- Shared Terminal Attributes (ConnectionTerminalFields) ---.
	ConnectionTypescriptRecordingDescription         = "Whether typescript recording is enabled."
	ConnectionTypescriptRecordingMarkdownDescription = "Whether **typescript recording** is enabled."

	ConnectionColorSchemeDescription         = "Color scheme for the connection terminal. Accepts named presets (black-white, gray-black, green-black, white-black) or a Guacamole terminal color scheme string."
	ConnectionColorSchemeMarkdownDescription = "**Color scheme** for the connection terminal. Accepts named presets (`black-white`, `gray-black`, `green-black`, `white-black`) or a Guacamole terminal color scheme string."

	ConnectionFontNameDescription         = "Font name for the connection terminal."
	ConnectionFontNameMarkdownDescription = "**Font name** for the connection terminal."

	ConnectionFontSizeDescription         = "Font size for the connection terminal."
	ConnectionFontSizeMarkdownDescription = "**Font size** for the connection terminal."

	ConnectionScrollbackDescription         = "Maximum scrollable size for the connection terminal."
	ConnectionScrollbackMarkdownDescription = "**Maximum scrollable size** for the connection terminal."

	// --- Shared Clipboard Attributes (ConnectionClipboardFields) ---.
	ConnectionDisableCopyDescription         = "Whether copy is disabled for the connection."
	ConnectionDisableCopyMarkdownDescription = "Whether **copy** is disabled for the connection."

	ConnectionDisablePasteDescription         = "Whether paste is disabled for the connection."
	ConnectionDisablePasteMarkdownDescription = "Whether **paste** is disabled for the connection."

	// --- Shared Rotate Attribute (ConnectionRotateFields) ---.
	ConnectionRotateOnTerminationDescription         = "Rotate launch credentials upon session termination."
	ConnectionRotateOnTerminationMarkdownDescription = "**Rotate** launch credentials upon session termination."

	// --- Kubernetes-only Connection Attributes ---.
	KubernetesUseSSLDescription         = "Use SSL/TLS for the Kubernetes connection."
	KubernetesUseSSLMarkdownDescription = "Use **SSL/TLS** for the Kubernetes connection."

	KubernetesIgnoreCertDescription         = "Ignore server certificate for the Kubernetes connection."
	KubernetesIgnoreCertMarkdownDescription = "**Ignore server certificate** for the Kubernetes connection."

	KubernetesCaCertDescription         = "Certificate Authority certificate for the Kubernetes connection."
	KubernetesCaCertMarkdownDescription = "**Certificate Authority certificate** for the Kubernetes connection."

	KubernetesClientCertDescription         = "Client certificate for the Kubernetes connection."
	KubernetesClientCertMarkdownDescription = "**Client certificate** for the Kubernetes connection."

	KubernetesClientKeyDescription         = "Client key for the Kubernetes connection."
	KubernetesClientKeyMarkdownDescription = "**Client key** for the Kubernetes connection."

	KubernetesNamespaceDescription         = "Container namespace for the Kubernetes connection."
	KubernetesNamespaceMarkdownDescription = "Container **namespace** for the Kubernetes connection."

	KubernetesPodDescription         = "Pod name for the Kubernetes connection."
	KubernetesPodMarkdownDescription = "**Pod name** for the Kubernetes connection."

	KubernetesContainerDescription         = "Container name for the Kubernetes connection."
	KubernetesContainerMarkdownDescription = "**Container name** for the Kubernetes connection."

	KubernetesCommandDescription         = "Execute command for the Kubernetes connection."
	KubernetesCommandMarkdownDescription = "**Execute command** for the Kubernetes connection."

	KubernetesBackspaceDescription         = "ASCII code sent by the backspace key. Must be \"127\" (DEL, default) or \"8\" (Ctrl+H)."
	KubernetesBackspaceMarkdownDescription = "ASCII code sent by the **backspace** key. Must be `127` (DEL, default) or `8` (Ctrl+H)."

	// --- Database-only Connection Attributes ---.
	DatabaseDisableCsvExportDescription         = "Whether CSV export is disabled for the database connection."
	DatabaseDisableCsvExportMarkdownDescription = "Whether **CSV export** is disabled for the database connection."

	DatabaseDisableCsvImportDescription         = "Whether CSV import is disabled for the database connection."
	DatabaseDisableCsvImportMarkdownDescription = "Whether **CSV import** is disabled for the database connection."

	DatabaseDatabaseDescription         = "Database name for the connection."
	DatabaseDatabaseMarkdownDescription = "**Database name** for the connection."

	// --- RDP-only Connection Attributes ---.
	RdpIgnoreCertDescription         = "Ignore server certificate for the RDP connection."
	RdpIgnoreCertMarkdownDescription = "**Ignore server certificate** for the RDP connection."

	RdpEnableFullWindowDragDescription         = "Enable full window drag during RDP sessions."
	RdpEnableFullWindowDragMarkdownDescription = "Enable **full window drag** during RDP sessions."

	RdpEnableWallpaperDescription         = "Enable wallpaper display during RDP sessions."
	RdpEnableWallpaperMarkdownDescription = "Enable **wallpaper display** during RDP sessions."

	RdpEnableThemingDescription         = "Enable theming during RDP sessions."
	RdpEnableThemingMarkdownDescription = "Enable **theming** during RDP sessions."

	RdpEnableFontSmoothingDescription         = "Enable font smoothing during RDP sessions."
	RdpEnableFontSmoothingMarkdownDescription = "Enable **font smoothing** during RDP sessions."

	RdpEnableDesktopCompositionDescription         = "Enable desktop composition during RDP sessions."
	RdpEnableDesktopCompositionMarkdownDescription = "Enable **desktop composition** during RDP sessions."

	RdpEnableMenuAnimationsDescription         = "Enable menu animations during RDP sessions."
	RdpEnableMenuAnimationsMarkdownDescription = "Enable **menu animations** during RDP sessions."

	RdpDisableBitmapCachingDescription         = "Disable bitmap caching for the RDP connection."
	RdpDisableBitmapCachingMarkdownDescription = "Disable **bitmap caching** for the RDP connection."

	RdpDisableOffscreenCachingDescription         = "Disable offscreen caching for the RDP connection."
	RdpDisableOffscreenCachingMarkdownDescription = "Disable **offscreen caching** for the RDP connection."

	RdpDisableGlyphCachingDescription         = "Disable glyph caching for the RDP connection."
	RdpDisableGlyphCachingMarkdownDescription = "Disable **glyph caching** for the RDP connection."

	RdpNormalizeClipboardDescription         = "Clipboard normalization mode for the RDP connection (e.g. preserve, unix, windows)."
	RdpNormalizeClipboardMarkdownDescription = "**Clipboard normalization** mode for the RDP connection (e.g. `preserve`, `unix`, `windows`)."

	RdpSecurityDescription         = "Security mode for the RDP connection (e.g. any, nla, tls, rdp)."
	RdpSecurityMarkdownDescription = "**Security mode** for the RDP connection (e.g. `any`, `nla`, `tls`, `rdp`)."

	RdpLoadBalanceInfoDescription         = "Load balance info for the RDP connection."
	RdpLoadBalanceInfoMarkdownDescription = "**Load balance info** for the RDP connection."

	RdpPreconnectionIdDescription         = "Pre-connection ID for the RDP connection."
	RdpPreconnectionIdMarkdownDescription = "**Pre-connection ID** for the RDP connection."

	RdpPreconnectionBlobDescription         = "Pre-connection blob for the RDP connection."
	RdpPreconnectionBlobMarkdownDescription = "**Pre-connection blob** for the RDP connection."

	RdpSftpDescription         = "SFTP settings for the RDP connection."
	RdpSftpMarkdownDescription = "**SFTP** settings for the RDP connection."

	// Shared SFTP field descriptions (used by RDP and VNC).
	SftpEnableDescription         = "Whether SFTP is enabled."
	SftpEnableMarkdownDescription = "Whether **SFTP** is enabled."

	SftpResourceUidDescription         = "UID of the SFTP resource record."
	SftpResourceUidMarkdownDescription = "UID of the **SFTP resource** record."

	SftpUserUidDescription         = "UID of the SFTP user record."
	SftpUserUidMarkdownDescription = "UID of the **SFTP user** record."

	SftpDirectoryDescription         = "SFTP root directory."
	SftpDirectoryMarkdownDescription = "**SFTP root directory**."

	SftpServerAliveIntervalDescription         = "SFTP server alive interval in seconds."
	SftpServerAliveIntervalMarkdownDescription = "**SFTP server alive interval** in seconds."

	RdpConsoleAudioDescription         = "Enable console audio for the RDP connection."
	RdpConsoleAudioMarkdownDescription = "Enable **console audio** for the RDP connection."

	RdpDisableAudioDescription         = "Disable audio for the RDP connection."
	RdpDisableAudioMarkdownDescription = "Disable **audio** for the RDP connection."

	RdpEnableAudioInputDescription         = "Enable audio input for the RDP connection."
	RdpEnableAudioInputMarkdownDescription = "Enable **audio input** for the RDP connection."

	RdpEnablePrintingDescription         = "Enable printing for the RDP connection."
	RdpEnablePrintingMarkdownDescription = "Enable **printing** for the RDP connection."

	RdpRedirectedPrinterNameDescription         = "Redirected printer name for the RDP connection."
	RdpRedirectedPrinterNameMarkdownDescription = "**Redirected printer name** for the RDP connection."

	RdpRemoteAppDescription         = "Remote application to launch."
	RdpRemoteAppMarkdownDescription = "**Remote application** to launch."

	RdpRemoteAppDirDescription         = "Working directory for the remote application."
	RdpRemoteAppDirMarkdownDescription = "**Working directory** for the remote application."

	RdpRemoteAppArgsDescription         = "Arguments for the remote application."
	RdpRemoteAppArgsMarkdownDescription = "**Arguments** for the remote application."

	RdpForceLosslessDescription         = "Force lossless compression for the RDP connection."
	RdpForceLosslessMarkdownDescription = "Force **lossless compression** for the RDP connection."

	RdpDpiDescription         = "DPI for the RDP connection display."
	RdpDpiMarkdownDescription = "**DPI** for the RDP connection display."

	RdpHeightDescription         = "Height in pixels for the RDP connection display."
	RdpHeightMarkdownDescription = "**Height** in pixels for the RDP connection display."

	RdpWidthDescription         = "Width in pixels for the RDP connection display."
	RdpWidthMarkdownDescription = "**Width** in pixels for the RDP connection display."

	RdpEnableTouchDescription         = "Enable touch input for the RDP connection."
	RdpEnableTouchMarkdownDescription = "Enable **touch input** for the RDP connection."

	RdpConsoleDescription         = "Use console session for the RDP connection."
	RdpConsoleMarkdownDescription = "Use **console session** for the RDP connection."

	RdpTimezoneDescription         = "Timezone for the RDP connection."
	RdpTimezoneMarkdownDescription = "**Timezone** for the RDP connection."

	RdpClientNameDescription         = "Client name for the RDP connection."
	RdpClientNameMarkdownDescription = "**Client name** for the RDP connection."

	RdpInitialProgramDescription         = "Initial program to run on RDP connect."
	RdpInitialProgramMarkdownDescription = "**Initial program** to run on RDP connect."

	RdpDisableAuthDescription         = "Disable authentication for the RDP connection."
	RdpDisableAuthMarkdownDescription = "Disable **authentication** for the RDP connection."

	RdpResizeMethodDescription         = "Resize method for the RDP connection (e.g. display-update, reconnect)."
	RdpResizeMethodMarkdownDescription = "**Resize method** for the RDP connection (e.g. `display-update`, `reconnect`)."

	RdpColorDepthDescription         = "Color depth in bits per pixel for the RDP connection (8, 16, 24, or 32)."
	RdpColorDepthMarkdownDescription = "**Color depth** in bits per pixel for the RDP connection (`8`, `16`, `24`, or `32`)."

	RdpServerLayoutDescription         = "Keyboard layout for the RDP connection (e.g. en-us-qwerty, de-de-qwertz)."
	RdpServerLayoutMarkdownDescription = "**Keyboard layout** for the RDP connection (e.g. `en-us-qwerty`, `de-de-qwertz`)."

	RdpDriveRedirectionModeDescription         = "Drive redirection mode for the RDP connection (e.g. none, user, resource)."
	RdpDriveRedirectionModeMarkdownDescription = "**Drive redirection mode** for the RDP connection (e.g. `none`, `user`, `resource`)."

	// --- SSH-only Connection Attributes ---.
	SshHostKeyDescription         = "Known host public key for the SSH connection."
	SshHostKeyMarkdownDescription = "**Known host public key** for the SSH connection."

	SshCommandDescription         = "Command to execute on the SSH server after connecting."
	SshCommandMarkdownDescription = "**Command** to execute on the SSH server after connecting."

	SshLocaleDescription         = "Locale for the SSH session. Use \"$LANG\" (default) to inherit client locale, or any valid POSIX locale string (e.g. en_US.UTF-8, fr_FR.UTF-8)."
	SshLocaleMarkdownDescription = "**Locale** for the SSH session. Use `$LANG` (default) to inherit client locale, or any valid POSIX locale string (e.g. `en_US.UTF-8`, `fr_FR.UTF-8`)."

	SshTimezoneDescription         = "Timezone for the SSH session. Use \"$TZ\" (default) to inherit client timezone, or any valid IANA timezone (e.g. America/New_York, Europe/London). See https://en.wikipedia.org/wiki/List_of_tz_database_time_zones."
	SshTimezoneMarkdownDescription = "**Timezone** for the SSH session. Use `$TZ` (default) to inherit client timezone, or any valid IANA timezone (e.g. `America/New_York`, `Europe/London`). See [IANA timezone list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)."

	SshServerAliveIntervalDescription         = "Interval in seconds between keep-alive messages sent to the SSH server."
	SshServerAliveIntervalMarkdownDescription = "Interval in seconds between **keep-alive messages** sent to the SSH server."

	SshBackspaceDescription         = "ASCII code sent by the backspace key. Must be \"127\" (DEL, default) or \"8\" (Ctrl+H)."
	SshBackspaceMarkdownDescription = "ASCII code sent by the **backspace** key. Must be `127` (DEL, default) or `8` (Ctrl+H)."

	SshTerminalTypeDescription         = "Terminal emulator type string for the SSH session (e.g. xterm, vt100)."
	SshTerminalTypeMarkdownDescription = "**Terminal emulator type** string for the SSH session (e.g. `xterm`, `vt100`)."

	SshSftpDescription         = "SFTP settings for the SSH connection."
	SshSftpMarkdownDescription = "**SFTP** settings for the SSH connection."

	SshSftpEnableDescription         = "Whether SFTP is enabled for the SSH connection."
	SshSftpEnableMarkdownDescription = "Whether **SFTP** is enabled for the SSH connection."

	// --- Telnet-only Connection Attributes ---.
	TelnetUsernameRegexDescription         = "Regular expression to detect the username prompt during Telnet login."
	TelnetUsernameRegexMarkdownDescription = "Regular expression to detect the **username prompt** during Telnet login."

	TelnetPasswordRegexDescription         = "Regular expression to detect the password prompt during Telnet login."
	TelnetPasswordRegexMarkdownDescription = "Regular expression to detect the **password prompt** during Telnet login."

	TelnetLoginSuccessRegexDescription         = "Regular expression to detect a successful login during Telnet login."
	TelnetLoginSuccessRegexMarkdownDescription = "Regular expression to detect a **successful login** during Telnet login."

	TelnetLoginFailureRegexDescription         = "Regular expression to detect a failed login during Telnet login."
	TelnetLoginFailureRegexMarkdownDescription = "Regular expression to detect a **failed login** during Telnet login."

	TelnetBackspaceDescription         = "ASCII code sent by the backspace key. Must be \"127\" (DEL, default) or \"8\" (Ctrl+H)."
	TelnetBackspaceMarkdownDescription = "ASCII code sent by the **backspace** key. Must be `127` (DEL, default) or `8` (Ctrl+H)."

	TelnetTerminalTypeDescription         = "Terminal emulator type string for the Telnet session (e.g. xterm, vt100)."
	TelnetTerminalTypeMarkdownDescription = "**Terminal emulator type** string for the Telnet session (e.g. `xterm`, `vt100`)."

	// --- VNC-only Connection Attributes ---.
	VncSwapRedBlueDescription         = "Swap red and blue color components in the VNC display."
	VncSwapRedBlueMarkdownDescription = "**Swap red and blue** color components in the VNC display."

	VncForceLosslessDescription         = "Force lossless compression for the VNC connection."
	VncForceLosslessMarkdownDescription = "Force **lossless compression** for the VNC connection."

	VncEnableAudioDescription         = "Enable audio for the VNC connection."
	VncEnableAudioMarkdownDescription = "Enable **audio** for the VNC connection."

	VncAudioServernameDescription         = "PulseAudio server name for VNC audio."
	VncAudioServernameMarkdownDescription = "**PulseAudio server name** for VNC audio."

	VncDestHostDescription         = "Destination host for the VNC connection."
	VncDestHostMarkdownDescription = "**Destination host** for the VNC connection."

	VncDestPortDescription         = "Destination port for the VNC connection."
	VncDestPortMarkdownDescription = "**Destination port** for the VNC connection."

	VncClipboardEncodingDescription         = "Clipboard encoding for the VNC connection. Must be one of: UTF-8 (default), UTF-16, ISO8859-1, CP1252."
	VncClipboardEncodingMarkdownDescription = "**Clipboard encoding** for the VNC connection. Must be one of: `UTF-8` (default), `UTF-16`, `ISO8859-1`, `CP1252`."

	VncCursorDescription         = "Cursor rendering mode for the VNC connection. Must be \"local\" or \"remote\"."
	VncCursorMarkdownDescription = "**Cursor rendering mode** for the VNC connection. Must be `local` or `remote`."

	VncColorDepthDescription         = "Color depth in bits per pixel for the VNC connection."
	VncColorDepthMarkdownDescription = "**Color depth** in bits per pixel for the VNC connection."

	VncSftpDescription         = "SFTP settings for the VNC connection."
	VncSftpMarkdownDescription = "**SFTP** settings for the VNC connection."

	PamSettingsTunnelDescription         = "Tunneling (port-forward) settings for the PAM record."
	PamSettingsTunnelMarkdownDescription = "**Tunneling** (port-forward) settings for the PAM record."

	PamSettingsConfigurationDescription         = "Configuration identifier for the PAM record."
	PamSettingsConfigurationMarkdownDescription = "**Configuration** identifier for the PAM record."

	PamSettingsAdministrativeCredentialsDescription         = "Linked PAM User credential used for connection and administrative operations."
	PamSettingsAdministrativeCredentialsMarkdownDescription = "**Linked PAM User credential** used for connection and administrative operations."

	// --- PAM Settings Tunnel ---.
	PamSettingsTunnelEnabledDescription         = "Whether tunneling is enabled for this PAM record."
	PamSettingsTunnelEnabledMarkdownDescription = "Whether **tunneling** is enabled for this PAM record."

	PamSettingsTunnelRemoteTargetPortDescription         = "Remote target port for the tunnel. Only applicable when tunneling is enabled."
	PamSettingsTunnelRemoteTargetPortMarkdownDescription = "**Remote target port** for the tunnel. Only applicable when tunneling is enabled."

	PamSettingsTunnelReUsePortDescription         = "Whether to reuse the port for tunneling. Only applicable when tunneling is enabled."
	PamSettingsTunnelReUsePortMarkdownDescription = "Whether to **reuse the port** for tunneling. Only applicable when tunneling is enabled."

	PamSettingsTunnelUseSpecifiedLocalPortDescription         = "Whether to use a specified local port for tunneling. Only applicable when tunneling is enabled."
	PamSettingsTunnelUseSpecifiedLocalPortMarkdownDescription = "Whether to use a **specified local port** for tunneling. Only applicable when tunneling is enabled."

	PamSettingsTunnelLocalPortDescription         = "Local port for tunneling. Only applicable when tunneling is enabled."
	PamSettingsTunnelLocalPortMarkdownDescription = "**Local port** for tunneling. Only applicable when tunneling is enabled."
)
