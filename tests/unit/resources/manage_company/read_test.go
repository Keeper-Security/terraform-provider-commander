// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/manage_company"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// companyListResponse returns JSON array for msp-info -m '<id>' (single company for Read).
// ManageCompanyResponse: company_id, company_name, node, node_name, plan, storage, addons, allocated.
func companyListResponse(companyID int, companyName, nodeName, plan, storage string, allocated int, addons []interface{}) interface{} {
	return []map[string]interface{}{
		{
			"company_id":   float64(companyID),
			"company_name": companyName,
			"node":         "Root",
			"node_name":    nodeName,
			"plan":         plan,
			"storage":      storage,
			"allocated":    float64(allocated),
			"addons":       addons,
		},
	}
}

func TestManageCompanyResource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "msp-down") {
			return "ok", nil
		}
		if strings.Contains(cmd, "msp-info -m") && strings.Contains(cmd, "--format json") {
			return "ok", companyListResponse(1169425105420462, "Test Company", "Root", "business", "100GB", 5, []interface{}{})
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyVals := newPlanStateValues(nil, nil, nil, nil, nil, nil, nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues("1169425105420462", "Test Company", "Root", nil, "business", nil, nil))
	emptyState := tftypes.NewValue(objType, emptyVals)

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestManageCompanyResource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "msp-down") {
			return "ok", nil
		}
		if strings.Contains(cmd, "msp-info -m") {
			return "ok", []interface{}{}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("999", "Unknown", "Root", nil, "business", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read (not found) should not add error: %v", resp.Diagnostics)
	}
}

func TestManageCompanyResource_Read_NoApiManager(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("1169425105420462", "Test Company", "Root", nil, "business", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestManageCompanyResource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("1169425105420462", "Test Company", "Root", nil, "business", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
