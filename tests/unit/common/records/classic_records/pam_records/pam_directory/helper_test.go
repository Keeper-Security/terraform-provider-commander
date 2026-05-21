// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"encoding/json"
	"testing"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapVaultRecordGetResponse_BasicFields(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-123",
		Type:      "pamDirectory",
		Title:     "My Directory",
		Notes:     "some notes",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"ldap.example.com","port":"636"}]`),
			},
			{Type: "checkbox", Label: "useSSL", Value: json.RawMessage(`[true]`)},
			{Type: "text", Label: "domainName", Value: json.RawMessage(`["example.com"]`)},
			{Type: "multiline", Label: "alternativeIPs", Value: json.RawMessage(`["10.0.0.1\n10.0.0.2"]`)},
			{Type: "text", Label: "directoryId", Value: json.RawMessage(`["dir-id-123"]`)},
			{Type: "directoryType", Value: json.RawMessage(`["active_directory"]`)},
			{Type: "text", Label: "userMatch", Value: json.RawMessage(`["OU=Users,DC=example,DC=com"]`)},
			{Type: "text", Label: "providerGroup", Value: json.RawMessage(`["Azure"]`)},
			{Type: "text", Label: "providerRegion", Value: json.RawMessage(`["us-west-2"]`)},
		},
	}

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.Id.ValueString() != "uid-123" {
		t.Errorf("expected id uid-123, got %s", state.Id.ValueString())
	}
	if state.Title.ValueString() != "My Directory" {
		t.Errorf("expected title My Directory, got %s", state.Title.ValueString())
	}
	if state.Notes.ValueString() != "some notes" {
		t.Errorf("expected notes some notes, got %s", state.Notes.ValueString())
	}
	if !state.UseSSL.ValueBool() {
		t.Error("expected use_ssl true")
	}
	if state.DomainName.ValueString() != "example.com" {
		t.Errorf("expected domain_name example.com, got %s", state.DomainName.ValueString())
	}
	if state.DirectoryId.ValueString() != "dir-id-123" {
		t.Errorf("expected directory_id dir-id-123, got %s", state.DirectoryId.ValueString())
	}
	if state.DirectoryType.ValueString() != "active_directory" {
		t.Errorf("expected directory_type active_directory, got %s", state.DirectoryType.ValueString())
	}
	if state.UserMatch.ValueString() != "OU=Users,DC=example,DC=com" {
		t.Errorf("expected user_match OU=Users,DC=example,DC=com, got %s", state.UserMatch.ValueString())
	}
	if state.ProviderGroup.ValueString() != "Azure" {
		t.Errorf("expected provider_group Azure, got %s", state.ProviderGroup.ValueString())
	}
	if state.ProviderRegion.ValueString() != "us-west-2" {
		t.Errorf("expected provider_region us-west-2, got %s", state.ProviderRegion.ValueString())
	}
}

func TestMapVaultRecordGetResponse_HostnameOrIP(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-host",
		Title:     "Host Test",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"10.0.0.1","port":"389"}]`),
			},
		},
	}

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.HostnameOrIP == nil {
		t.Fatal("expected non-nil hostname_or_ip")
	}
	if state.HostnameOrIP.HostName.ValueString() != "10.0.0.1" {
		t.Errorf("expected hostname 10.0.0.1, got %s", state.HostnameOrIP.HostName.ValueString())
	}
	if state.HostnameOrIP.AdministrativePort.ValueInt32() != 389 {
		t.Errorf("expected port 389, got %d", state.HostnameOrIP.AdministrativePort.ValueInt32())
	}
}

func TestMapVaultRecordGetResponse_HostnameWithoutPort(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-noport",
		Title:     "No Port",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamHostname",
				Value: json.RawMessage(`[{"hostName":"ldap.test.com","port":""}]`),
			},
		},
	}

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.HostnameOrIP == nil {
		t.Fatal("expected non-nil hostname_or_ip")
	}
	if state.HostnameOrIP.HostName.ValueString() != "ldap.test.com" {
		t.Errorf("expected hostname ldap.test.com, got %s", state.HostnameOrIP.HostName.ValueString())
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

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
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

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
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

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
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

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
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

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
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
				Value: json.RawMessage(`[{"allowSupplyHost":true,"portForward":{"port":"3389","reusePort":true},"connection":{"protocol":"ssh","port":"22"}}]`),
			},
		},
		PamSettingsEnabled: &utils.PamSettingsEnabledResponse{
			Connections: boolPtr(true),
			Tunneling:   boolPtr(true),
		},
		DagDebug: &utils.DagDebugResponse{
			AllEdges: []utils.DagDebugEdgeResponse{
				{Type: "link", HeadUID: "config-uid-456"},
			},
		},
		AssociatedCredentials: &utils.AssociatedCredentialsResponse{
			AdminCredential: strPtr("admin-uid"),
		},
	}

	var state commonpamdirectory.PamDirectoryResourceModel
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.PamSettings == nil {
		t.Fatal("expected non-nil pam_settings")
	}
	if !state.PamSettings.AllowSupplyHost.ValueBool() {
		t.Error("expected allow_supply_host true")
	}
	if state.PamSettings.Configuration.ValueString() != "config-uid-456" {
		t.Errorf("expected configuration config-uid-456, got %s", state.PamSettings.Configuration.ValueString())
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
	if state.PamSettings.Connection.Protocol.ValueString() != "ssh" {
		t.Errorf("expected protocol ssh, got %s", state.PamSettings.Connection.Protocol.ValueString())
	}
}

