// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var dsSettingsAttrTypes = map[string]tftypes.Type{
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
}

var dsAttrTypes = map[string]tftypes.Type{
	"remote_browser": tftypes.String,
	"id":             tftypes.String,
	"title":          tftypes.String,
	"url":            tftypes.String,
	"notes":          tftypes.String,
	"folder":         tftypes.String,
	"pam_remote_browser_settings": tftypes.Object{
		AttributeTypes: dsSettingsAttrTypes,
	},
}

func dsObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: dsAttrTypes}
}

func newDSConfigValues(remoteBrowser interface{}) map[string]tftypes.Value {
	settingsObjType := tftypes.Object{AttributeTypes: dsSettingsAttrTypes}
	return map[string]tftypes.Value{
		"remote_browser":              tftypes.NewValue(tftypes.String, remoteBrowser),
		"id":                          tftypes.NewValue(tftypes.String, nil),
		"title":                       tftypes.NewValue(tftypes.String, nil),
		"url":                         tftypes.NewValue(tftypes.String, nil),
		"notes":                       tftypes.NewValue(tftypes.String, nil),
		"folder":                      tftypes.NewValue(tftypes.String, nil),
		"pam_remote_browser_settings": tftypes.NewValue(settingsObjType, nil),
	}
}

func getDSSchema(t *testing.T) (dschema.Schema, tftypes.Object) {
	t.Helper()
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema, dsObjectType()
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *pamremotebrowser.PamRemoteBrowserDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startDSMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (string, interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

func dsVaultRecordJSON(uid, title, url, notes string) interface{} {
	return map[string]interface{}{
		"record_uid": uid,
		"type":       "pamRemoteBrowser",
		"title":      title,
		"notes":      notes,
		"fields": []map[string]interface{}{
			{
				"type":  "rbiUrl",
				"label": "URL",
				"value": json.RawMessage(`["` + url + `"]`),
			},
		},
	}
}

func TestRead_DS_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") && strings.Contains(cmd, "--format json") {
			return "ok", dsVaultRecordJSON("uid-abc", "My RBI", "https://example.com", "notes")
		}
		return "ok", nil
	}
	server := startDSMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues("uid-abc"))
	emptyState := tftypes.NewValue(objType, newDSConfigValues(nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestRead_DS_NoApiManager(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues("uid-abc"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestRead_DS_EmptyRecordUID(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startDSMockServer(mock, nil)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues(""))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record_uid is empty")
	}
}

func TestRead_DS_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues("uid-abc"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestRead_DS_ApiError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "get", nil)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues("uid-abc"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when get command fails")
	}
}

func TestRead_DS_WrongRecordType(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") {
			return "ok", map[string]interface{}{
				"record_uid": "uid-abc",
				"type":       "pamMachine",
				"title":      "Wrong",
				"notes":      "",
				"fields":     []interface{}{},
			}
		}
		return "ok", nil
	}
	server := startDSMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues("uid-abc"))
	emptyState := tftypes.NewValue(objType, newDSConfigValues(nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record type does not match")
	}
}

func TestRead_DS_NilResponse(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		return "ok", nil
	}
	server := startDSMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	configRaw := tftypes.NewValue(objType, newDSConfigValues("uid-abc"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when response data is nil")
	}
}
