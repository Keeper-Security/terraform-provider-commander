// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FetchVaultRecord(ctx context.Context, apiManager *api.ApiManager, recordUID string) (*api.RequestResultResponse, error) {
	command := fmt.Sprintf("%s '%s' %s %s", utils.CmdGetRecord, recordUID, utils.FlagFormatJSON, utils.FlagIncludeDag)
	apiResp, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryFetchVaultRecordFailed)

	return apiResp, err
}

// ExtractFirstTextFieldValue finds a field with type "text" and the given label,
// then returns the first string from its value array. Returns "" if not found.
// Reusable across pamMachine, pamDatabase, pamDirectory, etc.
func ExtractFirstTextFieldValue(fields []utils.VaultRecordFieldResponse, label string) string {
	for i := range fields {
		f := &fields[i]
		if f.Type != "text" || f.Label != label {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			continue
		}
		if len(vals) > 0 {
			return strings.TrimSpace(vals[0])
		}
	}
	return ""
}

func MoveRecordFromSourceToDestination(ctx context.Context, apiManager *api.ApiManager, recordUID string, planFolderData string) error {
	src := recordUID

	dest := planFolderData
	if dest == "" {
		dest = "/"
	}

	command := fmt.Sprintf("%s '%s' '%s' %s", utils.CmdMv, src, dest, utils.FlagForce)
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryMoveRecordFailed)
	if err != nil {
		return err
	}
	return nil
}

// ExtractPamSettingsFromResponse reads pamSettings from the API response and
// populates a CommonPamSettingsFieldResourceModel.
//   - allow_supply_host comes from fields[type=pamSettings].value[0].allowSupplyHost
//   - configuration comes from dagDebug.all_edges where type="link" → head_uid
//
// Fields not available from the API are preserved from existingState.
// Returns nil if no pamSettings field is found and existingState is nil.
func ExtractPamSettingsFromResponse(rec *utils.VaultRecordGetResponse, existingState *CommonPamSettingsFieldResourceModel) *CommonPamSettingsFieldResourceModel {
	var raw *utils.PamSettingsFieldValueResponse
	for i := range rec.Fields {
		f := &rec.Fields[i]
		if f.Type != "pamSettings" {
			continue
		}
		var vals []utils.PamSettingsFieldValueResponse
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			break
		}
		if len(vals) > 0 {
			raw = &vals[0]
		}
		break
	}

	configuration := extractConfigurationFromDagDebug(rec.DagDebug)
	adminCredential := extractAdminCredential(rec.AssociatedCredentials)

	if raw == nil && configuration == "" && adminCredential == "" && existingState == nil {
		return nil
	}

	result := &CommonPamSettingsFieldResourceModel{}

	// Preserve state values for fields the API does not return.
	if existingState != nil {
		result.Connection = existingState.Connection
		result.Tunnel = existingState.Tunnel
		result.Configuration = existingState.Configuration
		result.AdministrativeCredentials = existingState.AdministrativeCredentials
	}

	if raw != nil {
		result.AllowSupplyHost = types.BoolValue(raw.AllowSupplyHost)
	} else if existingState != nil {
		result.AllowSupplyHost = existingState.AllowSupplyHost
	}

	if configuration != "" {
		result.Configuration = types.StringValue(configuration)
	}

	if adminCredential != "" {
		result.AdministrativeCredentials = types.StringValue(adminCredential)
	} else {
		result.AdministrativeCredentials = types.StringNull()
	}

	result.Tunnel = extractTunnelFromResponse(raw, rec.PamSettingsEnabled)
	result.Connection = extractConnectionFromResponse(raw, rec.PamSettingsEnabled, rec.AssociatedCredentials, rec.DagDebug, existingState)

	return result
}

// extractTunnelFromResponse builds a CommonPamSettingsTunnelResourceModel from
// the pamSettings field value (portForward) and pamSettingsEnabled (tunneling).
func extractTunnelFromResponse(raw *utils.PamSettingsFieldValueResponse, pamEnabled *utils.PamSettingsEnabledResponse) *CommonPamSettingsTunnelResourceModel {
	hasTunnelingEnabled := pamEnabled != nil && pamEnabled.Tunneling != nil
	hasPortForward := raw != nil && raw.PortForward != nil

	if !hasTunnelingEnabled && !hasPortForward {
		return nil
	}

	tunnel := &CommonPamSettingsTunnelResourceModel{}

	if hasTunnelingEnabled {
		tunnel.Enable = types.BoolValue(*pamEnabled.Tunneling)
	} else {
		tunnel.Enable = types.BoolValue(false)
	}

	if hasPortForward {
		pf := raw.PortForward

		tunnel.RemoteTargetPort = parseStringToInt32(pf.Port)

		if pf.ReusePort != nil {
			tunnel.ReUsePort = types.BoolValue(*pf.ReusePort)
		} else {
			tunnel.ReUsePort = types.BoolNull()
		}

		if pf.UseSpecifiedLocalPort != nil {
			tunnel.UseSpecifiedLocalPort = types.BoolValue(*pf.UseSpecifiedLocalPort)
		} else {
			tunnel.UseSpecifiedLocalPort = types.BoolNull()
		}

		tunnel.LocalPort = parseStringToInt32(pf.LocalPort)
	} else {
		tunnel.RemoteTargetPort = types.Int32Null()
		tunnel.ReUsePort = types.BoolNull()
		tunnel.UseSpecifiedLocalPort = types.BoolNull()
		tunnel.LocalPort = types.Int32Null()
	}

	return tunnel
}

// parseStringToInt32 converts a string value (e.g. "1313") to types.Int32.
// Returns types.Int32Null() if the string is empty or not a valid integer.
func parseStringToInt32(s string) types.Int32 {
	s = strings.TrimSpace(s)
	if s == "" {
		return types.Int32Null()
	}
	parsed, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return types.Int32Null()
	}
	return types.Int32Value(int32(parsed))
}

// normalizeConnectionJSON converts bare numeric JSON values to quoted strings
// for fields that Go response structs declare as string but the API may return
// as numbers (e.g. port, destPort, colorDepth, fontSize).
func normalizeConnectionJSON(raw json.RawMessage) json.RawMessage {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return raw
	}

	stringFields := map[string]bool{
		"port": true, "destPort": true, "colorDepth": true, "fontSize": true,
	}

	changed := false
	for key, val := range obj {
		if !stringFields[key] {
			continue
		}
		trimmed := strings.TrimSpace(string(val))
		if len(trimmed) > 0 && trimmed[0] >= '0' && trimmed[0] <= '9' {
			obj[key] = json.RawMessage(`"` + trimmed + `"`)
			changed = true
		}
	}

	if !changed {
		return raw
	}

	result, err := json.Marshal(obj)
	if err != nil {
		return raw
	}
	return result
}

