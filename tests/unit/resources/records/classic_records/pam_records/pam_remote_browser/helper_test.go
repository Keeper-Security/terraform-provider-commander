// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var pamRemoteBrowserSettingsAttrTypes = map[string]tftypes.Type{
	"configuration":            tftypes.String,
	"remote_browser_isolation": tftypes.Bool,
	"connections_recording":    tftypes.Bool,
	"key_events":               tftypes.Bool,
	"allow_url_navigation":     tftypes.Bool,
	"ignore_server_cert":       tftypes.Bool,
	"allowed_urls":             tftypes.Set{ElementType: tftypes.String},
	"allowed_resource_urls":    tftypes.Set{ElementType: tftypes.String},
	"auto_fill_targets":        tftypes.Set{ElementType: tftypes.String},
	"auto_fill_credentials":    tftypes.String,
	"allow_copy":               tftypes.Bool,
	"allow_paste":              tftypes.Bool,
	"disable_audio":            tftypes.Bool,
	"audio_channels":           tftypes.Number,
	"audio_bit_depth":          tftypes.Number,
	"audio_sample_rate":        tftypes.Number,
	"session_persistence":      tftypes.String,
}

var shareElementAttrTypes = map[string]tftypes.Type{
	"can_share": tftypes.Bool,
	"can_edit":  tftypes.Bool,
}

var shareMapType = tftypes.Map{ElementType: tftypes.Object{AttributeTypes: shareElementAttrTypes}}

var pamRemoteBrowserAttrTypes = map[string]tftypes.Type{
	"id":              tftypes.String,
	"title":           tftypes.String,
	"url":             tftypes.String,
	"notes":           tftypes.String,
	"folder_location": tftypes.String,
	"pam_remote_browser_settings": tftypes.Object{
		AttributeTypes: pamRemoteBrowserSettingsAttrTypes,
	},
	"share": shareMapType,
}

func pamRemoteBrowserObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: pamRemoteBrowserAttrTypes}
}

func newSettingsValues(
	configuration, autoFillCredentials interface{},
	remoteBrowserIsolation, connectionsRecording, keyEvents, allowUrlNavigation, ignoreServerCert interface{},
	allowedUrls, allowedResourceUrls, autoFillTargets interface{},
	allowCopy, allowPaste, disableAudio interface{},
	audioChannels, audioBitDepth, audioSampleRate interface{},
) map[string]tftypes.Value {
	setType := tftypes.Set{ElementType: tftypes.String}
	makeSet := func(val interface{}) tftypes.Value {
		if val == nil {
			return tftypes.NewValue(setType, nil)
		}
		elems := val.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		return tftypes.NewValue(setType, vals)
	}

	return map[string]tftypes.Value{
		"configuration":            tftypes.NewValue(tftypes.String, configuration),
		"remote_browser_isolation": tftypes.NewValue(tftypes.Bool, remoteBrowserIsolation),
		"connections_recording":    tftypes.NewValue(tftypes.Bool, connectionsRecording),
		"key_events":               tftypes.NewValue(tftypes.Bool, keyEvents),
		"allow_url_navigation":     tftypes.NewValue(tftypes.Bool, allowUrlNavigation),
		"ignore_server_cert":       tftypes.NewValue(tftypes.Bool, ignoreServerCert),
		"allowed_urls":             makeSet(allowedUrls),
		"allowed_resource_urls":    makeSet(allowedResourceUrls),
		"auto_fill_targets":        makeSet(autoFillTargets),
		"auto_fill_credentials":    tftypes.NewValue(tftypes.String, autoFillCredentials),
		"allow_copy":               tftypes.NewValue(tftypes.Bool, allowCopy),
		"allow_paste":              tftypes.NewValue(tftypes.Bool, allowPaste),
		"disable_audio":            tftypes.NewValue(tftypes.Bool, disableAudio),
		"audio_channels":           tftypes.NewValue(tftypes.Number, audioChannels),
		"audio_bit_depth":          tftypes.NewValue(tftypes.Number, audioBitDepth),
		"audio_sample_rate":        tftypes.NewValue(tftypes.Number, audioSampleRate),
		"session_persistence":      tftypes.NewValue(tftypes.String, nil),
	}
}

