// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_team"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func teamListResponse(teamUID, name string, users, roles []interface{}) interface{} {
	return []map[string]interface{}{
		{
			"team_uid": teamUID,
			"name":     name,
			"users":    users,
			"roles":    roles,
		},
	}
}

func TestEnterpriseTeamDataSource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json") {
			return "ok", teamListResponse("team-uid-123", "Engineering", []interface{}{}, []interface{}{})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Engineering", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamDataSource_Read_WithUsersAndRoles(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json") {
			return "ok", teamListResponse("team-uid-456", "DevTeam", []interface{}{"user@example.com"}, []interface{}{"Admin"})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("DevTeam", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamDataSource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("NonExistent", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when team not found")
	}
}

func TestEnterpriseTeamDataSource_Read_NoApiManager(t *testing.T) {
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Engineering", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseTeamDataSource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Engineering", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
