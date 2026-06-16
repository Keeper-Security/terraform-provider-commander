// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"encoding/json"
	"testing"

	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

func TestMapVaultRecordGetResponse_BasicFields(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-123",
		Type:      "pamRemoteBrowser",
		Title:     "My RBI",
		Notes:     "some notes",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "rbiUrl",
				Value: json.RawMessage(`["https://example.com"]`),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.Id.ValueString() != "uid-123" {
		t.Errorf("expected id uid-123, got %s", state.Id.ValueString())
	}
	if state.Title.ValueString() != "My RBI" {
		t.Errorf("expected title My RBI, got %s", state.Title.ValueString())
	}
	if state.Url.ValueString() != "https://example.com" {
		t.Errorf("expected url https://example.com, got %s", state.Url.ValueString())
	}
	if state.Notes.ValueString() != "some notes" {
		t.Errorf("expected notes some notes, got %s", state.Notes.ValueString())
	}
	if state.PamRemoteBrowserSettings != nil {
		t.Error("expected nil settings when no pamRemoteBrowserSettings field")
	}
}

func TestMapVaultRecordGetResponse_EmptyFields(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "  ",
		Title:     "  ",
		Notes:     "  ",
		Fields:    nil,
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if !state.Title.IsNull() {
		t.Error("expected null title for whitespace-only")
	}
	if !state.Notes.IsNull() {
		t.Error("expected null notes for whitespace-only")
	}
	if !state.FolderLocation.IsNull() {
		t.Error("expected null folder for whitespace-only")
	}
	if !state.Url.IsNull() {
		t.Error("expected null url when no rbiUrl field")
	}
}

func TestMapVaultRecordGetResponse_WithSettings(t *testing.T) {
	settingsJSON := `[{"connection":{"configurationUid":"cfg-uid","httpCredentialsUid":"cred-uid","autofillConfiguration":"target1,target2","recordingIncludeKeys":true,"recordingScreens":false,"remoteBrowserIsolation":true,"allowUrlManipulation":true,"ignoreInitialSslCert":false,"allowedUrlPatterns":"https://allowed.com\nhttps://other.com","allowedResourceUrlPatterns":"https://res.com","disableCopy":true,"disablePaste":false,"disableAudio":false,"audioChannels":2,"audioBps":16,"audioSampleRate":44100}}]`

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-456",
		Type:      "pamRemoteBrowser",
		Title:     "Settings Test",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "rbiUrl",
				Value: json.RawMessage(`["https://app.example.com"]`),
			},
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	s := state.PamRemoteBrowserSettings

	if s.KeyEvents.ValueBool() != true {
		t.Error("expected key_events true")
	}
	if s.AllowUrlNavigation.ValueBool() != true {
		t.Error("expected allow_url_navigation true")
	}
	if s.IgnoreServerCert.ValueBool() != false {
		t.Error("expected ignore_server_cert false")
	}
	// AllowCopy is inverse of DisableCopy; DisableCopy=true => AllowCopy=false
	if s.AllowCopy.ValueBool() != false {
		t.Errorf("expected allow_copy false (inverse of disable_copy=true), got %v", s.AllowCopy.ValueBool())
	}
	if s.AllowPaste.ValueBool() != true {
		t.Errorf("expected allow_paste true (inverse of disable_paste=false), got %v", s.AllowPaste.ValueBool())
	}
	if s.AutoFillCredentials.ValueString() != "cred-uid" {
		t.Errorf("expected auto_fill_credentials cred-uid, got %s", s.AutoFillCredentials.ValueString())
	}
	if s.AudioChannels.ValueInt32() != 2 {
		t.Errorf("expected audio_channels 2, got %d", s.AudioChannels.ValueInt32())
	}
	if s.AudioBitDepth.ValueInt64() != 16 {
		t.Errorf("expected audio_bit_depth 16, got %d", s.AudioBitDepth.ValueInt64())
	}
	if s.AudioSampleRate.ValueInt64() != 44100 {
		t.Errorf("expected audio_sample_rate 44100, got %d", s.AudioSampleRate.ValueInt64())
	}

	if s.AllowedUrls.IsNull() {
		t.Error("expected non-null allowed_urls")
	}
	if s.AllowedResourceUrls.IsNull() {
		t.Error("expected non-null allowed_resource_urls")
	}
	if s.AutoFillTargets.IsNull() {
		t.Error("expected non-null auto_fill_targets")
	}

	if !s.Configuration.IsNull() {
		t.Error("expected null configuration when pam_configuration_uid not in response")
	}
	if !s.RemoteBrowserIsolation.IsNull() {
		t.Error("expected null remote_browser_isolation when configuration_allowed_settings not in response")
	}
	if !s.ConnectionsRecording.IsNull() {
		t.Error("expected null connections_recording when configuration_allowed_settings not in response")
	}
}

