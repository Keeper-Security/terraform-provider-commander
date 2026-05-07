// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_role"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func roleListResponse(roleID int, name string, users, teams []interface{}) interface{} {
	return []map[string]interface{}{
		{
			"role_id": float64(roleID),
			"name":    name,
			"users":   users,
			"teams":   teams,
		},
	}
}

func TestEnterpriseRoleDataSource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json") {
			return "ok", roleListResponse(123, "Admin", []interface{}{}, []interface{}{})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Admin", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleDataSource_Read_WithUsersAndTeams(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json") {
			return "ok", roleListResponse(456, "Developer", []interface{}{"user@example.com"}, []interface{}{"Team1"})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Developer", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleDataSource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") {
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
		t.Error("expected diagnostics when role not found")
	}
}

func TestEnterpriseRoleDataSource_Read_NoApiManager(t *testing.T) {
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Admin", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseRoleDataSource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Admin", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