// unmarshalConnectionByProtocol performs a two-pass unmarshal on the raw
// connection JSON. First it peeks at the "protocol" field, then unmarshals
// into the correct per-protocol struct.
// The raw JSON is normalized before unmarshaling to handle numeric values
// that should be strings (e.g. port, destPort, colorDepth, fontSize).
// Returns the protocol string, the base response (for shared fields like port),
// and the typed per-protocol struct as an interface{}.
func unmarshalConnectionByProtocol(raw json.RawMessage) (string, *utils.PamSettingsConnectionBaseResponse, interface{}) {
	if len(raw) == 0 {
		return "", nil, nil
	}

	normalized := normalizeConnectionJSON(raw)

	var base utils.PamSettingsConnectionBaseResponse
	if err := json.Unmarshal(normalized, &base); err != nil {
		return "", nil, nil
	}

	switch base.Protocol {
	case ConnectionProtocolKubernetes:
		var k8s utils.KubernetesConnectionResponse
		if err := json.Unmarshal(normalized, &k8s); err != nil {
			return base.Protocol, &base, nil
		}
		return base.Protocol, &base, &k8s
	case ConnectionProtocolMysql, ConnectionProtocolPostgreSql, ConnectionProtocolSqlServer:
		var db utils.DatabaseConnectionResponse
		if err := json.Unmarshal(normalized, &db); err != nil {
			return base.Protocol, &base, nil
		}
		return base.Protocol, &base, &db
	case ConnectionProtocolRdp:
		var rdp utils.RdpConnectionResponse
		if err := json.Unmarshal(normalized, &rdp); err != nil {
			return base.Protocol, &base, nil
		}
		return base.Protocol, &base, &rdp
	case ConnectionProtocolSsh:
		var ssh utils.SshConnectionResponse
		if err := json.Unmarshal(normalized, &ssh); err != nil {
			return base.Protocol, &base, nil
		}
		return base.Protocol, &base, &ssh
	case ConnectionProtocolTelnet:
		var telnet utils.TelnetConnectionResponse
		if err := json.Unmarshal(normalized, &telnet); err != nil {
			return base.Protocol, &base, nil
		}
		return base.Protocol, &base, &telnet
	case ConnectionProtocolVnc:
		var vnc utils.VncConnectionResponse
		if err := json.Unmarshal(normalized, &vnc); err != nil {
			return base.Protocol, &base, nil
		}
		return base.Protocol, &base, &vnc
	default:
		return base.Protocol, &base, nil
	}
}

// extractConnectionFromResponse builds a CommonPamSettingsConnectionResourceModel
// from the API response fields.
func extractConnectionFromResponse(
	raw *utils.PamSettingsFieldValueResponse,
	pamEnabled *utils.PamSettingsEnabledResponse,
	creds *utils.AssociatedCredentialsResponse,
	dagDebug *utils.DagDebugResponse,
	existingState *CommonPamSettingsFieldResourceModel,
) *CommonPamSettingsConnectionResourceModel {
	hasConnectionEnabled := pamEnabled != nil && pamEnabled.Connections != nil
	hasConnectionData := raw != nil && len(raw.Connection) > 0

	if !hasConnectionEnabled && !hasConnectionData && (existingState == nil || existingState.Connection == nil) {
		return nil
	}

	conn := &CommonPamSettingsConnectionResourceModel{}

	if existingState != nil && existingState.Connection != nil {
		*conn = *existingState.Connection
	}

	if hasConnectionEnabled {
		conn.Enable = types.BoolValue(*pamEnabled.Connections)
	} else if hasConnectionData {
		conn.Enable = types.BoolValue(true)
	} else {
		conn.Enable = types.BoolValue(false)
	}

	if creds != nil && creds.LaunchCredential != nil && strings.TrimSpace(*creds.LaunchCredential) != "" {
		conn.LaunchCredential = types.StringValue(strings.TrimSpace(*creds.LaunchCredential))
	} else {
		conn.LaunchCredential = types.StringNull()
	}

	if hasConnectionData {
		protocol, base, typed := unmarshalConnectionByProtocol(raw.Connection)

		if protocol != "" {
			conn.Protocol = types.StringValue(protocol)
		}

		if base != nil && base.Port != "" {
			conn.ConnectionPort = parseStringToInt32(base.Port)
		}

		switch protocol {
		case ConnectionProtocolKubernetes:
			if k8s, ok := typed.(*utils.KubernetesConnectionResponse); ok {
				conn.Kubernetes = extractKubernetesFromResponse(k8s, pamEnabled, dagDebug)
			}
		case ConnectionProtocolMysql:
			if db, ok := typed.(*utils.DatabaseConnectionResponse); ok {
				conn.Mysql = extractDatabaseConnectionFromResponse(db, pamEnabled)
			}
		case ConnectionProtocolPostgreSql:
			if db, ok := typed.(*utils.DatabaseConnectionResponse); ok {
				conn.PostgreSql = extractDatabaseConnectionFromResponse(db, pamEnabled)
			}
		case ConnectionProtocolSqlServer:
			if db, ok := typed.(*utils.DatabaseConnectionResponse); ok {
				conn.SqlServer = extractDatabaseConnectionFromResponse(db, pamEnabled)
			}
		case ConnectionProtocolRdp:
			if rdp, ok := typed.(*utils.RdpConnectionResponse); ok {
				var existingRdp *ConnectionRdpModel
				if existingState != nil && existingState.Connection != nil {
					existingRdp = existingState.Connection.Rdp
				}
				conn.Rdp = extractRdpConnectionFromResponse(rdp, pamEnabled, existingRdp)
			}
		case ConnectionProtocolSsh:
			if ssh, ok := typed.(*utils.SshConnectionResponse); ok {
				conn.Ssh = extractSshConnectionFromResponse(ssh, pamEnabled)
			}
		case ConnectionProtocolTelnet:
			if telnet, ok := typed.(*utils.TelnetConnectionResponse); ok {
				conn.Telnet = extractTelnetConnectionFromResponse(telnet, pamEnabled)
			}
		case ConnectionProtocolVnc:
			if vnc, ok := typed.(*utils.VncConnectionResponse); ok {
				var existingVnc *ConnectionVncModel
				if existingState != nil && existingState.Connection != nil {
					existingVnc = existingState.Connection.Vnc
				}
				conn.Vnc = extractVncConnectionFromResponse(vnc, pamEnabled, existingVnc)
			}
		}
	}

	return conn
}

// extractKubernetesFromResponse builds a ConnectionKubernetesModel from
// the per-protocol API struct.
func extractKubernetesFromResponse(
	k8sConn *utils.KubernetesConnectionResponse,
	pamEnabled *utils.PamSettingsEnabledResponse,
	dagDebug *utils.DagDebugResponse,
) *ConnectionKubernetesModel {
	k8s := &ConnectionKubernetesModel{}

	if pamEnabled != nil {
		k8s.SessionRecording = optionalBoolValue(pamEnabled.SessionRecording)
		k8s.TypescriptRecording = optionalBoolValue(pamEnabled.TypescriptRecording)
	} else {
		k8s.SessionRecording = types.BoolNull()
		k8s.TypescriptRecording = types.BoolNull()
	}

	if dagDebug != nil && dagDebug.VertexContent != nil {
		k8s.RotateOnTermination = types.BoolValue(dagDebug.VertexContent.RotateOnTermination)
	} else {
		k8s.RotateOnTermination = types.BoolNull()
	}

	if k8sConn != nil {
		k8s.RecordingIncludeKeys = optionalBoolValue(k8sConn.RecordingIncludeKeys)
		k8s.AllowSupplyUser = optionalBoolValue(k8sConn.AllowSupplyUser)
		k8s.UseSSL = optionalBoolValue(k8sConn.UseSSL)
		k8s.IgnoreCert = optionalBoolValue(k8sConn.IgnoreCert)
		k8s.ReadOnly = optionalBoolValue(k8sConn.ReadOnly)
		k8s.CaCert = setStringOrNull(k8sConn.CaCert)
		k8s.ClientCert = setStringOrNull(k8sConn.ClientCert)
		k8s.ClientKey = setStringOrNull(k8sConn.ClientKey)
		k8s.Namespace = setStringOrNull(k8sConn.Namespace)
		k8s.Pod = setStringOrNull(k8sConn.Pod)
		k8s.Container = setStringOrNull(k8sConn.Container)
		k8s.Command = setStringOrNull(k8sConn.Command)
		k8s.ColorScheme = setStringOrNull(k8sConn.ColorScheme)
		k8s.FontName = setStringOrNull(k8sConn.FontName)
		k8s.FontSize = parseStringToInt32(k8sConn.FontSize)
		k8s.Scrollback = types.Int32Value(int32(k8sConn.Scrollback))
		k8s.Backspace = setStringOrNull(k8sConn.Backspace)
	}

	return k8s
}

