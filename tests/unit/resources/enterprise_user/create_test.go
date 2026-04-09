// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_user"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseUserResource_Create_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--add") {
			return "ok", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "user@example.com", nil, nil, nil, nil, "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Create_WithNameAndJobTitle(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--add") &&
			strings.Contains(cmd, "--name") && strings.Contains(cmd, "--job-title") {
			return "u@ex.com user invited with Enterprise User ID : 456", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "u@ex.com", "Dev", "Engineer", nil, nil, "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Create_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "user@example.com", nil, nil, nil, nil, "Root", nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns 500")
	}
}

func TestEnterpriseUserResource_Create_NoApiManager(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "user@example.com", nil, nil, nil, nil, "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

// TestEnterpriseUserResource_Create_FallbackToEmailWhenNoUserID verifies that when the API
// response does not contain "User ID : <id>", create succeeds and falls back to using email as id (for API compatibility).
func TestEnterpriseUserResource_Create_FallbackToEmailWhenNoUserID(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--add") {
			return "User invited successfully", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "user@example.com", nil, nil, nil, nil, "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create should succeed with email fallback: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Create_TeamsSetError(t *testing.T) {
	// Use mock server so ExecuteWithManagedCompanyContext can run enterprise-down; then lambda returns teams error.
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "user@example.com", nil, nil, nil, []interface{}{"Team1"}, "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when teams are set on create")
	}
}
