// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseTeamResource_Create_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--add") {
			return "Team ID: team-uid-456", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Engineering", nil, nil, nil, nil, nil, "Root", nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamResource_Create_WithRestrictFlags(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--add") &&
			strings.Contains(cmd, "--restrict-edit on") && strings.Contains(cmd, "--restrict-share on") {
			return "Team ID: team-789", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "RestrictedTeam", true, true, nil, nil, nil, "Root", nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamResource_Create_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Engineering", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns 500")
	}
}

func TestEnterpriseTeamResource_Create_NoApiManager(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Engineering", nil, nil, nil, nil, nil, "Root", nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseTeamResource_Create_ExtractTeamIdFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--add") {
			return "No team ID in response", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Engineering", nil, nil, nil, nil, nil, "Root", nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when team ID cannot be extracted from response")
	}
}