// extractDatabaseConnectionFromResponse builds a ConnectionDatabaseModel from
// the DatabaseConnectionResponse API struct. Used for mysql, postgresql, sql-server.
func extractDatabaseConnectionFromResponse(dbConn *utils.DatabaseConnectionResponse, pamEnabled *utils.PamSettingsEnabledResponse) *ConnectionDatabaseModel {
	db := &ConnectionDatabaseModel{}

	if pamEnabled != nil {
		db.SessionRecording = optionalBoolValue(pamEnabled.SessionRecording)
		db.TypescriptRecording = optionalBoolValue(pamEnabled.TypescriptRecording)
	} else {
		db.SessionRecording = types.BoolNull()
		db.TypescriptRecording = types.BoolNull()
	}

	if dbConn == nil {
		return db
	}

	db.AllowSupplyUser = optionalBoolValue(dbConn.AllowSupplyUser)
	db.RecordingIncludeKeys = optionalBoolValue(dbConn.RecordingIncludeKeys)
	db.ReadOnly = optionalBoolValue(dbConn.ReadOnly)
	db.DisableCopy = optionalBoolValueWithDefault(dbConn.DisableCopy)
	db.DisablePaste = optionalBoolValueWithDefault(dbConn.DisablePaste)
	db.DisableCsvExport = optionalBoolValueWithDefault(dbConn.DisableCsvExport)
	db.DisableCsvImport = optionalBoolValueWithDefault(dbConn.DisableCsvImport)
	db.Database = setStringOrNull(dbConn.Database)
	db.ColorScheme = setStringOrNull(dbConn.ColorScheme)
	db.FontName = setStringOrNull(dbConn.FontName)
	db.FontSize = parseStringToInt32(dbConn.FontSize)
	db.Scrollback = types.Int32Value(int32(dbConn.Scrollback))

	return db
}

// extractRdpConnectionFromResponse builds a ConnectionRdpModel from the API response.
func extractRdpConnectionFromResponse(rdpConn *utils.RdpConnectionResponse, pamEnabled *utils.PamSettingsEnabledResponse, existingRdp *ConnectionRdpModel) *ConnectionRdpModel {
	rdp := &ConnectionRdpModel{}

	// TODO: ColorDepth and ServerLayout are not yet returned by the read CLI response.
	// Once the CLI supports them, replace this state-preservation with API values.
	if existingRdp != nil {
		rdp.ColorDepth = existingRdp.ColorDepth
		rdp.ServerLayout = existingRdp.ServerLayout
	}

	if pamEnabled != nil {
		rdp.SessionRecording = optionalBoolValue(pamEnabled.SessionRecording)
	} else {
		rdp.SessionRecording = types.BoolNull()
	}

	if rdpConn == nil {
		return rdp
	}

	rdp.RecordingIncludeKeys = optionalBoolValue(rdpConn.RecordingIncludeKeys)
	rdp.AllowSupplyUser = optionalBoolValue(rdpConn.AllowSupplyUser)
	rdp.IgnoreCert = optionalBoolValue(rdpConn.IgnoreCert)
	rdp.EnableFullWindowDrag = optionalBoolValue(rdpConn.EnableFullWindowDrag)
	rdp.EnableWallpaper = optionalBoolValue(rdpConn.EnableWallpaper)
	rdp.EnableTheming = optionalBoolValue(rdpConn.EnableTheming)
	rdp.EnableFontSmoothing = optionalBoolValue(rdpConn.EnableFontSmoothing)
	rdp.EnableDesktopComposition = optionalBoolValue(rdpConn.EnableDesktopComposition)
	rdp.EnableMenuAnimations = optionalBoolValue(rdpConn.EnableMenuAnimations)
	rdp.DisableBitmapCaching = optionalBoolValue(rdpConn.DisableBitmapCaching)
	rdp.DisableOffscreenCaching = optionalBoolValue(rdpConn.DisableOffscreenCaching)
	rdp.DisableGlyphCaching = optionalBoolValue(rdpConn.DisableGlyphCaching)
	rdp.ConsoleAudio = optionalBoolValue(rdpConn.ConsoleAudio)
	rdp.DisableAudio = optionalBoolValue(rdpConn.DisableAudio)
	rdp.EnableAudioInput = optionalBoolValue(rdpConn.EnableAudioInput)
	rdp.EnablePrinting = optionalBoolValue(rdpConn.EnablePrinting)
	rdp.ForceLossless = optionalBoolValue(rdpConn.ForceLossless)
	rdp.ReadOnly = optionalBoolValue(rdpConn.ReadOnly)
	rdp.EnableTouch = optionalBoolValue(rdpConn.EnableTouch)
	rdp.Console = optionalBoolValue(rdpConn.Console)
	rdp.DisableAuth = optionalBoolValue(rdpConn.DisableAuth)
	rdp.DisableCopy = optionalBoolValueWithDefault(rdpConn.DisableCopy)
	rdp.DisablePaste = optionalBoolValueWithDefault(rdpConn.DisablePaste)

	rdp.NormalizeClipboard = setStringOrNull(rdpConn.NormalizeClipboard)
	rdp.Security = setStringOrNull(rdpConn.Security)
	rdp.LoadBalanceInfo = setStringOrNull(rdpConn.LoadBalanceInfo)
	rdp.PreconnectionId = setStringOrNull(rdpConn.PreconnectionId)
	rdp.PreconnectionBlob = setStringOrNull(rdpConn.PreconnectionBlob)
	rdp.RedirectedPrinterName = setStringOrNull(rdpConn.RedirectedPrinterName)
	rdp.RemoteApp = setStringOrNull(rdpConn.RemoteApp)
	rdp.RemoteAppDir = setStringOrNull(rdpConn.RemoteAppDir)
	rdp.RemoteAppArgs = setStringOrNull(rdpConn.RemoteAppArgs)
	rdp.Timezone = setStringOrNull(rdpConn.Timezone)
	rdp.ClientName = setStringOrNull(rdpConn.ClientName)
	rdp.InitialProgram = setStringOrNull(rdpConn.InitialProgram)
	rdp.ResizeMethod = setStringOrNull(rdpConn.ResizeMethod)

	if rdpConn.Dpi > 0 {
		rdp.Dpi = types.Int32Value(int32(rdpConn.Dpi))
	} else {
		rdp.Dpi = types.Int32Null()
	}
	if rdpConn.Height > 0 {
		rdp.Height = types.Int32Value(int32(rdpConn.Height))
	} else {
		rdp.Height = types.Int32Null()
	}
	if rdpConn.Width > 0 {
		rdp.Width = types.Int32Value(int32(rdpConn.Width))
	} else {
		rdp.Width = types.Int32Null()
	}

	if rdpConn.Sftp != nil {
		rdp.Sftp = extractSftpFromResponse(rdpConn.Sftp)
	}

	return rdp
}

