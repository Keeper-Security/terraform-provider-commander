// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUpdate_Success_TitleChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Old Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", "https://example.com", nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_Success_SettingsChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	oldSettings := newSettingsValues(
		"config-uid", nil,
		true, false, false, false, false,
		nil, nil, []interface{}{"#login-form"},
		true, true, false,
		float64(2), float64(16), float64(44100),
	)
	newSettings := newSettingsValues(
		"config-uid-new", "cred-uid",
		false, true, true, true, true,
		[]interface{}{"https://allowed.example.com"}, nil, nil,
		false, false, true,
		float64(1), float64(8), float64(48000),
	)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, oldSettings,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, newSettings,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_NoApiManager(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", "https://example.com", nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestUpdate_EmptyId(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"", "Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"", "New Title", "https://example.com", nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when id is empty")
	}
}

func TestUpdate_RecordUpdateFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "record-update", nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Old Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", "https://example.com", nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record-update fails")
	}
}

func TestUpdate_PamRbiEditFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "pam rbi edit", nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	newSettings := newSettingsValues(
		"config-uid", nil,
		true, false, false, false, false,
		nil, nil, nil,
		true, true, false,
		float64(2), float64(16), float64(44100),
	)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, newSettings,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when pam rbi edit fails")
	}
}

func TestUpdate_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Old Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", "https://example.com", nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}

func TestUpdate_NoMutations_NoSettingsChange(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", "https://example.com", nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update should succeed with no mutations: %v", resp.Diagnostics)
	}
	if mock.CommandCount() != 1 {
		t.Errorf("expected 1 command (sync-down only), got %d", mock.CommandCount())
	}
}
