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

	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestCreate_Success_WithoutSettings(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "sync-down") {
			return "ok", nil
		}
		if strings.Contains(cmd, "record-add") {
			return "ok", map[string]interface{}{"record_uid": "new-uid-123"}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "My RBI", "https://example.com", nil, nil, nil,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestCreate_Success_WithSettings(t *testing.T) {
	mock := &helpers.CommandServer{}
	callCount := 0
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		callCount++
		if strings.Contains(cmd, "sync-down") {
			return "ok", nil
		}
		if strings.Contains(cmd, "record-add") {
			return "ok", map[string]interface{}{"record_uid": "new-uid-456"}
		}
		if strings.Contains(cmd, "pam rbi edit") {
			return "ok", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	settingsVals := newSettingsValues(
		"config-uid", nil,
		true, false, true, false, false,
		nil, []interface{}{"https://res.example.com"}, nil,
		true, true, false,
		float64(2), float64(16), float64(44100),
	)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "My RBI", "https://example.com", "some notes", "folder-uid",
		settingsVals,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestCreate_NoApiManager(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", "https://example.com", nil, nil, nil,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestCreate_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestCreate_RecordAddFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "record-add", func(cmd string, idx int) (string, interface{}) {
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record-add fails")
	}
}

func TestCreate_RecordAddNoUID(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "record-add") {
			return "ok", map[string]interface{}{"other_field": "value"}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", "https://example.com", nil, nil, nil,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record_uid is not in response")
	}
}

func TestCreate_PamRbiEditFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(
		mock,
		func(cmd string, idx int) (string, interface{}) {
			if strings.Contains(cmd, "record-add") {
				return "ok", map[string]interface{}{"record_uid": "uid-789"}
			}
			return "ok", nil
		},
		func(cmd string, idx int) (int, []byte) {
			if strings.Contains(cmd, "pam rbi edit") {
				return http.StatusInternalServerError, []byte(`{"message":"pam rbi edit failed"}`)
			}
			return 0, nil
		},
	)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	settingsVals := newSettingsValues(
		"config-uid", nil,
		true, false, false, false, false,
		nil, nil, nil,
		true, true, false,
		float64(2), float64(16), float64(44100),
	)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", "https://example.com", nil, nil,
		settingsVals,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when pam rbi edit fails")
	}
}

func vaultRecordGetJSON(uid, title, url, notes, folder string) interface{} {
	rec := map[string]interface{}{
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
	if folder != "" {
		rec["folder"] = folder
	}
	return rec
}