func TestMapVaultRecordGetResponse_WithConfigurationAndAllowedSettings(t *testing.T) {
	settingsJSON := `[{"connection":{"recordingIncludeKeys":false,"allowUrlManipulation":false,"ignoreInitialSslCert":false,"disableCopy":false,"disablePaste":false,"disableAudio":false}}]`

	rec := &utils.VaultRecordGetResponse{
		RecordUID:           "uid-cfg",
		Type:                "pamRemoteBrowser",
		Title:               "Config Test",
		PamConfigurationUID: "08OV7gNVRky9BtStSsBGEw",
		ConfigurationAllowedSettings: &utils.ConfigurationAllowedSettingsResponse{
			ConnectionsRecording:   true,
			RemoteBrowserIsolation: true,
		},
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "rbiUrl",
				Value: json.RawMessage(`["https://example.com"]`),
			},
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	s := state.PamRemoteBrowserSettings

	if s.Configuration.ValueString() != "08OV7gNVRky9BtStSsBGEw" {
		t.Errorf("expected configuration 08OV7gNVRky9BtStSsBGEw, got %s", s.Configuration.ValueString())
	}
	if s.RemoteBrowserIsolation.ValueBool() != true {
		t.Error("expected remote_browser_isolation true")
	}
	if s.ConnectionsRecording.ValueBool() != true {
		t.Error("expected connections_recording true")
	}
}

