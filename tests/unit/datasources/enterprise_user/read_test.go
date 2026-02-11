// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_user"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func userListResponse(userID int, name, email, jobTitle, status string, roles, teams []interface{}) interface{} {
	return []map[string]interface{}{
		{
			"user_id":   float64(userID),
			"name":      name,
			"email":     email,
			"job_title": jobTitle,
			"status":    status,
			"roles":     roles,
			"teams":     teams,
		},
	}
}

func TestEnterpriseUserDataSource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json") {
			return "ok", userListResponse(1001, "Jane Doe", "jane@example.com", "Engineer", "active", []interface{}{}, []interface{}{})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("jane@example.com", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserDataSource_Read_WithRolesAndTeams(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json") {
			return "ok", userListResponse(1002, "John Smith", "john@example.com", "Admin", "active", []interface{}{"Admin"}, []interface{}{"Engineering"})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("john@example.com", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserDataSource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("nonexistent@example.com", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when user not found")
	}
}

func TestEnterpriseUserDataSource_Read_NoApiManager(t *testing.T) {
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("jane@example.com", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseUserDataSource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("jane@example.com", nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
