// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"encoding/json"
	"testing"

	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapVaultRecordGetResponse_BasicFields(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-123",
		Type:      "pamMachine",
		Title:     "My Machine",
		Notes:     "some notes",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"localhost.com","port":"1111"}]`),
			},
			{Type: "text", Label: "operatingSystem", Value: json.RawMessage(`["Linux"]`)},
			{Type: "text", Label: "instanceName", Value: json.RawMessage(`["my-instance"]`)},
			{Type: "text", Label: "instanceId", Value: json.RawMessage(`["i-12345"]`)},
			{Type: "text", Label: "providerGroup", Value: json.RawMessage(`["AWS"]`)},
			{Type: "text", Label: "providerRegion", Value: json.RawMessage(`["us-east-1"]`)},
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.Id.ValueString() != "uid-123" {
		t.Errorf("expected id uid-123, got %s", state.Id.ValueString())
	}
	if state.Title.ValueString() != "My Machine" {
		t.Errorf("expected title My Machine, got %s", state.Title.ValueString())
	}
	if state.Notes.ValueString() != "some notes" {
		t.Errorf("expected notes some notes, got %s", state.Notes.ValueString())
	}
	if state.OperatingSystem.ValueString() != "Linux" {
		t.Errorf("expected operating_system Linux, got %s", state.OperatingSystem.ValueString())
	}
	if state.InstanceName.ValueString() != "my-instance" {
		t.Errorf("expected instance_name my-instance, got %s", state.InstanceName.ValueString())
	}
	if state.InstanceId.ValueString() != "i-12345" {
		t.Errorf("expected instance_id i-12345, got %s", state.InstanceId.ValueString())
	}
	if state.ProviderGroup.ValueString() != "AWS" {
		t.Errorf("expected provider_group AWS, got %s", state.ProviderGroup.ValueString())
	}
	if state.ProviderRegion.ValueString() != "us-east-1" {
		t.Errorf("expected provider_region us-east-1, got %s", state.ProviderRegion.ValueString())
	}
}

func TestMapVaultRecordGetResponse_HostnameOrIP(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-host",
		Title:     "Host Test",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"10.0.0.1","port":"8080"}]`),
			},
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.HostnameOrIP == nil {
		t.Fatal("expected non-nil hostname_or_ip")
	}
	if state.HostnameOrIP.HostName.ValueString() != "10.0.0.1" {
		t.Errorf("expected hostname 10.0.0.1, got %s", state.HostnameOrIP.HostName.ValueString())
	}
	if state.HostnameOrIP.AdministrativePort.ValueInt32() != 8080 {
		t.Errorf("expected port 8080, got %d", state.HostnameOrIP.AdministrativePort.ValueInt32())
	}
}

func TestMapVaultRecordGetResponse_HostnameWithoutPort(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-noport",
		Title:     "No Port",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"example.com","port":""}]`),
			},
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.HostnameOrIP == nil {
		t.Fatal("expected non-nil hostname_or_ip")
	}
	if state.HostnameOrIP.HostName.ValueString() != "example.com" {
		t.Errorf("expected hostname example.com, got %s", state.HostnameOrIP.HostName.ValueString())
	}
	if !state.HostnameOrIP.AdministrativePort.IsNull() {
		t.Error("expected null port for empty string")
	}
}

func TestMapVaultRecordGetResponse_InvalidPortString(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-badport",
		Title:     "Bad Port",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"host","port":"notanumber"}]`),
			},
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.HostnameOrIP == nil {
		t.Fatal("expected non-nil hostname_or_ip")
	}
	if !state.HostnameOrIP.AdministrativePort.IsNull() {
		t.Error("expected null port for invalid number string")
	}
}

func TestMapVaultRecordGetResponse_EmptyFields(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "  ",
		Title:     "  ",
		Notes:     "  ",
		Fields:    nil,
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if !state.Title.IsNull() {
		t.Error("expected null title for whitespace-only")
	}
	if !state.Notes.IsNull() {
		t.Error("expected null notes for whitespace-only")
	}
	if state.HostnameOrIP != nil {
		t.Error("expected nil hostname_or_ip when no pamHostname field")
	}
}

func TestMapVaultRecordGetResponse_BadHostnameJSON(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-badjson",
		Title:     "Bad JSON",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`{invalid`),
			},
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.HostnameOrIP != nil {
		t.Error("expected nil hostname_or_ip for bad JSON")
	}
}