func TestMapVaultRecordGetResponse_ConfigAllowedSettingsFalse(t *testing.T) {
	settingsJSON := `[{"connection":{"disableCopy":false,"disablePaste":false,"disableAudio":false}}]`

	rec := &utils.VaultRecordGetResponse{
		RecordUID:           "uid-cfg-false",
		Type:                "pamRemoteBrowser",
		Title:               "Config False Test",
		PamConfigurationUID: "some-uid",
		ConfigurationAllowedSettings: &utils.ConfigurationAllowedSettingsResponse{
			ConnectionsRecording:   false,
			RemoteBrowserIsolation: false,
		},
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	s := state.PamRemoteBrowserSettings
	if s == nil {
		t.Fatal("expected non-nil settings")
		return
	}
	if s.RemoteBrowserIsolation.ValueBool() != false {
		t.Error("expected remote_browser_isolation false")
	}
	if s.ConnectionsRecording.ValueBool() != false {
		t.Error("expected connections_recording false")
	}
}

func TestMapVaultRecordGetResponse_EmptyPamConfigurationUID(t *testing.T) {
	settingsJSON := `[{"connection":{"disableCopy":false,"disablePaste":false,"disableAudio":false}}]`

	rec := &utils.VaultRecordGetResponse{
		RecordUID:           "uid-empty-cfg",
		Type:                "pamRemoteBrowser",
		Title:               "Empty Config UID",
		PamConfigurationUID: "",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}

	s := state.PamRemoteBrowserSettings
	if s == nil {
		t.Fatal("expected non-nil settings")
		return
	}
	if !s.Configuration.IsNull() {
		t.Error("expected null configuration for empty pam_configuration_uid")
	}
}

func TestMapVaultRecordGetResponse_WithFolder(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-789",
		Title:     "Folder Test",
		FolderLocation: &utils.FolderLocationResponse{
			UID:  "my-folder",
			Path: "Test/My Folder",
		},
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "rbiUrl",
				Value: json.RawMessage(`["https://example.com"]`),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.FolderLocation.ValueString() != "Test/My Folder" {
		t.Errorf("expected folder Test/My Folder, got %s", state.FolderLocation.ValueString())
	}
}

func TestMapVaultRecordGetResponse_BadRbiUrlJSON(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-bad",
		Title:     "Bad URL",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "rbiUrl",
				Value: json.RawMessage(`{invalid`),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if !diags.HasError() && len(diags.Warnings()) == 0 {
		t.Error("expected warning for bad rbiUrl JSON")
	}
}

func TestMapVaultRecordGetResponse_BadSettingsJSON(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-bad2",
		Title:     "Bad Settings",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(`{invalid`),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if !diags.HasError() && len(diags.Warnings()) == 0 {
		t.Error("expected warning for bad settings JSON")
	}
}

func TestMapVaultRecordGetResponse_EmptyRbiUrlArray(t *testing.T) {
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-empty-url",
		Title:     "Empty URL",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "rbiUrl",
				Value: json.RawMessage(`[]`),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if !state.Url.IsNull() {
		t.Error("expected null url for empty rbiUrl array")
	}
}

func TestMapVaultRecordGetResponse_ZeroAudioValues(t *testing.T) {
	settingsJSON := `[{"connection":{"audioChannels":0,"audioBps":0,"audioSampleRate":0}}]`
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-zero-audio",
		Title:     "Zero Audio",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	if !state.PamRemoteBrowserSettings.AudioChannels.IsNull() {
		t.Error("expected null audio_channels for 0")
	}
	if !state.PamRemoteBrowserSettings.AudioBitDepth.IsNull() {
		t.Error("expected null audio_bit_depth for 0")
	}
	if !state.PamRemoteBrowserSettings.AudioSampleRate.IsNull() {
		t.Error("expected null audio_sample_rate for 0")
	}
}

func TestMapVaultRecordGetResponse_EmptyHttpCredentialsUID(t *testing.T) {
	settingsJSON := `[{"connection":{"httpCredentialsUid":""}}]`
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-no-creds",
		Title:     "No Creds",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	if !state.PamRemoteBrowserSettings.AutoFillCredentials.IsNull() {
		t.Error("expected null auto_fill_credentials for empty httpCredentialsUid")
	}
}

func TestMapVaultRecordGetResponse_EmptyAutofillConfiguration(t *testing.T) {
	settingsJSON := `[{"connection":{"autofillConfiguration":""}}]`
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-no-autofill",
		Title:     "No Autofill",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	if !state.PamRemoteBrowserSettings.AutoFillTargets.IsNull() {
		t.Error("expected null auto_fill_targets for empty autofillConfiguration")
	}
}

func TestMapVaultRecordGetResponse_DuplicateAutofillTargets(t *testing.T) {
	settingsJSON := `[{"connection":{"autofillConfiguration":"target1,target1,target2"}}]`
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-dupe-targets",
		Title:     "Dupe Targets",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	targets := state.PamRemoteBrowserSettings.AutoFillTargets
	if targets.IsNull() {
		t.Fatal("expected non-null auto_fill_targets")
	}
	if len(targets.Elements()) != 2 {
		t.Errorf("expected 2 unique targets after dedup, got %d", len(targets.Elements()))
	}
}

func TestMapVaultRecordGetResponse_MultiLineAllowedUrlPatterns(t *testing.T) {
	settingsJSON := `[{"connection":{"allowedUrlPatterns":"https://a.com\r\nhttps://b.com\n\nhttps://c.com"}}]`
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-multiline",
		Title:     "Multiline",
		Fields: []utils.VaultRecordFieldResponse{
			{
				Type:  "pamRemoteBrowserSettings",
				Value: json.RawMessage(settingsJSON),
			},
		},
	}

	var state commonpamremotebrowser.PamRemoteBrowserResourceModel
	diags := commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(context.Background(), rec, &state)
	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags)
	}
	if state.PamRemoteBrowserSettings == nil {
		t.Fatal("expected non-nil settings")
	}
	urls := state.PamRemoteBrowserSettings.AllowedUrls
	if urls.IsNull() {
		t.Fatal("expected non-null allowed_urls")
	}
	if len(urls.Elements()) != 3 {
		t.Errorf("expected 3 allowed urls, got %d", len(urls.Elements()))
	}
}
