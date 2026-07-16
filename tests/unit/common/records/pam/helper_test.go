// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords_test

import (
	"encoding/json"
	"testing"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
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
