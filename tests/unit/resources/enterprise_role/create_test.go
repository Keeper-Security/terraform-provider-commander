// Copyright (c) HashiCorp, Inc.
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

func TestEnterpriseRoleResource_Create_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "--add") {
			return "Role ID : 456", nil
		}
		if strings.Contains(cmd, "enterprise-info -n --format json") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Admin", "Root", nil, nil, nil, nil, nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleResource_Create_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns 500")
	}
}

func TestEnterpriseRoleResource_Create_NoApiManager(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Admin", "Root", nil, nil, nil, nil, nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseRoleResource_Create_ExtractRoleIdFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "--add") {
			return "No role ID in response", nil
		}
		if strings.Contains(cmd, "enterprise-info -n --format json") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Admin", "Root", nil, nil, nil, nil, nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when role ID cannot be extracted from response")
	}
}

func TestEnterpriseRoleResource_Create_TeamsAndManagingNodesMutualExclusivity(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Both teams and managing_nodes set -> validation error
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil)
	managingNodeElem := tftypes.NewValue(managingNodesElemType, map[string]tftypes.Value{
		"privileges": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"cascade":    tftypes.NewValue(tftypes.Bool, false),
	})
	planVals := newPlanStateValues(nil, "Admin", "Root", nil, []interface{}{"Team1"}, map[string]tftypes.Value{"Node1": managingNodeElem}, nil, nil)
	rawPlan := tftypes.NewValue(objType, planVals)
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when both teams and managing_nodes are set")
	}
}
