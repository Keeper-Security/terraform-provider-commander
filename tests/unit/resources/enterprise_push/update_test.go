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

func TestEnterprisePushResource_Update_Success_WithAddedEmail(t *testing.T) {
	filePath := createTempJSONFile(t)
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-push") && strings.Contains(cmd, "--email=") {
			return "ok", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Plan: add one email (carol). State: alice, bob.
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"alice@ex.com", "bob@ex.com", "carol@ex.com"}, []interface{}{"Engineering"}))
	rawState := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"alice@ex.com", "bob@ex.com"}, []interface{}{"Engineering"}))

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

func TestEnterprisePushResource_Update_NoPushWhenOnlyRemovals(t *testing.T) {
	filePath := createTempJSONFile(t)
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Plan: remove one email. State has two emails. No new targets → no push.
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"alice@ex.com"}, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"alice@ex.com", "bob@ex.com"}, nil))

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

func TestEnterprisePushResource_Update_NoApiManager(t *testing.T) {
	filePath := createTempJSONFile(t)
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"a@ex.com"}, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"a@ex.com"}, nil))

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

func TestEnterprisePushResource_Update_FileReadFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("id1", "/nonexistent/file.json", "sha256abc", []interface{}{"a@ex.com"}, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("id1", "/nonexistent/file.json", "sha256abc", []interface{}{"a@ex.com"}, nil))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when file cannot be read")
	}
}

func TestEnterprisePushResource_Update_ApiError(t *testing.T) {
	filePath := createTempJSONFile(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"a@ex.com", "b@ex.com"}, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("id1", filePath, "sha256abc", []interface{}{"a@ex.com"}, nil))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns 500")
	}
}
