// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func boolPtr(b bool) *bool { return &b }

func TestExtractPamSettingsFromResponse_NoPhantomConnectionOrTunnelWhenOnlyEnableFlags(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamSettings",
				Value: json.RawMessage(`[{"allowSupplyHost":false}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
			Tunneling:   boolPtr(false),
		},
		DagDebug: &utils.DagDebugResponse{
			AllEdges: []utils.DagDebugEdgeResponse{
				{Type: "link", HeadUID: "config-uid"},
			},
		},
	}

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil {
		t.Fatal("expected pam_settings")
	}
	if settings.Connection != nil {
		t.Fatalf("expected nil connection when API only reports enable flags, got enable=%v", settings.Connection.Enable)
	}
	if settings.Tunnel != nil {
		t.Fatalf("expected nil tunnel when API only reports enable flags, got enable=%v", settings.Tunnel.Enable)
	}
	if settings.Configuration.ValueString() != "config-uid" {
		t.Fatalf("configuration = %q, want config-uid", settings.Configuration.ValueString())
	}
}

func TestExtractPamSettingsFromResponse_PreservesExistingConnectionAndTunnel(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamSettings",
				Value: json.RawMessage(`[{"allowSupplyHost":false}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
			Tunneling:   boolPtr(false),
		},
		DagDebug: &utils.DagDebugResponse{
			AllEdges: []utils.DagDebugEdgeResponse{
				{Type: "link", HeadUID: "config-uid"},
			},
		},
	}

	conn := &commonpamrecords.CommonPamSettingsConnectionResourceModel{}
	conn.Protocol = types.StringValue(commonpamrecords.ConnectionProtocolSsh)

	existing := &commonpamrecords.CommonPamSettingsFieldResourceModel{
		Connection: conn,
	}
	existing.Tunnel = &commonpamrecords.CommonPamSettingsTunnelResourceModel{
		Enable: types.BoolValue(true),
	}

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, existing)
	if settings == nil {
		t.Fatal("expected pam_settings")
	}
	if settings.Connection == nil || settings.Connection.Protocol.ValueString() != commonpamrecords.ConnectionProtocolSsh {
		t.Fatalf("expected existing connection preserved, got %+v", settings.Connection)
	}
	if settings.Tunnel == nil || !settings.Tunnel.Enable.ValueBool() {
		t.Fatalf("expected existing tunnel preserved, got %+v", settings.Tunnel)
	}
}