func extractSftpFromResponse(s *utils.SftpResponse) *ConnectionSftpModel {
	sftp := &ConnectionSftpModel{}
	sftp.EnableSftp = optionalBoolValue(s.EnableSftp)
	sftp.SftpResourceUid = setStringOrNull(s.SftpResourceUid)
	sftp.SftpUserUid = setStringOrNull(s.SftpUserUid)
	sftp.SftpDirectory = setStringOrNull(s.SftpDirectory)
	if s.SftpServerAliveInterval > 0 {
		sftp.SftpServerAliveInterval = types.Int32Value(int32(s.SftpServerAliveInterval))
	} else {
		sftp.SftpServerAliveInterval = types.Int32Null()
	}
	return sftp
}

func optionalBoolValueWithDefault(b *bool) types.Bool {
	if b == nil {
		return types.BoolValue(false)
	}
	return types.BoolValue(*b)
}

func optionalBoolValue(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

func setStringOrNull(val string) types.String {
	if strings.TrimSpace(val) == "" {
		return types.StringNull()
	}
	return types.StringValue(val)
}

// extractAdminCredential returns the adminCredential value from
// associatedCredentials, or "" if not set.
func extractAdminCredential(creds *utils.AssociatedCredentialsResponse) string {
	if creds == nil || creds.AdminCredential == nil {
		return ""
	}
	return strings.TrimSpace(*creds.AdminCredential)
}

// extractConfigurationFromDagDebug finds the "link" edge in dagDebug.all_edges
// and returns its head_uid, which is the configuration UID.
func extractConfigurationFromDagDebug(dagDebug *utils.DagDebugResponse) string {
	if dagDebug == nil {
		return ""
	}
	for _, edge := range dagDebug.AllEdges {
		if edge.Type == "link" {
			return strings.TrimSpace(edge.HeadUID)
		}
	}
	return ""
}

func ApplyPamSettings(ctx context.Context, apiManager *api.ApiManager, recordUID string, pamSettings *CommonPamSettingsFieldResourceModel) error {
	if pamSettings == nil {
		return nil
	}

	configuration := pamSettings.Configuration.ValueString()

	if pamSettings.Tunnel != nil {
		if err := applyPamTunnelSettings(ctx, apiManager, recordUID, configuration, pamSettings.Tunnel); err != nil {
			return err
		}
	}

	if pamSettings.Connection != nil {
		if err := runPamConnectionEditCommand(ctx, apiManager, recordUID, configuration, pamSettings.AdministrativeCredentials, pamSettings.Connection); err != nil {
			return err
		}
	}

	if pamSettings.Tunnel == nil && pamSettings.Connection == nil {
		if err := applyPamConfiguration(ctx, apiManager, recordUID, configuration, pamSettings.AdministrativeCredentials.ValueString()); err != nil {
			return err
		}
	}

	return applyPamSettingsFieldUpdate(ctx, apiManager, recordUID, pamSettings)
}

// applyPamConfiguration applies only the configuration via `pam connection edit`
// when no tunnel or connection block is present.
func applyPamConfiguration(ctx context.Context, apiManager *api.ApiManager, recordUID string, configuration string, adminCredential string) error {
	parts := []string{
		utils.CmdPamConnectionEdit,
		fmt.Sprintf("'%s'", recordUID),
		fmt.Sprintf("%s '%s'", utils.FlagConfiguration, configuration),
	}

	if adminCredential != "" {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagAdminCredential, adminCredential))
	}

	command := strings.Join(parts, " ")
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryApplyPamSettingsFailed)
	return err
}

// applyPamSettingsFieldUpdate runs a single `record-update` with the COMPLETE
// pamSettings JSON containing all top-level keys (allowSupplyHost, portForward,
// connection). The CLI replaces the entire pamSettings field value, so every
// key must be present to avoid wiping sibling data.
func applyPamSettingsFieldUpdate(ctx context.Context, apiManager *api.ApiManager, recordUID string, pamSettings *CommonPamSettingsFieldResourceModel) error {
	pamJSON := buildFullPamSettingsJSON(pamSettings)
	if pamJSON == "" {
		return nil
	}

	command := fmt.Sprintf("%s %s '%s' '%s=$JSON:%s'",
		utils.CmdRecordUpdate,
		utils.FlagRecord,
		recordUID,
		utils.FlagPamSettings,
		pamJSON,
	)
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryApplyPamConnectionFieldUpdateFailed)
	return err
}

// applyPamTunnelSettings builds and executes a `pam tunnel edit` command.
//
// CLI pattern:
//
//	pam tunnel edit '<recordUID>' --configuration '<config>'
//	  [--enable-tunneling | --disable-tunneling]
//	  [--tunneling-override-port <port> | --remove-tunneling-override-port]
func applyPamTunnelSettings(ctx context.Context, apiManager *api.ApiManager, recordUID string, configuration string, tunnel *CommonPamSettingsTunnelResourceModel) error {
	parts := []string{
		utils.CmdPamTunnelEdit,
		fmt.Sprintf("'%s'", recordUID),
		fmt.Sprintf("%s '%s'", utils.FlagConfiguration, configuration),
	}

	if !tunnel.Enable.IsNull() && !tunnel.Enable.IsUnknown() && tunnel.Enable.ValueBool() {
		parts = append(parts, utils.FlagEnableTunneling)
	} else {
		parts = append(parts, utils.FlagDisableTunneling)
	}

	if !tunnel.RemoteTargetPort.IsNull() && !tunnel.RemoteTargetPort.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %d", utils.FlagTunnelingOverridePort, tunnel.RemoteTargetPort.ValueInt32()))
	} else {
		parts = append(parts, utils.FlagRemoveTunnelingOverridePort)
	}

	command := strings.Join(parts, " ")
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryApplyPamTunnelSettingsFailed)
	return err
}