func newPlanStateValues(
	id, title, url, notes, folder interface{},
	settings interface{},
) map[string]tftypes.Value {
	settingsObjType := tftypes.Object{AttributeTypes: pamRemoteBrowserSettingsAttrTypes}
	var settingsVal tftypes.Value
	if settings == nil {
		settingsVal = tftypes.NewValue(settingsObjType, nil)
	} else {
		settingsVal = tftypes.NewValue(settingsObjType, settings)
	}
	return map[string]tftypes.Value{
		"id":                          tftypes.NewValue(tftypes.String, id),
		"title":                       tftypes.NewValue(tftypes.String, title),
		"url":                         tftypes.NewValue(tftypes.String, url),
		"notes":                       tftypes.NewValue(tftypes.String, notes),
		"folder_location":             tftypes.NewValue(tftypes.String, folder),
		"pam_remote_browser_settings": settingsVal,
		"share":                       tftypes.NewValue(shareMapType, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := pamRemoteBrowserObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *pamremotebrowser.PamRemoteBrowserResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

func TestBuildPamRbiEditCommand_NilSettings(t *testing.T) {
	cmd := commonpamremotebrowser.BuildPamRbiEditCommand("rec-uid-123", nil)
	if cmd != "pam rbi edit --record 'rec-uid-123'" {
		t.Errorf("unexpected command: %s", cmd)
	}
}

func TestBuildPamRbiEditCommand_WithSettings(t *testing.T) {
	settings := &commonpamremotebrowser.PamRemoteBrowserSettingsModel{
		Configuration:          types.StringValue("config-uid"),
		RemoteBrowserIsolation: types.BoolValue(true),
		ConnectionsRecording:   types.BoolValue(false),
		KeyEvents:              types.BoolNull(),
		AllowUrlNavigation:     types.BoolValue(true),
		IgnoreServerCert:       types.BoolUnknown(),
		AllowedUrls:            types.SetNull(types.StringType),
		AllowedResourceUrls:    types.SetNull(types.StringType),
		AutoFillTargets:        types.SetNull(types.StringType),
		AutoFillCredentials:    types.StringNull(),
		AllowCopy:              types.BoolValue(true),
		AllowPaste:             types.BoolValue(false),
		DisableAudio:           types.BoolValue(false),
		AudioChannels:          types.Int32Value(2),
		AudioBitDepth:          types.Int64Value(16),
		AudioSampleRate:        types.Int64Value(44100),
	}
	cmd := commonpamremotebrowser.BuildPamRbiEditCommand("uid-1", settings)
	if cmd == "" {
		t.Fatal("expected non-empty command")
	}
	// --record must be present
	if !contains(cmd, "--record 'uid-1'") {
		t.Errorf("command should contain --record: %s", cmd)
	}
	// on/off booleans
	if !contains(cmd, "--remote-browser-isolation on") {
		t.Errorf("expected --remote-browser-isolation on: %s", cmd)
	}
	if !contains(cmd, "--connections-recording off") {
		t.Errorf("expected --connections-recording off: %s", cmd)
	}
	// null bool => off
	if !contains(cmd, "--key-events off") {
		t.Errorf("expected --key-events off (null): %s", cmd)
	}
	// unknown bool => omitted
	if contains(cmd, "--ignore-server-cert") {
		t.Errorf("unknown --ignore-server-cert should be omitted: %s", cmd)
	}
	if !contains(cmd, "--audio-channels 2") {
		t.Errorf("expected --audio-channels 2: %s", cmd)
	}
}

func TestBuildPamRbiEditCommand_UnknownSet_Omitted(t *testing.T) {
	settings := &commonpamremotebrowser.PamRemoteBrowserSettingsModel{
		Configuration:          types.StringValue("cfg"),
		RemoteBrowserIsolation: types.BoolNull(),
		ConnectionsRecording:   types.BoolNull(),
		KeyEvents:              types.BoolNull(),
		AllowUrlNavigation:     types.BoolNull(),
		IgnoreServerCert:       types.BoolNull(),
		AllowedUrls:            types.SetUnknown(types.StringType),
		AllowedResourceUrls:    types.SetUnknown(types.StringType),
		AutoFillTargets:        types.SetUnknown(types.StringType),
		AutoFillCredentials:    types.StringUnknown(),
		AllowCopy:              types.BoolNull(),
		AllowPaste:             types.BoolNull(),
		DisableAudio:           types.BoolNull(),
		AudioChannels:          types.Int32Null(),
		AudioBitDepth:          types.Int64Null(),
		AudioSampleRate:        types.Int64Null(),
	}
	cmd := commonpamremotebrowser.BuildPamRbiEditCommand("uid-2", settings)
	if contains(cmd, "--allowed-urls") {
		t.Errorf("unknown --allowed-urls should be omitted: %s", cmd)
	}
	if contains(cmd, "--autofill-credentials") {
		t.Errorf("unknown --autofill-credentials should be omitted: %s", cmd)
	}
}

func TestAppendPamRbiEditSettingsFlags_NilSettings(t *testing.T) {
	var parts []string
	commonpamremotebrowser.AppendPamRbiEditSettingsFlags(&parts, nil)
	if len(parts) != 0 {
		t.Errorf("expected empty parts for nil settings, got %v", parts)
	}
}

func TestAppendPamRbiEditSettingsFlags_WithNonNullSets(t *testing.T) {
	ctx := context.Background()
	urlSet, _ := types.SetValueFrom(ctx, types.StringType, []string{"https://example.com"})
	settings := &commonpamremotebrowser.PamRemoteBrowserSettingsModel{
		Configuration:          types.StringValue("cfg-uid"),
		RemoteBrowserIsolation: types.BoolValue(true),
		ConnectionsRecording:   types.BoolValue(true),
		KeyEvents:              types.BoolValue(true),
		AllowUrlNavigation:     types.BoolValue(false),
		IgnoreServerCert:       types.BoolValue(false),
		AllowedUrls:            urlSet,
		AllowedResourceUrls:    types.SetNull(types.StringType),
		AutoFillTargets:        types.SetNull(types.StringType),
		AutoFillCredentials:    types.StringValue("cred-uid"),
		AllowCopy:              types.BoolValue(true),
		AllowPaste:             types.BoolValue(true),
		DisableAudio:           types.BoolValue(false),
		AudioChannels:          types.Int32Value(1),
		AudioBitDepth:          types.Int64Value(8),
		AudioSampleRate:        types.Int64Value(48000),
	}
	var parts []string
	commonpamremotebrowser.AppendPamRbiEditSettingsFlags(&parts, settings)
	if len(parts) == 0 {
		t.Fatal("expected parts to have flags")
	}
	joined := joinParts(parts)
	if !contains(joined, "--configuration 'cfg-uid'") {
		t.Errorf("expected --configuration flag: %s", joined)
	}
	if !contains(joined, "--allowed-urls 'https://example.com'") {
		t.Errorf("expected --allowed-urls flag: %s", joined)
	}
	if !contains(joined, "--audio-channels 1") {
		t.Errorf("expected --audio-channels 1: %s", joined)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || containsSub(s, sub))
}

func containsSub(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func joinParts(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}
	return result
}
