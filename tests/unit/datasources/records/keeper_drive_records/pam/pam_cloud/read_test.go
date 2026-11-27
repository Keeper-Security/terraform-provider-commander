// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamcloud_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_cloud"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func getDSSchema(t *testing.T) dschema.Schema {
	t.Helper()
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *newpamcloud.PamCloudDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
	}
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

// dsVaultRecordJSON mirrors the `nsf-get <uid> --format json` payload
// Commander returns for a pamCloudResource record. Unlike `get`, NSF records
// return an empty label on the rbiUrl field, so mapping must match on type.
func dsVaultRecordJSON(uid string) map[string]interface{} {
	return map[string]interface{}{
		"record_uid": uid,
		"type":       "pamCloudResource",
		"title":      "Demo PAM Cloud",
		"notes":      "demo notes",
		"fields": []map[string]interface{}{
			{"type": "rbiUrl", "label": "", "value": json.RawMessage(`["https://console.demo-cloud.example.com"]`)},
			{"type": "accountNumber", "label": "pamCloudResourceAccountIdentifier", "value": json.RawMessage(`["demo-account-123456"]`)},
			{"type": "pamRemoteBrowserSettings", "label": "", "value": json.RawMessage(`[]`)},
		},
	}
}

func newDSConfigRaw(t *testing.T, sch dschema.Schema, pamCloud string) tftypes.Value {
	t.Helper()
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("expected tftypes.Object")
	}
	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		if name == "pam_cloud" {
			vals[name] = tftypes.NewValue(tftypes.String, pamCloud)
		} else {
			vals[name] = tftypes.NewValue(attrType, nil)
		}
	}
	return tftypes.NewValue(objType, vals)
}

func newDSEmptyState(t *testing.T, sch dschema.Schema) tftypes.Value {
	t.Helper()
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("expected tftypes.Object")
	}
	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(objType, vals)
}

func readDS(t *testing.T, d *newpamcloud.PamCloudDataSource, uid string) datasource.ReadResponse {
	t.Helper()
	sch := getDSSchema(t)
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: newDSConfigRaw(t, sch, uid)}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: newDSEmptyState(t, sch)}}
	d.Read(context.Background(), req, &resp)
	return resp
}

func TestRead_DS_Success_MapsFields(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.Contains(cmd, "nsf-get") && strings.Contains(cmd, "--format json") {
			return "ok", dsVaultRecordJSON("uid-abc")
		}
		return "ok", nil
	})
	defer server.Close()

	resp := readDS(t, newConfiguredDataSource(t, server), "uid-abc")
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}

	want := map[string]string{
		"id":                 "uid-abc",
		"pam_cloud":          "uid-abc",
		"title":              "Demo PAM Cloud",
		"notes":              "demo notes",
		"url":                "https://console.demo-cloud.example.com",
		"account_identifier": "demo-account-123456",
	}
	for attr, expected := range want {
		var got string
		if diags := resp.State.GetAttribute(context.Background(), path.Root(attr), &got); diags.HasError() {
			t.Fatalf("GetAttribute(%s) failed: %v", attr, diags)
		}
		if got != expected {
			t.Errorf("%s = %q, want %q", attr, got, expected)
		}
	}
}

func TestRead_DS_NoApiManager(t *testing.T) {
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	if resp := readDS(t, d, "uid-abc"); !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestRead_DS_EmptyRecordUID(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	if resp := readDS(t, newConfiguredDataSource(t, server), ""); !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when pam_cloud is empty")
	}
}

func TestRead_DS_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	if resp := readDS(t, newConfiguredDataSource(t, server), "uid-abc"); !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestRead_DS_ApiError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "nsf-get", nil)
	defer server.Close()

	if resp := readDS(t, newConfiguredDataSource(t, server), "uid-abc"); !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when nsf-get command fails")
	}
}

func TestRead_DS_WrongRecordType(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.Contains(cmd, "nsf-get") {
			rec := dsVaultRecordJSON("uid-abc")
			rec["type"] = "pamDatabase"
			return "ok", rec
		}
		return "ok", nil
	})
	defer server.Close()

	if resp := readDS(t, newConfiguredDataSource(t, server), "uid-abc"); !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record type does not match")
	}
}

func TestRead_DS_NilResponse(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(string, int) (string, interface{}) { return "ok", nil })
	defer server.Close()

	if resp := readDS(t, newConfiguredDataSource(t, server), "uid-abc"); !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when response data is nil")
	}
}
