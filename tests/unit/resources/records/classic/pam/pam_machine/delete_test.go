// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDelete_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", resp.Diagnostics)
	}
}

func TestDelete_NoApiManager(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestDelete_ApiFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "rm", nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when rm command fails")
	}
}

func TestDelete_SyncDownFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"sync failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when sync-down fails")
	}
}
