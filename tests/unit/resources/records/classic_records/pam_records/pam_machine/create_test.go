// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_machine"
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

	hostVals := newHostnameValues("localhost.com", float64(22))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "My Machine",
		hostVals,
		nil, nil, nil,
		nil, nil,
		nil, nil,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestCreate_Success_WithAllFields(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "sync-down") {
			return "ok", nil
		}
		if strings.Contains(cmd, "record-add") {
			return "ok", map[string]interface{}{"record_uid": "new-uid-456"}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("10.0.0.1", float64(3389))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Full Machine",
		hostVals,
		"Linux", "my-instance", "i-12345",
		"AWS", "us-east-1",
		"some notes", "folder-uid",
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestCreate_NoApiManager(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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

	hostVals := newHostnameValues("host.com", nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
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

	hostVals := newHostnameValues("host.com", nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
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

	hostVals := newHostnameValues("host.com", nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		nil, "Title", hostVals, nil, nil, nil, nil, nil, nil, nil,
	))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record_uid is not in response")
	}
}

func vaultRecordGetJSON(uid, title, notes string, hostname, port string) interface{} {
	rec := map[string]interface{}{
		"record_uid": uid,
		"type":       "pamMachine",
		"title":      title,
		"notes":      notes,
		"fields": []map[string]interface{}{
			{
				"type":  "pamHostname",
				"label": "pamHostname",
				"value": json.RawMessage(`[{"hostName":"` + hostname + `","port":"` + port + `"}]`),
			},
		},
	}
	return rec
}