func sshPamSettingsRecord(t *testing.T, vertexContent json.RawMessage) *utils.VaultRecordGetResponse {
	t.Helper()
	rec := &utils.VaultRecordGetResponse{
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type: "pamSettings",
				Value: json.RawMessage(`[{
					"allowSupplyHost": false,
					"connection": {
						"protocol": "ssh",
						"port": "22"
					}
				}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
		},
		DagDebug: &utils.DagDebugResponse{
			AllEdges: []utils.DagDebugEdgeResponse{
				{Type: "link", HeadUID: "config-uid"},
			},
		},
	}
	if len(vertexContent) > 0 {
		var content utils.DagDebugVertexContentResponse
		if err := json.Unmarshal(vertexContent, &content); err != nil {
			t.Fatalf("vertex content json: %v", err)
		}
		rec.DagDebug.VertexContent = &content
	}
	return rec
}

func TestExtractPamSettingsFromResponse_SshRotateOnTerminationNullWhenAPIFieldMissing(t *testing.T) {
	rec := sshPamSettingsRecord(t, json.RawMessage(`{"allowedSettings":{"connections":true}}`))

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil || settings.Connection == nil || settings.Connection.Ssh == nil {
		t.Fatal("expected ssh connection in pam_settings")
	}
	if !settings.Connection.Ssh.RotateOnTermination.IsNull() {
		t.Fatalf("expected null rotate_on_termination when API omits field, got %v", settings.Connection.Ssh.RotateOnTermination)
	}
}

func TestExtractPamSettingsFromResponse_SshRotateOnTerminationFalseWhenAPIFieldFalse(t *testing.T) {
	rec := sshPamSettingsRecord(t, json.RawMessage(`{"rotateOnTermination":false}`))

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil || settings.Connection == nil || settings.Connection.Ssh == nil {
		t.Fatal("expected ssh connection in pam_settings")
	}
	if settings.Connection.Ssh.RotateOnTermination.IsNull() || settings.Connection.Ssh.RotateOnTermination.ValueBool() {
		t.Fatalf("expected false rotate_on_termination, got %v", settings.Connection.Ssh.RotateOnTermination)
	}
}

func TestExtractPamSettingsFromResponse_SshRotateOnTerminationTrueWhenAPIFieldTrue(t *testing.T) {
	rec := sshPamSettingsRecord(t, json.RawMessage(`{"rotateOnTermination":true}`))

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil || settings.Connection == nil || settings.Connection.Ssh == nil {
		t.Fatal("expected ssh connection in pam_settings")
	}
	if settings.Connection.Ssh.RotateOnTermination.IsNull() || !settings.Connection.Ssh.RotateOnTermination.ValueBool() {
		t.Fatalf("expected true rotate_on_termination, got %v", settings.Connection.Ssh.RotateOnTermination)
	}
}

func TestExtractPamSettingsFromResponse_ExtractsKeeperDbFromConnection(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type: "pamSettings",
				Value: json.RawMessage(`[{
					"allowSupplyHost": false,
					"connection": {
						"protocol": "mariadb",
						"port": "3306",
						"recordingIncludeKeys": false,
						"dbConnectionMethod": "keeperDB",
						"rbiSettings": {
							"sessionPersistence": "none"
						}
					},
					"portForward": {
						"port": "3306",
						"reusePort": true
					}
				}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
		},
	}

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil || settings.KeeperDB == nil {
		t.Fatal("expected keeper_db to be extracted from connection")
	}
	if settings.KeeperDB.DbConnectionMethod.ValueString() != "keeperDB" {
		t.Errorf("dbConnectionMethod = %q, want keeperDB", settings.KeeperDB.DbConnectionMethod.ValueString())
	}
	if settings.KeeperDB.SessionPersistence.ValueString() != "none" {
		t.Errorf("sessionPersistence = %q, want none", settings.KeeperDB.SessionPersistence.ValueString())
	}
}

func TestExtractPamSettingsFromResponse_NoKeeperDbWhenAbsentFromConnection(t *testing.T) {
	rec := sshPamSettingsRecord(t, nil)

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil {
		t.Fatal("expected pam_settings")
	}
	if settings.KeeperDB != nil {
		t.Fatalf("expected nil keeper_db when connection has no KeeperDB fields, got %+v", settings.KeeperDB)
	}
}

func TestExtractPamSettingsFromResponse_PreservesExistingKeeperDbWhenAbsentFromResponse(t *testing.T) {
	rec := sshPamSettingsRecord(t, nil)

	existing := &commonpamrecords.CommonPamSettingsFieldResourceModel{
		KeeperDB: &commonpamrecords.DatabaseKeeperDbResourceModel{
			DbConnectionMethod: types.StringValue("keeperDB"),
			SessionPersistence: types.StringValue("user"),
		},
	}

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, existing)
	if settings == nil || settings.KeeperDB == nil {
		t.Fatal("expected existing keeper_db to be preserved")
	}
	if settings.KeeperDB.DbConnectionMethod.ValueString() != "keeperDB" || settings.KeeperDB.SessionPersistence.ValueString() != "user" {
		t.Fatalf("expected existing keeper_db preserved, got %+v", settings.KeeperDB)
	}
}

func TestApplyPamSettings_IncludesKeeperDbFieldsInConnectionJSON(t *testing.T) {
	mock := &helpers.CommandServer{}
	var capturedCommands []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		capturedCommands = append(capturedCommands, cmd)
		return "", nil
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	conn := &commonpamrecords.CommonPamSettingsConnectionResourceModel{}
	conn.Enable = types.BoolValue(true)
	conn.Protocol = types.StringValue(commonpamrecords.ConnectionProtocolMariaDb)
	conn.ConnectionPort = types.Int32Value(3306)
	conn.MariaDb = &commonpamrecords.ConnectionMariaDbOracleDatabaseModel{}

	pamSettings := &commonpamrecords.CommonPamSettingsFieldResourceModel{
		Connection: conn,
		KeeperDB: &commonpamrecords.DatabaseKeeperDbResourceModel{
			DbConnectionMethod: types.StringValue("keeperDB"),
			SessionPersistence: types.StringValue("none"),
		},
	}
	pamSettings.Configuration = types.StringValue("config-uid")

	if err := commonpamrecords.ApplyPamSettings(context.Background(), apiManager, utils.CmdRecordUpdate, "record-uid", pamSettings, nil); err != nil {
		t.Fatalf("ApplyPamSettings failed: %v", err)
	}

	var jsonCmd string
	for _, c := range capturedCommands {
		if strings.Contains(c, utils.FlagPamSettings+"=$JSON:") {
			jsonCmd = c
			break
		}
	}
	if jsonCmd == "" {
		t.Fatalf("expected a command setting pamSettings, got: %v", capturedCommands)
	}
	if !strings.Contains(jsonCmd, `"dbConnectionMethod":"keeperDB"`) {
		t.Errorf("expected dbConnectionMethod in connection JSON, got: %s", jsonCmd)
	}
	if !strings.Contains(jsonCmd, `"rbiSettings":{"sessionPersistence":"none"}`) {
		t.Errorf("expected rbiSettings.sessionPersistence in connection JSON, got: %s", jsonCmd)
	}
}

