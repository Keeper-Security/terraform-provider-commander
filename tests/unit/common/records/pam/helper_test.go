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
