// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRead_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") && strings.Contains(cmd, "--format json") {
			return "ok", vaultRecordGetJSON("uid-abc", "My RBI", "https://example.com", "notes", "")
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-abc", "My RBI", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestRead_NoApiManager(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-abc", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestRead_EmptyId(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when id is empty")
	}
}

func TestRead_ResourceNotFound_RemovesFromState(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, idx int) (int, []byte) {
		if strings.Contains(cmd, "get") {
			return http.StatusInternalServerError, []byte(`{"message":"record not found"}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-missing", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error for not-found (resource removed from state): %v", resp.Diagnostics)
	}
}

func TestRead_ApiError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, idx int) (int, []byte) {
		if strings.Contains(cmd, "get") {
			return http.StatusInternalServerError, []byte(`{"message":"command execution failed"}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-abc", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when get returns 500")
	}
}

func TestRead_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-abc", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestRead_WrongRecordType(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") {
			return "ok", map[string]interface{}{
				"record_uid": "uid-abc",
				"type":       "pamMachine",
				"title":      "Wrong Type",
				"notes":      "",
				"fields":     []interface{}{},
			}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-abc", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record type does not match")
	}
}
