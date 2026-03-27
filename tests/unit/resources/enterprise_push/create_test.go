// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterprisePushResource_Create_Success_WithEmail(t *testing.T) {
	filePath := createTempJSONFile(t)
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-push") && strings.Contains(cmd, "FILEDATA") && strings.Contains(cmd, "--email=") {
			return "ok", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, filePath, nil, []interface{}{"user@example.com"}, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterprisePushResource_Create_Success_WithTeam(t *testing.T) {
	filePath := createTempJSONFile(t)
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, filePath, nil, nil, []interface{}{"Engineering"}))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterprisePushResource_Create_NoEmailOrTeam_Error(t *testing.T) {
	filePath := createTempJSONFile(t)
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Neither email nor team set (both null)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, filePath, nil, nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when neither email nor team is specified")
	}
	// Ensure our specific error message is present
	var found bool
	for _, d := range resp.Diagnostics.Errors() {
		if strings.Contains(d.Detail(), "At least one of email or team must be specified") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'At least one of email or team must be specified' in diagnostics, got: %v", resp.Diagnostics)
	}
}

func TestEnterprisePushResource_Create_NoApiManager(t *testing.T) {
	filePath := createTempJSONFile(t)
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, filePath, nil, []interface{}{"user@example.com"}, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterprisePushResource_Create_FileReadFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Non-existent file
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "/nonexistent/path/to/file.json", nil, []interface{}{"user@example.com"}, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when file cannot be read")
	}
}

func TestEnterprisePushResource_Create_ApiError(t *testing.T) {
	filePath := createTempJSONFile(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, filePath, nil, []interface{}{"user@example.com"}, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns 500")
	}
}
