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
	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func getResSchema(t *testing.T) rschema.Schema {
	t.Helper()
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema
}

func newConfiguredRes(t *testing.T, server *httptest.Server) *pamconfiguration.PamConfigurationResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func newResStateRaw(t *testing.T, sch rschema.Schema, id string) tftypes.Value {
	t.Helper()
	tfType := sch.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", tfType)
	}
	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		switch name {
		case "id":
			vals[name] = tftypes.NewValue(tftypes.String, id)
		case "environment":
			vals[name] = tftypes.NewValue(tftypes.String, "local")
		case "title":
			vals[name] = tftypes.NewValue(tftypes.String, "Test")
		case "gateway":
			vals[name] = tftypes.NewValue(tftypes.String, "gw-uid")
		case "application_folder":
			vals[name] = tftypes.NewValue(tftypes.String, "folder-uid")
		default:
			vals[name] = tftypes.NewValue(attrType, nil)
		}
	}
	return tftypes.NewValue(objType, vals)
}

func pamConfigListJSON() interface{} {
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

func TestResRead_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "pam config") && strings.Contains(cmd, "list") {
			return "ok", pamConfigListJSON()
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "cfg-uid-123")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestResRead_NoApiManager(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "cfg-uid-123")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestResRead_EmptyId(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when id is empty")
	}
}

func TestResRead_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "cfg-uid-123")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestResRead_ApiError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, idx int) (int, []byte) {
		if strings.Contains(cmd, "pam config") {
			return http.StatusInternalServerError, []byte(`{"message":"command execution failed"}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "cfg-uid-123")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

func TestResRead_NilResponse_RemovesFromState(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "cfg-uid-123")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error for nil response (resource removed): %v", resp.Diagnostics)
	}
}

func TestResRead_ResourceNotFound_RemovesFromState(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, idx int) (int, []byte) {
		if strings.Contains(cmd, "pam config") {
			return http.StatusInternalServerError, []byte(`{"message":"record not found"}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)
	rawState := newResStateRaw(t, sch, "cfg-missing")

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error for not-found: %v", resp.Diagnostics)
	}
}
