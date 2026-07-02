// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func getDSSchema(t *testing.T) dschema.Schema {
	t.Helper()
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *pamdatabase.PamDatabaseDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startDSMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (string, interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

func dsVaultRecordJSON(uid, title, notes, hostname, port string) interface{} {
	return map[string]interface{}{
		"record_uid": uid,
		"type":       "pamDatabase",
		"title":      title,
		"notes":      notes,
		"fields": []map[string]interface{}{
			{
				"type":  "pamHostname",
				"label": "pamHostname",
				"value": json.RawMessage(fmt.Sprintf(`[{"hostName":"%s","port":"%s"}]`, hostname, port)),
			},
			{
				"type":  "checkbox",
				"label": "useSSL",
				"value": json.RawMessage(`[true]`),
			},
			{
				"type":  "text",
				"label": "databaseId",
				"value": json.RawMessage(`["test_id"]`),
			},
			{
				"type":  "databaseType",
				"value": json.RawMessage(`["mongodb"]`),
			},
		},
	}
}

func TestRead_DS_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") && strings.Contains(cmd, "--format json") {
			return "ok", dsVaultRecordJSON("uid-abc", "My Database", "notes", "10.0.0.1", "5432")
		}
		return "ok", nil
	}
	server := startDSMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")
	emptyState := newDSEmptyState(t, sch)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestRead_DS_NoApiManager(t *testing.T) {
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")

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
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "")

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
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")

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
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")

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
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")
	emptyState := newDSEmptyState(t, sch)

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
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when response data is nil")
	}
}

func newDSConfigRaw(t *testing.T, sch dschema.Schema, pamDatabase string) tftypes.Value {
	t.Helper()
	tfType := sch.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", tfType)
	}

	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		if name == "pam_database" {
			vals[name] = tftypes.NewValue(tftypes.String, pamDatabase)
		} else {
			vals[name] = tftypes.NewValue(attrType, nil)
		}
	}
	return tftypes.NewValue(objType, vals)
}

func newDSEmptyState(t *testing.T, sch dschema.Schema) tftypes.Value {
	t.Helper()
	tfType := sch.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", tfType)
	}

	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(objType, vals)
}