func TestMapVaultRecordGetResponse_EmptyHostnameArray(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-empty",
		Title:     "Empty Array",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[]`),
			},
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.HostnameOrIP != nil {
		t.Error("expected nil hostname_or_ip for empty array")
	}
}

func TestMapVaultRecordGetResponse_NoPamSettings(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-nosettings",
		Title:     "No Settings",
		Fields:    []utils.VaultRecordFieldResponse{},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.PamSettings != nil {
		t.Error("expected nil pam_settings when no pamSettings field")
	}
}

func TestMapVaultRecordGetResponse_WithPamSettings(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-settings",
		Title:     "With Settings",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamSettings",
				Value: json.RawMessage(`[{"allowSupplyHost":true,"portForward":{"port":"3389","reusePort":true},"connection":{"protocol":"vnc","port":"445"}}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
			Tunneling:   boolPtr(true),
		},
		DagDebug: &utils.DagDebugResponse{
			AllEdges: []utils.DagDebugEdgeResponse{
				{Type: "link", HeadUID: "config-uid-123"},
			},
		},
		AssociatedCredentials: &utils.AssociatedCredentialsResponse{
			AdminCredential: strPtr("admin-uid"),
		},
	}

	var state commonpammachine.PamMachineResourceModel
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.PamSettings == nil {
		t.Fatal("expected non-nil pam_settings")
	}
	if !state.PamSettings.AllowSupplyHost.ValueBool() {
		t.Error("expected allow_supply_host true")
	}
	if state.PamSettings.Configuration.ValueString() != "config-uid-123" {
		t.Errorf("expected configuration config-uid-123, got %s", state.PamSettings.Configuration.ValueString())
	}
	if state.PamSettings.AdministrativeCredentials.ValueString() != "admin-uid" {
		t.Errorf("expected administrative_credentials admin-uid, got %s", state.PamSettings.AdministrativeCredentials.ValueString())
	}
	if state.PamSettings.Tunnel == nil {
		t.Fatal("expected non-nil tunnel")
	}
	if !state.PamSettings.Tunnel.Enable.ValueBool() {
		t.Error("expected tunnel enable true")
	}
	if state.PamSettings.Connection == nil {
		t.Fatal("expected non-nil connection")
	}
	if !state.PamSettings.Connection.Enable.ValueBool() {
		t.Error("expected connection enable true")
	}
	if state.PamSettings.Connection.Protocol.ValueString() != "vnc" {
		t.Errorf("expected protocol vnc, got %s", state.PamSettings.Connection.Protocol.ValueString())
	}
}

func TestExtractPamHostnameFieldValue_Success(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "pamHostname", Value: json.RawMessage(`[{"hostName":"test.com","port":"443"}]`)},
	}
	model := commonpammachine.ExtractPamHostnameFieldValue(fields)
	if model == nil {
		t.Fatal("expected non-nil model")
		return
	}
	if model.HostName.ValueString() != "test.com" {
		t.Errorf("expected hostname test.com, got %s", model.HostName.ValueString())
	}
	if model.AdministrativePort.ValueInt32() != 443 {
		t.Errorf("expected port 443, got %d", model.AdministrativePort.ValueInt32())
	}
}

func TestExtractPamHostnameFieldValue_NoField(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "text", Label: "other", Value: json.RawMessage(`["value"]`)},
	}
	model := commonpammachine.ExtractPamHostnameFieldValue(fields)
	if model != nil {
		t.Error("expected nil model when no pamHostname field")
	}
}

func TestExtractPamHostnameFieldValue_BadJSON(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "pamHostname", Value: json.RawMessage(`{bad`)},
	}
	model := commonpammachine.ExtractPamHostnameFieldValue(fields)
	if model != nil {
		t.Error("expected nil model for bad JSON")
	}
}

func TestExtractPamHostnameFieldValue_EmptyArray(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "pamHostname", Value: json.RawMessage(`[]`)},
	}
	model := commonpammachine.ExtractPamHostnameFieldValue(fields)
	if model != nil {
		t.Error("expected nil model for empty array")
	}
}

func TestMapVaultRecordGetResponse_FolderFromResponse_StateMatchesUID(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-folder",
		Title:     "Folder Test",
		FolderLocation: &utils.FolderLocationResponse{
			UID:  "folder-uid-123",
			Path: "Test/My Folder",
		},
	}
	var state commonpammachine.PamMachineResourceModel
	state.FolderLocation = types.StringValue("folder-uid-123")
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.FolderLocation.ValueString() != "folder-uid-123" {
		t.Errorf("expected folder folder-uid-123 (preserved), got %s", state.FolderLocation.ValueString())
	}
}