// runPamConnectionEditCommand builds and executes:
//
//	pam connection edit '<recordUID>' --configuration '<config>'
//	  --connections=on|off --connections-recording=on|off
//	  --typescript-recording=on|off --key-events=on|off
//	  --protocol=<protocol> [--launch-user <uid>] [--admin-user <uid>]
func runPamConnectionEditCommand(ctx context.Context, apiManager *api.ApiManager, recordUID string, configuration string, adminCredentials types.String, connection *CommonPamSettingsConnectionResourceModel) error {
	parts := []string{
		utils.CmdPamConnectionEdit,
		fmt.Sprintf("'%s'", recordUID),
		fmt.Sprintf("%s '%s'", utils.FlagConfiguration, configuration),
	}

	parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnections, boolToOnOff(connection.Enable)))

	if !connection.Protocol.IsNull() && !connection.Protocol.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagProtocol, connection.Protocol.ValueString()))
	}

	switch {
	case connection.Kubernetes != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.Kubernetes.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagTypescriptRecording, boolToOnOff(connection.Kubernetes.TypescriptRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.Kubernetes.RecordingIncludeKeys)))
	case connection.Mysql != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.Mysql.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagTypescriptRecording, boolToOnOff(connection.Mysql.TypescriptRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.Mysql.RecordingIncludeKeys)))
	case connection.PostgreSql != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.PostgreSql.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagTypescriptRecording, boolToOnOff(connection.PostgreSql.TypescriptRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.PostgreSql.RecordingIncludeKeys)))
	case connection.SqlServer != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.SqlServer.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagTypescriptRecording, boolToOnOff(connection.SqlServer.TypescriptRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.SqlServer.RecordingIncludeKeys)))
	case connection.Rdp != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.Rdp.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.Rdp.RecordingIncludeKeys)))
	case connection.Ssh != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.Ssh.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagTypescriptRecording, boolToOnOff(connection.Ssh.TypescriptRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.Ssh.RecordingIncludeKeys)))
	case connection.Telnet != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.Telnet.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagTypescriptRecording, boolToOnOff(connection.Telnet.TypescriptRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.Telnet.RecordingIncludeKeys)))
	case connection.Vnc != nil:
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagConnectionsRecording, boolToOnOff(connection.Vnc.SessionRecording)))
		parts = append(parts, fmt.Sprintf("%s=%s", utils.FlagKeyEvents, boolToOnOff(connection.Vnc.RecordingIncludeKeys)))
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagLaunchCredential, connection.LaunchCredential.ValueString()))
	}

	if !adminCredentials.IsNull() && !adminCredentials.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagAdminCredential, adminCredentials.ValueString()))
	}

	command := strings.Join(parts, " ")
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryApplyPamConnectionSettingsFailed)
	return err
}

func boolToOnOff(b types.Bool) string {
	if !b.IsNull() && !b.IsUnknown() && b.ValueBool() {
		return utils.ValueOn
	}
	return utils.ValueOff
}

// buildFullPamSettingsJSON constructs the complete pamSettings JSON for
// record-update. The CLI replaces the entire field value, so all top-level
// keys (allowSupplyHost, portForward, connection) must be present.
func buildFullPamSettingsJSON(pamSettings *CommonPamSettingsFieldResourceModel) string {
	if pamSettings == nil {
		return ""
	}

	root := map[string]interface{}{
		"allowSupplyHost": !pamSettings.AllowSupplyHost.IsNull() && !pamSettings.AllowSupplyHost.IsUnknown() && pamSettings.AllowSupplyHost.ValueBool(),
	}

	if pamSettings.Tunnel != nil {
		root["portForward"] = buildPortForwardMap(pamSettings.Tunnel)
	}

	if pamSettings.Connection != nil {
		connMap := buildConnectionMap(pamSettings.Connection)
		if connMap != nil {
			root["connection"] = connMap
		}
	}

	b, err := json.Marshal(root)
	if err != nil {
		return ""
	}
	return string(b)
}

func buildPortForwardMap(tunnel *CommonPamSettingsTunnelResourceModel) map[string]interface{} {
	pf := map[string]interface{}{}

	if !tunnel.RemoteTargetPort.IsNull() && !tunnel.RemoteTargetPort.IsUnknown() {
		pf["port"] = fmt.Sprintf("%d", tunnel.RemoteTargetPort.ValueInt32())
	} else {
		pf["port"] = ""
	}

	if !tunnel.ReUsePort.IsNull() && !tunnel.ReUsePort.IsUnknown() {
		pf["reusePort"] = tunnel.ReUsePort.ValueBool()
	}

	if !tunnel.UseSpecifiedLocalPort.IsNull() && !tunnel.UseSpecifiedLocalPort.IsUnknown() {
		pf["useSpecifiedLocalPort"] = tunnel.UseSpecifiedLocalPort.ValueBool()
	}

	if !tunnel.LocalPort.IsNull() && !tunnel.LocalPort.IsUnknown() {
		pf["localPort"] = fmt.Sprintf("%d", tunnel.LocalPort.ValueInt32())
	} else {
		pf["localPort"] = ""
	}

	return pf
}

func buildConnectionMap(connection *CommonPamSettingsConnectionResourceModel) map[string]interface{} {
	if connection == nil {
		return nil
	}

	protocol := connection.Protocol.ValueString()

	switch protocol {
	case ConnectionProtocolKubernetes:
		return buildKubernetesConnectionMap(connection)
	case ConnectionProtocolMysql:
		return buildDatabaseConnectionMap(connection, connection.Mysql)
	case ConnectionProtocolPostgreSql:
		return buildDatabaseConnectionMap(connection, connection.PostgreSql)
	case ConnectionProtocolSqlServer:
		return buildDatabaseConnectionMap(connection, connection.SqlServer)
	case ConnectionProtocolRdp:
		return buildRdpConnectionMap(connection)
	case ConnectionProtocolSsh:
		return buildSshConnectionMap(connection)
	case ConnectionProtocolTelnet:
		return buildTelnetConnectionMap(connection)
	case ConnectionProtocolVnc:
		return buildVncConnectionMap(connection)
	default:
		return map[string]interface{}{
			"protocol": protocol,
		}
	}
}

func buildKubernetesConnectionMap(connection *CommonPamSettingsConnectionResourceModel) map[string]interface{} {
	k8s := connection.Kubernetes
	if k8s == nil {
		return map[string]interface{}{
			"protocol": connection.Protocol.ValueString(),
		}
	}

	connMap := map[string]interface{}{
		"protocol": connection.Protocol.ValueString(),
	}

	if !connection.ConnectionPort.IsNull() && !connection.ConnectionPort.IsUnknown() {
		connMap["port"] = fmt.Sprintf("%d", connection.ConnectionPort.ValueInt32())
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		connMap["userRecords"] = []string{connection.LaunchCredential.ValueString()}
	}

	setOptionalBoolField(connMap, "allowSupplyUser", k8s.AllowSupplyUser)
	setOptionalBoolField(connMap, "recordingIncludeKeys", k8s.RecordingIncludeKeys)
	setOptionalBoolField(connMap, "useSSL", k8s.UseSSL)
	setOptionalBoolField(connMap, "ignoreCert", k8s.IgnoreCert)
	setOptionalBoolField(connMap, "readOnly", k8s.ReadOnly)
	setOptionalStringField(connMap, "caCert", k8s.CaCert)
	setOptionalStringField(connMap, "clientCert", k8s.ClientCert)
	setOptionalStringField(connMap, "clientKey", k8s.ClientKey)
	setOptionalStringField(connMap, "namespace", k8s.Namespace)
	setOptionalStringField(connMap, "pod", k8s.Pod)
	setOptionalStringField(connMap, "container", k8s.Container)
	setOptionalStringField(connMap, "command", k8s.Command)
	setOptionalStringField(connMap, "backspace", k8s.Backspace)
	setOptionalStringField(connMap, "colorScheme", k8s.ColorScheme)
	setOptionalStringField(connMap, "fontName", k8s.FontName)

	if !k8s.FontSize.IsNull() && !k8s.FontSize.IsUnknown() {
		connMap["fontSize"] = fmt.Sprintf("%d", k8s.FontSize.ValueInt32())
	}

	if !k8s.Scrollback.IsNull() && !k8s.Scrollback.IsUnknown() {
		connMap["scrollback"] = k8s.Scrollback.ValueInt32()
	}

	return connMap
}