func TestExtractPamSettingsFromResponse_ExtractsClipboardDataSizeLimitFromConnection(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type: "pamSettings",
				Value: json.RawMessage(`[{
					"allowSupplyHost": false,
					"connection": {
						"allowSupplyUser": false,
						"colorScheme": "black-white",
						"database": "app_production",
						"dbConnectionMethod": "both",
						"disableCopy": false,
						"disableCsvExport": true,
						"disableCsvImport": true,
						"disablePaste": false,
						"fontSize": "12",
						"port": "3306",
						"protocol": "mysql",
						"rbiSettings": {
							"sessionPersistence": "user"
						},
						"readOnly": false,
						"recordingIncludeKeys": true,
						"userRecords": ["nDd5kukmr1Ujg3o6HAqeYA"],
						"clipboardBufferSize": 10485760
					}
				}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
		},
	}

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil || settings.Connection == nil || settings.Connection.Mysql == nil {
		t.Fatal("expected mysql connection in pam_settings")
	}
	if got := settings.Connection.Mysql.ClipboardDataSizeLimit.ValueString(); got != "10MB" {
		t.Errorf("clipboard_data_size_limit = %q, want 10MB", got)
	}
}

func TestExtractPamSettingsFromResponse_NoClipboardDataSizeLimitWhenBytesUnrecognized(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type: "pamSettings",
				Value: json.RawMessage(`[{
					"allowSupplyHost": false,
					"connection": {
						"protocol": "mysql",
						"port": "3306",
						"clipboardBufferSize": 999
					}
				}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
		},
	}

	settings := commonpamrecords.ExtractPamSettingsFromResponse(rec, nil)
	if settings == nil || settings.Connection == nil || settings.Connection.Mysql == nil {
		t.Fatal("expected mysql connection in pam_settings")
	}
	if !settings.Connection.Mysql.ClipboardDataSizeLimit.IsNull() {
		t.Errorf("expected null clipboard_data_size_limit for an unrecognized byte count, got %q", settings.Connection.Mysql.ClipboardDataSizeLimit.ValueString())
	}
}

func TestApplyPamSettings_IncludesClipboardBufferSizeInConnectionJSON(t *testing.T) {
	mock := &helpers.CommandServer{}
	var capturedCommands []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		capturedCommands = append(capturedCommands, cmd)
		return "", nil
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	conn := &commonpamrecords.CommonPamSettingsConnectionResourceModel{}
	conn.Enable = types.BoolValue(true)
	conn.Protocol = types.StringValue(commonpamrecords.ConnectionProtocolMysql)
	conn.ConnectionPort = types.Int32Value(3306)
	conn.Mysql = &commonpamrecords.ConnectionDatabaseModel{}
	conn.Mysql.ClipboardDataSizeLimit = types.StringValue("10MB")

	pamSettings := &commonpamrecords.CommonPamSettingsFieldResourceModel{Connection: conn}
	pamSettings.Configuration = types.StringValue("config-uid")

	if err := commonpamrecords.ApplyPamSettings(context.Background(), apiManager, utils.CmdRecordUpdate, "record-uid", pamSettings, nil); err != nil {
		t.Fatalf("ApplyPamSettings failed: %v", err)
	}

	var jsonCmd string
	for _, c := range capturedCommands {
		if strings.Contains(c, utils.FlagPamSettings+"=$JSON:") {
			jsonCmd = c
			break
		}
	}
	if jsonCmd == "" {
		t.Fatalf("expected a command setting pamSettings, got: %v", capturedCommands)
	}
	if !strings.Contains(jsonCmd, `"clipboardBufferSize":10485760`) {
		t.Errorf("expected clipboardBufferSize in connection JSON, got: %s", jsonCmd)
	}
}

func TestAppendOptionalOnOffFlag(t *testing.T) {
	t.Run("null omits flag", func(t *testing.T) {
		var parts []string
		commonpamrecords.AppendOptionalOnOffFlag(&parts, "--rotate-on-termination", types.BoolNull())
		if len(parts) != 0 {
			t.Fatalf("expected no flags, got %v", parts)
		}
	})

	t.Run("false appends off", func(t *testing.T) {
		var parts []string
		commonpamrecords.AppendOptionalOnOffFlag(&parts, "--rotate-on-termination", types.BoolValue(false))
		if len(parts) != 1 || parts[0] != "--rotate-on-termination=off" {
			t.Fatalf("got %v, want [--rotate-on-termination=off]", parts)
		}
	})

	t.Run("true appends on", func(t *testing.T) {
		var parts []string
		commonpamrecords.AppendOptionalOnOffFlag(&parts, "--connections-recording", types.BoolValue(true))
		if len(parts) != 1 || parts[0] != "--connections-recording=on" {
			t.Fatalf("got %v, want [--connections-recording=on]", parts)
		}
	})
}
