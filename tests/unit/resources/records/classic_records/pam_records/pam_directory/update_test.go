// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_directory"
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

	hostVals := newHostnameValues("host.com", float64(636))
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Old Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

func TestUpdate_Success_HostnameChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	oldHost := newHostnameValues("old-host.com", float64(389))
	newHost := newHostnameValues("new-host.com", float64(636))
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", oldHost, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", newHost, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

func TestUpdate_Success_OptionalFieldsChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals,
		true, "old-domain.com", nil,
		"old-dir-id", "active_directory", "OU=OldUsers",
		"old-group", "us-west-1",
		"old notes", nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals,
		false, "new-domain.com", nil,
		"new-dir-id", "openldap", "OU=NewUsers",
		"new-group", "eu-west-1",
		"new notes", nil,
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
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"", "New Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Old Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

func TestUpdate_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Old Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "New Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

func TestUpdate_NoMutations(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", float64(636))
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

func TestUpdate_NotesClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, "old notes", nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

func TestUpdate_FolderChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, "old-folder",
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, "new-folder",
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