func TestExtractPamHostnameFieldValue_Success(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "pamHostname", Value: json.RawMessage(`[{"hostName":"test.com","port":"443"}]`)},
	}
	model := commonpamdirectory.ExtractPamHostnameFieldValue(fields)
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
	model := commonpamdirectory.ExtractPamHostnameFieldValue(fields)
	if model != nil {
		t.Error("expected nil model when no pamHostname field")
	}
}

func TestExtractPamHostnameFieldValue_BadJSON(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "pamHostname", Value: json.RawMessage(`{bad`)},
	}
	model := commonpamdirectory.ExtractPamHostnameFieldValue(fields)
	if model != nil {
		t.Error("expected nil model for bad JSON")
	}
}

func TestExtractPamHostnameFieldValue_EmptyArray(t *testing.T) {
	fields := []utils.VaultRecordFieldResponse{
		{Type: "pamHostname", Value: json.RawMessage(`[]`)},
	}
	model := commonpamdirectory.ExtractPamHostnameFieldValue(fields)
	if model != nil {
		t.Error("expected nil model for empty array")
	}
}

func TestExtractCheckboxFieldValue_True(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-ssl-true",
		Title:     "SSL True",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "checkbox", Label: "useSSL", Value: json.RawMessage(`[true]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.UseSSL.ValueBool() {
		t.Error("expected use_ssl true")
	}
}

func TestExtractCheckboxFieldValue_False(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-ssl-false",
		Title:     "SSL False",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "checkbox", Label: "useSSL", Value: json.RawMessage(`[false]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if state.UseSSL.ValueBool() {
		t.Error("expected use_ssl false")
	}
}

func TestExtractCheckboxFieldValue_BadJSON(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-ssl-bad",
		Title:     "SSL Bad",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "checkbox", Label: "useSSL", Value: json.RawMessage(`{invalid`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.UseSSL.IsNull() {
		t.Error("expected null use_ssl for bad JSON")
	}
}

func TestExtractCheckboxFieldValue_EmptyArray(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-ssl-empty",
		Title:     "SSL Empty",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "checkbox", Label: "useSSL", Value: json.RawMessage(`[]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.UseSSL.IsNull() {
		t.Error("expected null use_ssl for empty array")
	}
}

func TestExtractCheckboxFieldValue_NoMatch(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-no-ssl",
		Title:     "No SSL",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "checkbox", Label: "otherCheckbox", Value: json.RawMessage(`[true]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.UseSSL.IsNull() {
		t.Error("expected null use_ssl when label does not match")
	}
}

func TestExtractDirectoryTypeFieldValue_ActiveDirectory(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dt-ad",
		Title:     "DT AD",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "directoryType", Value: json.RawMessage(`["active_directory"]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if state.DirectoryType.ValueString() != "active_directory" {
		t.Errorf("expected active_directory, got %s", state.DirectoryType.ValueString())
	}
}

func TestExtractDirectoryTypeFieldValue_OpenLDAP(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dt-ldap",
		Title:     "DT LDAP",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "directoryType", Value: json.RawMessage(`["openldap"]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if state.DirectoryType.ValueString() != "openldap" {
		t.Errorf("expected openldap, got %s", state.DirectoryType.ValueString())
	}
}

func TestExtractDirectoryTypeFieldValue_BadJSON(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dt-bad",
		Title:     "DT Bad",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "directoryType", Value: json.RawMessage(`{invalid`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.DirectoryType.IsNull() {
		t.Error("expected null directory_type for bad JSON")
	}
}

func TestExtractDirectoryTypeFieldValue_EmptyArray(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dt-empty",
		Title:     "DT Empty",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "directoryType", Value: json.RawMessage(`[]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.DirectoryType.IsNull() {
		t.Error("expected null directory_type for empty array")
	}
}

func TestExtractDirectoryTypeFieldValue_EmptyString(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dt-es",
		Title:     "DT EmptyStr",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "directoryType", Value: json.RawMessage(`[""]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.DirectoryType.IsNull() {
		t.Error("expected null directory_type for empty string value")
	}
}

func TestExtractDirectoryTypeFieldValue_NoField(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dt-nf",
		Title:     "DT NoField",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "text", Label: "domainName", Value: json.RawMessage(`["example.com"]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.DirectoryType.IsNull() {
		t.Error("expected null directory_type when no directoryType field")
	}
}

func TestExtractMultilineAsSet_WithValues(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-alt-ips",
		Title:     "Alt IPs",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "multiline", Label: "alternativeIPs", Value: json.RawMessage(`["10.0.0.1\n10.0.0.2\n10.0.0.3"]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if state.AlternativeIPs.IsNull() {
		t.Fatal("expected non-null alternative_ips")
	}
	elems := state.AlternativeIPs.Elements()
	if len(elems) != 3 {
		t.Errorf("expected 3 elements, got %d", len(elems))
	}
}

func TestExtractMultilineAsSet_EmptyValue(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-alt-empty",
		Title:     "Alt Empty",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "multiline", Label: "alternativeIPs", Value: json.RawMessage(`[""]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.AlternativeIPs.IsNull() {
		t.Error("expected null alternative_ips for empty string")
	}
}

func TestExtractMultilineAsSet_BadJSON(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-alt-bad",
		Title:     "Alt Bad",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "multiline", Label: "alternativeIPs", Value: json.RawMessage(`{invalid`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.AlternativeIPs.IsNull() {
		t.Error("expected null alternative_ips for bad JSON")
	}
}

func TestExtractMultilineAsSet_EmptyArrayValue(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-alt-ea",
		Title:     "Alt EA",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "multiline", Label: "alternativeIPs", Value: json.RawMessage(`[]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.AlternativeIPs.IsNull() {
		t.Error("expected null alternative_ips for empty array")
	}
}

func TestExtractMultilineAsSet_NoMatchingLabel(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-alt-nml",
		Title:     "Alt NoMatch",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "multiline", Label: "otherField", Value: json.RawMessage(`["10.0.0.1"]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.AlternativeIPs.IsNull() {
		t.Error("expected null alternative_ips when label does not match")
	}
}

func TestExtractMultilineAsSet_WhitespaceOnlyLines(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-alt-ws",
		Title:     "Alt WS",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "multiline", Label: "alternativeIPs", Value: json.RawMessage(`["  \n  \n  "]`)},
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if !state.AlternativeIPs.IsNull() {
		t.Error("expected null alternative_ips when all lines are whitespace")
	}
}

func TestMapVaultRecordGetResponse_FolderFromResponse_StateMatchesUID(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-folder",
		Title:     "Folder Test",
		Folder: &utils.RecordFolderResponse{
			UID:  "folder-uid-123",
			Path: "Test/My Folder",
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	state.Folder = types.StringValue("folder-uid-123")
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.Folder.ValueString() != "folder-uid-123" {
		t.Errorf("expected folder folder-uid-123 (preserved), got %s", state.Folder.ValueString())
	}
}

func TestMapVaultRecordGetResponse_FolderFromResponse_StateMatchesPath(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-folder2",
		Title:     "Folder Path Test",
		Folder: &utils.RecordFolderResponse{
			UID:  "folder-uid-456",
			Path: "Test/My Folder",
		},
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	state.Folder = types.StringValue("Test/My Folder")
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.Folder.ValueString() != "Test/My Folder" {
		t.Errorf("expected folder Test/My Folder (preserved), got %s", state.Folder.ValueString())
	}
}

func TestMapVaultRecordGetResponse_FolderNilInResponse(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-nofolder",
		Title:     "No Folder",
	}
	var state commonpamdirectory.PamDirectoryResourceModel
	state.Folder = types.StringValue("some-folder")
	diags := commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if !state.Folder.IsNull() {
		t.Error("expected null folder when API response has no folder")
	}
}

func TestExtractFolderValue_NilFolder(t *testing.T) {
	result := commonpamrecords.ExtractFolderValue(nil, types.StringValue("any"))
	if !result.IsNull() {
		t.Error("expected null for nil folder response")
	}
}

func boolPtr(v bool) *bool    { return &v }
func strPtr(v string) *string { return &v }
