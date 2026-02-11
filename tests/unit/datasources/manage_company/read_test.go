// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/manage_company"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func companyListResponse(companyID int, companyName, node, plan, storage string) interface{} {
	return []map[string]interface{}{
		{
			"company_id":   float64(companyID),
			"company_name": companyName,
			"node":         node,
			"plan":         plan,
			"storage":      storage,
		},
	}
}

func manageCompanyResponseForCommand(cmd string, _ int, companyID int, companyName, node, plan, storage string) (string, interface{}) {
	if cmd == "msp-down" {
		return "ok", nil
	}
	if strings.Contains(cmd, "msp-info") && strings.Contains(cmd, "-m") && strings.Contains(cmd, "--format json") {
		return "ok", companyListResponse(companyID, companyName, node, plan, storage)
	}
	return "ok", nil
}

func TestManageCompanyDataSource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		return manageCompanyResponseForCommand(cmd, idx, 100, "Acme Corp", "node-1", "business", "1TB")
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Acme Corp"))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestManageCompanyDataSource_Read_ById(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		return manageCompanyResponseForCommand(cmd, idx, 200, "Other Company", "node-2", "enterprise", "5TB")
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("200"))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestManageCompanyDataSource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if cmd == "msp-down" {
			return "ok", nil
		}
		if strings.Contains(cmd, "msp-info") && strings.Contains(cmd, "-m") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("NonExistent"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when managed company not found")
	}
}

func TestManageCompanyDataSource_Read_NoApiManager(t *testing.T) {
	d := managecompany.NewManageCompanyDataSource().(*managecompany.ManageCompanyDataSource)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Acme Corp"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestManageCompanyDataSource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("Acme Corp"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
