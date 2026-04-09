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

// userListResponse returns JSON array for enterprise-info '<id>' -u (single user for Read).
func userListResponse(userID int, email, name, jobTitle, status, node string, teams, roles []interface{}) interface{} {
	return []map[string]interface{}{
		{
			"user_id":   float64(userID),
			"email":     email,
			"name":      name,
			"job_title": jobTitle,
			"status":    status,
			"node":      node,
			"teams":     teams,
			"roles":     roles,
		},
	}
}

func nodesListResponse() interface{} {
	return []map[string]interface{}{
		{"node_id": float64(1), "name": "Root", "parent_node": "", "parent_id": float64(0)},
	}
}

func TestEnterpriseUserResource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json") {
			return "ok", userListResponse(123, "user@example.com", "Dev", "Engineer", "Invited", "Root", []interface{}{}, []interface{}{})
		}
		if strings.Contains(cmd, "enterprise-info -n --format json -v -q") {
			return "ok", nodesListResponse()
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "Dev", "Engineer", nil, nil, "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("999", "old@ex.com", nil, nil, nil, nil, "Root", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read (not found) should not add error: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Read_NoApiManager(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", nil, nil, nil, nil, "Root", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseUserResource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", nil, nil, nil, nil, "Root", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