func buildDatabaseConnectionMap(connection *CommonPamSettingsConnectionResourceModel, db *ConnectionDatabaseModel) map[string]interface{} {
	connMap := map[string]interface{}{
		"protocol": connection.Protocol.ValueString(),
	}

	if !connection.ConnectionPort.IsNull() && !connection.ConnectionPort.IsUnknown() {
		connMap["port"] = fmt.Sprintf("%d", connection.ConnectionPort.ValueInt32())
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		connMap["userRecords"] = []string{connection.LaunchCredential.ValueString()}
	}

	if db == nil {
		return connMap
	}

	setOptionalBoolField(connMap, "allowSupplyUser", db.AllowSupplyUser)
	setOptionalBoolField(connMap, "recordingIncludeKeys", db.RecordingIncludeKeys)
	setOptionalBoolField(connMap, "readOnly", db.ReadOnly)
	setOptionalBoolField(connMap, "disableCopy", db.DisableCopy)
	setOptionalBoolField(connMap, "disablePaste", db.DisablePaste)
	setOptionalBoolField(connMap, "disableCsvExport", db.DisableCsvExport)
	setOptionalBoolField(connMap, "disableCsvImport", db.DisableCsvImport)
	setOptionalStringField(connMap, "database", db.Database)
	setOptionalStringField(connMap, "colorScheme", db.ColorScheme)
	setOptionalStringField(connMap, "fontName", db.FontName)

	if !db.FontSize.IsNull() && !db.FontSize.IsUnknown() {
		connMap["fontSize"] = fmt.Sprintf("%d", db.FontSize.ValueInt32())
	}

	if !db.Scrollback.IsNull() && !db.Scrollback.IsUnknown() {
		connMap["scrollback"] = db.Scrollback.ValueInt32()
	}

	return connMap
}

func buildRdpConnectionMap(connection *CommonPamSettingsConnectionResourceModel) map[string]interface{} {
	connMap := map[string]interface{}{
		"protocol": connection.Protocol.ValueString(),
	}

	if !connection.ConnectionPort.IsNull() && !connection.ConnectionPort.IsUnknown() {
		connMap["port"] = fmt.Sprintf("%d", connection.ConnectionPort.ValueInt32())
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		connMap["userRecords"] = []string{connection.LaunchCredential.ValueString()}
	}

	rdp := connection.Rdp
	if rdp == nil {
		return connMap
	}

	setOptionalBoolField(connMap, "allowSupplyUser", rdp.AllowSupplyUser)
	setOptionalBoolField(connMap, "recordingIncludeKeys", rdp.RecordingIncludeKeys)
	setOptionalBoolField(connMap, "ignoreCert", rdp.IgnoreCert)
	setOptionalBoolField(connMap, "enableFullWindowDrag", rdp.EnableFullWindowDrag)
	setOptionalBoolField(connMap, "enableWallpaper", rdp.EnableWallpaper)
	setOptionalBoolField(connMap, "enableTheming", rdp.EnableTheming)
	setOptionalBoolField(connMap, "enableFontSmoothing", rdp.EnableFontSmoothing)
	setOptionalBoolField(connMap, "enableDesktopComposition", rdp.EnableDesktopComposition)
	setOptionalBoolField(connMap, "enableMenuAnimations", rdp.EnableMenuAnimations)
	setOptionalBoolField(connMap, "disableBitmapCaching", rdp.DisableBitmapCaching)
	setOptionalBoolField(connMap, "disableOffscreenCaching", rdp.DisableOffscreenCaching)
	setOptionalBoolField(connMap, "disableGlyphCaching", rdp.DisableGlyphCaching)
	setOptionalBoolField(connMap, "consoleAudio", rdp.ConsoleAudio)
	setOptionalBoolField(connMap, "disableAudio", rdp.DisableAudio)
	setOptionalBoolField(connMap, "enableAudioInput", rdp.EnableAudioInput)
	setOptionalBoolField(connMap, "enablePrinting", rdp.EnablePrinting)
	setOptionalBoolField(connMap, "forceLossless", rdp.ForceLossless)
	setOptionalBoolField(connMap, "readOnly", rdp.ReadOnly)
	setOptionalBoolField(connMap, "enableTouch", rdp.EnableTouch)
	setOptionalBoolField(connMap, "console", rdp.Console)
	setOptionalBoolField(connMap, "disableAuth", rdp.DisableAuth)
	setOptionalBoolField(connMap, "disableCopy", rdp.DisableCopy)
	setOptionalBoolField(connMap, "disablePaste", rdp.DisablePaste)

	setOptionalStringField(connMap, "normalizeClipboard", rdp.NormalizeClipboard)
	setOptionalStringField(connMap, "security", rdp.Security)
	setOptionalStringField(connMap, "loadBalanceInfo", rdp.LoadBalanceInfo)
	setOptionalStringField(connMap, "preconnectionId", rdp.PreconnectionId)
	setOptionalStringField(connMap, "preconnectionBlob", rdp.PreconnectionBlob)
	setOptionalStringField(connMap, "redirectedPrinterName", rdp.RedirectedPrinterName)
	setOptionalStringField(connMap, "remoteApp", rdp.RemoteApp)
	setOptionalStringField(connMap, "remoteAppDir", rdp.RemoteAppDir)
	setOptionalStringField(connMap, "remoteAppArgs", rdp.RemoteAppArgs)
	setOptionalStringField(connMap, "timezone", rdp.Timezone)
	setOptionalStringField(connMap, "clientName", rdp.ClientName)
	setOptionalStringField(connMap, "initialProgram", rdp.InitialProgram)
	setOptionalStringField(connMap, "resizeMethod", rdp.ResizeMethod)
	setOptionalStringField(connMap, "serverLayout", rdp.ServerLayout)
	setOptionalInt32Field(connMap, "colorDepth", rdp.ColorDepth)

	setOptionalInt32Field(connMap, "dpi", rdp.Dpi)
	setOptionalInt32Field(connMap, "height", rdp.Height)
	setOptionalInt32Field(connMap, "width", rdp.Width)

	if rdp.Sftp != nil {
		sftpMap := map[string]interface{}{}
		setOptionalBoolField(sftpMap, "enableSftp", rdp.Sftp.EnableSftp)
		setOptionalStringField(sftpMap, "sftpResourceUid", rdp.Sftp.SftpResourceUid)
		setOptionalStringField(sftpMap, "sftpUserUid", rdp.Sftp.SftpUserUid)
		setOptionalStringField(sftpMap, "sftpDirectory", rdp.Sftp.SftpDirectory)
		setOptionalInt32Field(sftpMap, "sftpServerAliveInterval", rdp.Sftp.SftpServerAliveInterval)
		connMap["sftp"] = sftpMap
	}

	return connMap
}

