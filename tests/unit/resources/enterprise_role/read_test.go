// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func roleListResponse(roleID int, name, node string, users, teams []interface{}, managedNodesPerms []interface{}, enforcements map[string]interface{}) interface{} {
	return []map[string]interface{}{
		{
			"role_id":                   float64(roleID),
			"name":                      name,
			"node":                      node,
			"users":                     users,
			"teams":                     teams,
			"managed_nodes_permissions": managedNodesPerms,
			"enforcements":              enforcements,
		},
	}
}

func nodesListResponse() interface{} {
	return []map[string]interface{}{
		{"node_id": float64(1), "name": "Root", "parent_node": "", "parent_id": float64(0)},
	}
}

func TestEnterpriseRoleResource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json") {
			return "ok", roleListResponse(123, "Admin", "Root", []interface{}{}, []interface{}{}, []interface{}{}, map[string]interface{}{})
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
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Admin", "Root", nil, nil, nil, nil, nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleResource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("999", "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read (not found) should not add error: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleResource_Read_NoApiManager(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseRoleResource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