func TestMapVaultRecordGetResponse_FolderFromResponse_StateMatchesPath(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-folder2",
		Title:     "Folder Path Test",
		FolderLocation: &utils.FolderLocationResponse{
			UID:  "folder-uid-456",
			Path: "Test/My Folder",
		},
	}
	var state commonpammachine.PamMachineResourceModel
	state.FolderLocation = types.StringValue("Test/My Folder")
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.FolderLocation.ValueString() != "Test/My Folder" {
		t.Errorf("expected folder Test/My Folder (preserved), got %s", state.FolderLocation.ValueString())
	}
}

func TestMapVaultRecordGetResponse_FolderFromResponse_StateDoesNotMatch(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-folder3",
		Title:     "Folder Mismatch",
		FolderLocation: &utils.FolderLocationResponse{
			UID:  "folder-uid-789",
			Path: "Test/Other Folder",
		},
	}
	var state commonpammachine.PamMachineResourceModel
	state.FolderLocation = types.StringValue("old-folder-uid")
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.FolderLocation.ValueString() != "Test/Other Folder" {
		t.Errorf("expected folder Test/Other Folder (from response path), got %s", state.FolderLocation.ValueString())
	}
}

func TestMapVaultRecordGetResponse_FolderNilInResponse(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-nofolder",
		Title:     "No Folder",
	}
	var state commonpammachine.PamMachineResourceModel
	state.FolderLocation = types.StringValue("some-folder")
	diags := commonpammachine.MapVaultRecordGetResponseToPamMachineModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if !state.FolderLocation.IsNull() {
		t.Error("expected null folder when API response has no folder")
	}
}

func TestExtractFolderValue_NilFolder(t *testing.T) {
	result := utils.ExtractFolderValue(nil, types.StringValue("any"))
	if !result.IsNull() {
		t.Error("expected null for nil folder response")
	}
}

func TestExtractFolderValue_EmptyUIDAndPath(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "", Path: "  "}
	result := utils.ExtractFolderValue(folder, types.StringValue("any"))
	if result.ValueString() != "  " {
		t.Errorf("expected raw path '  ' (preserved), got %q", result.ValueString())
	}
}

func TestExtractFolderValue_StateMatchesUID(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Prod/Servers"}
	state := types.StringValue("abc-123")
	result := utils.ExtractFolderValue(folder, state)
	if result.ValueString() != "abc-123" {
		t.Errorf("expected abc-123 (preserved), got %s", result.ValueString())
	}
}

func TestExtractFolderValue_StateMatchesPath(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Prod/Servers"}
	state := types.StringValue("Prod/Servers")
	result := utils.ExtractFolderValue(folder, state)
	if result.ValueString() != "Prod/Servers" {
		t.Errorf("expected Prod/Servers (preserved), got %s", result.ValueString())
	}
}

func TestExtractFolderValue_StateDoesNotMatch(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Prod/Servers"}
	state := types.StringValue("old-uid")
	result := utils.ExtractFolderValue(folder, state)
	if result.ValueString() != "Prod/Servers" {
		t.Errorf("expected Prod/Servers (path fallback), got %s", result.ValueString())
	}
}

func TestExtractFolderValue_NullState(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Prod/Servers"}
	result := utils.ExtractFolderValue(folder, types.StringNull())
	if result.ValueString() != "Prod/Servers" {
		t.Errorf("expected Prod/Servers (path, null state), got %s", result.ValueString())
	}
}

func TestExtractFolderValue_UnknownState(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Prod/Servers"}
	result := utils.ExtractFolderValue(folder, types.StringUnknown())
	if result.ValueString() != "Prod/Servers" {
		t.Errorf("expected Prod/Servers (path, unknown state), got %s", result.ValueString())
	}
}

func TestExtractFolderValue_StateMatchesPathWithSpacesAroundSlash(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Test/My Folder"}
	state := types.StringValue("Test / My Folder")
	result := utils.ExtractFolderValue(folder, state)
	if result.ValueString() != "Test / My Folder" {
		t.Errorf("expected Test / My Folder (preserved), got %s", result.ValueString())
	}
}

func TestExtractFolderValue_StateMatchesNestedPathWithSpaces(t *testing.T) {
	folder := &utils.FolderLocationResponse{UID: "abc-123", Path: "Prod/PAM/Servers"}
	state := types.StringValue("Prod / PAM / Servers")
	result := utils.ExtractFolderValue(folder, state)
	if result.ValueString() != "Prod / PAM / Servers" {
		t.Errorf("expected Prod / PAM / Servers (preserved), got %s", result.ValueString())
	}
}

func boolPtr(v bool) *bool    { return &v }
func strPtr(v string) *string { return &v }