func buildSshConnectionMap(connection *CommonPamSettingsConnectionResourceModel) map[string]interface{} {
	connMap := map[string]interface{}{
		"protocol": connection.Protocol.ValueString(),
	}

	if !connection.ConnectionPort.IsNull() && !connection.ConnectionPort.IsUnknown() {
		connMap["port"] = fmt.Sprintf("%d", connection.ConnectionPort.ValueInt32())
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		connMap["userRecords"] = []string{connection.LaunchCredential.ValueString()}
	}

	ssh := connection.Ssh
	if ssh == nil {
		return connMap
	}

	setOptionalBoolField(connMap, "allowSupplyUser", ssh.AllowSupplyUser)
	setOptionalBoolField(connMap, "recordingIncludeKeys", ssh.RecordingIncludeKeys)
	setOptionalBoolField(connMap, "readOnly", ssh.ReadOnly)
	setOptionalBoolField(connMap, "disableCopy", ssh.DisableCopy)
	setOptionalBoolField(connMap, "disablePaste", ssh.DisablePaste)
	setOptionalStringField(connMap, "colorScheme", ssh.ColorScheme)
	setOptionalStringField(connMap, "fontName", ssh.FontName)
	setOptionalStringField(connMap, "hostKey", ssh.HostKey)
	setOptionalStringField(connMap, "command", ssh.Command)
	setOptionalStringField(connMap, "locale", ssh.Locale)
	setOptionalStringField(connMap, "timezone", ssh.Timezone)
	setOptionalStringField(connMap, "backspace", ssh.Backspace)
	setOptionalStringField(connMap, "terminalType", ssh.TerminalType)

	if !ssh.FontSize.IsNull() && !ssh.FontSize.IsUnknown() {
		connMap["fontSize"] = fmt.Sprintf("%d", ssh.FontSize.ValueInt32())
	}

	if !ssh.Scrollback.IsNull() && !ssh.Scrollback.IsUnknown() {
		connMap["scrollback"] = ssh.Scrollback.ValueInt32()
	}

	setOptionalInt32Field(connMap, "serverAliveInterval", ssh.ServerAliveInterval)

	if ssh.Sftp != nil {
		sftpMap := map[string]interface{}{}
		setOptionalBoolField(sftpMap, "enableSftp", ssh.Sftp.EnableSftp)
		connMap["sftp"] = sftpMap
	}

	return connMap
}

// extractSshConnectionFromResponse builds a ConnectionSshModel from the API response.
func extractSshConnectionFromResponse(sshConn *utils.SshConnectionResponse, pamEnabled *utils.PamSettingsEnabledResponse) *ConnectionSshModel {
	ssh := &ConnectionSshModel{}

	if pamEnabled != nil {
		ssh.SessionRecording = optionalBoolValue(pamEnabled.SessionRecording)
		ssh.TypescriptRecording = optionalBoolValue(pamEnabled.TypescriptRecording)
	} else {
		ssh.SessionRecording = types.BoolNull()
		ssh.TypescriptRecording = types.BoolNull()
	}

	if sshConn == nil {
		return ssh
	}

	ssh.AllowSupplyUser = optionalBoolValue(sshConn.AllowSupplyUser)
	ssh.RecordingIncludeKeys = optionalBoolValue(sshConn.RecordingIncludeKeys)
	ssh.ReadOnly = optionalBoolValue(sshConn.ReadOnly)
	ssh.DisableCopy = optionalBoolValueWithDefault(sshConn.DisableCopy)
	ssh.DisablePaste = optionalBoolValueWithDefault(sshConn.DisablePaste)
	ssh.ColorScheme = setStringOrNull(sshConn.ColorScheme)
	ssh.FontName = setStringOrNull(sshConn.FontName)
	ssh.FontSize = parseStringToInt32(sshConn.FontSize)
	if sshConn.Scrollback > 0 {
		ssh.Scrollback = types.Int32Value(int32(sshConn.Scrollback))
	} else {
		ssh.Scrollback = types.Int32Null()
	}
	ssh.HostKey = setStringOrNull(sshConn.HostKey)
	ssh.Command = setStringOrNull(sshConn.Command)
	ssh.Locale = setStringOrNull(sshConn.Locale)
	ssh.Timezone = setStringOrNull(sshConn.Timezone)
	if sshConn.ServerAliveInterval > 0 {
		ssh.ServerAliveInterval = types.Int32Value(int32(sshConn.ServerAliveInterval))
	} else {
		ssh.ServerAliveInterval = types.Int32Null()
	}
	ssh.Backspace = setStringOrNull(sshConn.Backspace)
	ssh.TerminalType = setStringOrNull(sshConn.TerminalType)

	if sshConn.Sftp != nil {
		ssh.Sftp = &ConnectionSshSftpModel{
			EnableSftp: optionalBoolValue(sshConn.Sftp.EnableSftp),
		}
	}

	return ssh
}

func buildTelnetConnectionMap(connection *CommonPamSettingsConnectionResourceModel) map[string]interface{} {
	connMap := map[string]interface{}{
		"protocol": connection.Protocol.ValueString(),
	}

	if !connection.ConnectionPort.IsNull() && !connection.ConnectionPort.IsUnknown() {
		connMap["port"] = fmt.Sprintf("%d", connection.ConnectionPort.ValueInt32())
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		connMap["userRecords"] = []string{connection.LaunchCredential.ValueString()}
	}

	telnet := connection.Telnet
	if telnet == nil {
		return connMap
	}

	setOptionalBoolField(connMap, "allowSupplyUser", telnet.AllowSupplyUser)
	setOptionalBoolField(connMap, "recordingIncludeKeys", telnet.RecordingIncludeKeys)
	setOptionalBoolField(connMap, "readOnly", telnet.ReadOnly)
	setOptionalBoolField(connMap, "disableCopy", telnet.DisableCopy)
	setOptionalBoolField(connMap, "disablePaste", telnet.DisablePaste)
	setOptionalStringField(connMap, "colorScheme", telnet.ColorScheme)
	setOptionalStringField(connMap, "fontName", telnet.FontName)
	setOptionalStringField(connMap, "usernameRegex", telnet.UsernameRegex)
	setOptionalStringField(connMap, "passwordRegex", telnet.PasswordRegex)
	setOptionalStringField(connMap, "loginSuccessRegex", telnet.LoginSuccessRegex)
	setOptionalStringField(connMap, "loginFailureRegex", telnet.LoginFailureRegex)
	setOptionalStringField(connMap, "backspace", telnet.Backspace)
	setOptionalStringField(connMap, "terminalType", telnet.TerminalType)

	if !telnet.FontSize.IsNull() && !telnet.FontSize.IsUnknown() {
		connMap["fontSize"] = fmt.Sprintf("%d", telnet.FontSize.ValueInt32())
	}

	if !telnet.Scrollback.IsNull() && !telnet.Scrollback.IsUnknown() {
		connMap["scrollback"] = telnet.Scrollback.ValueInt32()
	}

	return connMap
}

