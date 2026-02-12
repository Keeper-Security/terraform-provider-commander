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

// teamListResponse returns JSON array for enterprise-info '<id>' -t (single team for Read).
// Restricts: "W" = edit, "S" = share, "R" = view (enable_privacy_screen).
func teamListResponse(teamUID, name, restricts, node string, users, roles []interface{}) interface{} {
	return []map[string]interface{}{
		{
			"team_uid":  teamUID,
			"name":      name,
			"restricts": restricts,
			"node":      node,
			"users":     users,
			"roles":     roles,
		},
	}
}

func nodesListResponse() interface{} {
	return []map[string]interface{}{
		{"node_id": float64(1), "name": "Root", "parent_node": "", "parent_id": float64(0)},
	}
}

func TestEnterpriseTeamResource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json") {
			return "ok", teamListResponse("team-uid-123", "Engineering", "W S", "Root", []interface{}{}, []interface{}{})
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
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", true, true, nil, nil, nil, "Root", nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamResource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("unknown-uid", "OldTeam", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read (not found) should not add error: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamResource_Read_NoApiManager(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseTeamResource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
