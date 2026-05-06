// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func getDSSchema(t *testing.T) dschema.Schema {
	t.Helper()
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

func newConfiguredDS(t *testing.T, server *httptest.Server) *pamconfiguration.PamConfigurationDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func newDSConfigRaw(t *testing.T, sch dschema.Schema, pamConfig string) tftypes.Value {
	t.Helper()
	tfType := sch.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", tfType)
	}
	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		if name == "pam_configuration" {
			vals[name] = tftypes.NewValue(tftypes.String, pamConfig)
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

func dsPamConfigListJSON() interface{} {
	return map[string]interface{}{
		"uid":          "cfg-uid-123",
		"name":         "Test Config",
		"config_type":  "pamNetworkConfiguration",
		"gateway_uid":  "gw-uid",
		"gateway_name": "My Gateway",
		"shared_folder": map[string]interface{}{
			"uid":  "folder-uid",
			"name": "PAM Folder",
		},
		"allowed_settings": map[string]interface{}{
			"connections":                       true,
			"tunneling":                         false,
			"rotation":                          true,
			"remote_browser_isolation":          false,
			"connections_recording":             true,
			"typescript_recording":              false,
			"ai_threat_detection":               false,
			"ai_terminate_session_on_detection": false,
		},
		"fields": map[string]interface{}{
			"networkId":   []string{"DC-East-1"},
			"networkCIDR": []string{"10.0.0.0/16"},
		},
	}
}

func TestDSRead_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "pam config") && strings.Contains(cmd, "list") {
			return "ok", dsPamConfigListJSON()
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDS(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "cfg-uid-123")
	emptyState := newDSEmptyState(t, sch)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestDSRead_NoApiManager(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "cfg-uid-123")

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestDSRead_EmptyConfigUID(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	d := newConfiguredDS(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "")

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when pam_configuration is empty")
	}
}

func TestDSRead_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	d := newConfiguredDS(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "cfg-uid-123")

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestDSRead_ApiError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "pam config", nil)
	defer server.Close()

	d := newConfiguredDS(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "cfg-uid-123")

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

func TestDSRead_NilResponse(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDS(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "cfg-uid-123")

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when response data is nil")
	}
}

func TestDSRead_GcpConfig(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "pam config") && strings.Contains(cmd, "list") {
			return "ok", map[string]interface{}{
				"uid":          "cfg-gcp-001",
				"name":         "GCP Config",
				"config_type":  "pamGcpConfiguration",
				"gateway_uid":  "gw-gcp",
				"gateway_name": "GCP Gateway",
				"shared_folder": map[string]interface{}{
					"uid":  "sf-uid",
					"name": "GCP Folder",
				},
				"allowed_settings": map[string]interface{}{
					"connections":                       true,
					"tunneling":                         true,
					"rotation":                          false,
					"remote_browser_isolation":          false,
					"connections_recording":             false,
					"typescript_recording":              false,
					"ai_threat_detection":               false,
					"ai_terminate_session_on_detection": false,
				},
				"fields": map[string]interface{}{
					"pamGcpId":             []string{"my-project"},
					"pamServiceAccountKey": []string{`{"type":"service_account"}`},
					"pamGoogleAdminEmail":  []string{"admin@example.com"},
					"pamGcpRegionName":     []string{"us-central1"},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDS(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "cfg-gcp-001")
	emptyState := newDSEmptyState(t, sch)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed for GCP config: %v", resp.Diagnostics)
	}
}