// extractTelnetConnectionFromResponse builds a ConnectionTelnetModel from the API response.
func extractTelnetConnectionFromResponse(telnetConn *utils.TelnetConnectionResponse, pamEnabled *utils.PamSettingsEnabledResponse) *ConnectionTelnetModel {
	telnet := &ConnectionTelnetModel{}

	if pamEnabled != nil {
		telnet.SessionRecording = optionalBoolValue(pamEnabled.SessionRecording)
		telnet.TypescriptRecording = optionalBoolValue(pamEnabled.TypescriptRecording)
	} else {
		telnet.SessionRecording = types.BoolNull()
		telnet.TypescriptRecording = types.BoolNull()
	}

	if telnetConn == nil {
		return telnet
	}

	telnet.AllowSupplyUser = optionalBoolValue(telnetConn.AllowSupplyUser)
	telnet.RecordingIncludeKeys = optionalBoolValue(telnetConn.RecordingIncludeKeys)
	telnet.ReadOnly = optionalBoolValue(telnetConn.ReadOnly)
	telnet.DisableCopy = optionalBoolValueWithDefault(telnetConn.DisableCopy)
	telnet.DisablePaste = optionalBoolValueWithDefault(telnetConn.DisablePaste)
	telnet.ColorScheme = setStringOrNull(telnetConn.ColorScheme)
	telnet.FontName = setStringOrNull(telnetConn.FontName)
	telnet.FontSize = parseStringToInt32(telnetConn.FontSize)
	if telnetConn.Scrollback > 0 {
		telnet.Scrollback = types.Int32Value(int32(telnetConn.Scrollback))
	} else {
		telnet.Scrollback = types.Int32Null()
	}
	telnet.UsernameRegex = setStringOrNull(telnetConn.UsernameRegex)
	telnet.PasswordRegex = setStringOrNull(telnetConn.PasswordRegex)
	telnet.LoginSuccessRegex = setStringOrNull(telnetConn.LoginSuccessRegex)
	telnet.LoginFailureRegex = setStringOrNull(telnetConn.LoginFailureRegex)
	telnet.Backspace = setStringOrNull(telnetConn.Backspace)
	telnet.TerminalType = setStringOrNull(telnetConn.TerminalType)

	return telnet
}

func buildVncConnectionMap(connection *CommonPamSettingsConnectionResourceModel) map[string]interface{} {
	connMap := map[string]interface{}{
		"protocol": connection.Protocol.ValueString(),
	}

	if !connection.ConnectionPort.IsNull() && !connection.ConnectionPort.IsUnknown() {
		connMap["port"] = fmt.Sprintf("%d", connection.ConnectionPort.ValueInt32())
	}

	if !connection.LaunchCredential.IsNull() && !connection.LaunchCredential.IsUnknown() {
		connMap["userRecords"] = []string{connection.LaunchCredential.ValueString()}
	}

	vnc := connection.Vnc
	if vnc == nil {
		return connMap
	}

	setOptionalBoolField(connMap, "allowSupplyUser", vnc.AllowSupplyUser)
	setOptionalBoolField(connMap, "recordingIncludeKeys", vnc.RecordingIncludeKeys)
	setOptionalBoolField(connMap, "readOnly", vnc.ReadOnly)
	setOptionalBoolField(connMap, "disableCopy", vnc.DisableCopy)
	setOptionalBoolField(connMap, "disablePaste", vnc.DisablePaste)
	setOptionalBoolField(connMap, "swapRedBlue", vnc.SwapRedBlue)
	setOptionalBoolField(connMap, "forceLossless", vnc.ForceLossless)
	setOptionalBoolField(connMap, "enableAudio", vnc.EnableAudio)
	setOptionalStringField(connMap, "audioServername", vnc.AudioServername)
	setOptionalStringField(connMap, "destHost", vnc.DestHost)
	setOptionalStringField(connMap, "clipboardEncoding", vnc.ClipboardEncoding)
	setOptionalStringField(connMap, "cursor", vnc.Cursor)

	if !vnc.DestPort.IsNull() && !vnc.DestPort.IsUnknown() {
		connMap["destPort"] = fmt.Sprintf("%d", vnc.DestPort.ValueInt32())
	}

	setOptionalInt32Field(connMap, "colorDepth", vnc.ColorDepth)

	if vnc.Sftp != nil {
		sftpMap := map[string]interface{}{}
		setOptionalBoolField(sftpMap, "enableSftp", vnc.Sftp.EnableSftp)
		setOptionalStringField(sftpMap, "sftpResourceUid", vnc.Sftp.SftpResourceUid)
		setOptionalStringField(sftpMap, "sftpUserUid", vnc.Sftp.SftpUserUid)
		setOptionalStringField(sftpMap, "sftpDirectory", vnc.Sftp.SftpDirectory)
		setOptionalInt32Field(sftpMap, "sftpServerAliveInterval", vnc.Sftp.SftpServerAliveInterval)
		connMap["sftp"] = sftpMap
	}

	return connMap
}

// extractVncConnectionFromResponse builds a ConnectionVncModel from the API response.
func extractVncConnectionFromResponse(vncConn *utils.VncConnectionResponse, pamEnabled *utils.PamSettingsEnabledResponse, existingVnc *ConnectionVncModel) *ConnectionVncModel {
	vnc := &ConnectionVncModel{}

	// TODO: ColorDepth may not be returned by the read CLI response.
	// Preserve from existing state until confirmed.
	if existingVnc != nil {
		vnc.ColorDepth = existingVnc.ColorDepth
	}

	if pamEnabled != nil {
		vnc.SessionRecording = optionalBoolValue(pamEnabled.SessionRecording)
	} else {
		vnc.SessionRecording = types.BoolNull()
	}

	if vncConn == nil {
		return vnc
	}

	vnc.AllowSupplyUser = optionalBoolValue(vncConn.AllowSupplyUser)
	vnc.RecordingIncludeKeys = optionalBoolValue(vncConn.RecordingIncludeKeys)
	vnc.ReadOnly = optionalBoolValue(vncConn.ReadOnly)
	vnc.DisableCopy = optionalBoolValueWithDefault(vncConn.DisableCopy)
	vnc.DisablePaste = optionalBoolValueWithDefault(vncConn.DisablePaste)
	vnc.SwapRedBlue = optionalBoolValue(vncConn.SwapRedBlue)
	vnc.ForceLossless = optionalBoolValue(vncConn.ForceLossless)
	vnc.EnableAudio = optionalBoolValue(vncConn.EnableAudio)
	vnc.AudioServername = setStringOrNull(vncConn.AudioServername)
	vnc.DestHost = setStringOrNull(vncConn.DestHost)
	vnc.DestPort = parseStringToInt32(vncConn.DestPort)
	vnc.ClipboardEncoding = setStringOrNull(vncConn.ClipboardEncoding)
	vnc.Cursor = setStringOrNull(vncConn.Cursor)

	if vncConn.ColorDepth != "" {
		vnc.ColorDepth = parseStringToInt32(vncConn.ColorDepth)
	}

	if vncConn.Sftp != nil {
		vnc.Sftp = extractSftpFromResponse(vncConn.Sftp)
	}

	return vnc
}

func setOptionalInt32Field(m map[string]interface{}, key string, v types.Int32) {
	if !v.IsNull() && !v.IsUnknown() {
		m[key] = v.ValueInt32()
	}
}

func setOptionalBoolField(m map[string]interface{}, key string, v types.Bool) {
	if !v.IsNull() && !v.IsUnknown() {
		m[key] = v.ValueBool()
	}
}

func setOptionalStringField(m map[string]interface{}, key string, v types.String) {
	if !v.IsNull() && !v.IsUnknown() {
		m[key] = v.ValueString()
	}
}
