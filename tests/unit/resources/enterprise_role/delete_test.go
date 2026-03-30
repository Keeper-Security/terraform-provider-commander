// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseRoleResource_Delete_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleResource_Delete_NoApiManager(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseRoleResource_Delete_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Admin", "Root", nil, nil, nil, nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

func TestEnterpriseRoleResource_Delete_StateGetError(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, _ := getSchema(t)
	wrongTypes := map[string]tftypes.Type{
		"id": tftypes.Number, "name": tftypes.String, "node": tftypes.String,
		"users": tftypes.Set{ElementType: tftypes.String}, "teams": tftypes.Set{ElementType: tftypes.String},
		"managing_nodes":       tftypes.Map{ElementType: managingNodesElemType},
		"enforcement_policies": tftypes.Map{ElementType: tftypes.String}, "managed_company": tftypes.String,
	}
	rawState := tftypes.NewValue(tftypes.Object{AttributeTypes: wrongTypes}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.Number, 123), "name": tftypes.NewValue(tftypes.String, "Admin"),
		"node":                 tftypes.NewValue(tftypes.String, "Root"),
		"users":                tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"teams":                tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"managing_nodes":       tftypes.NewValue(tftypes.Map{ElementType: managingNodesElemType}, nil),
		"enforcement_policies": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		"managed_company":      tftypes.NewValue(tftypes.String, nil),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when state get fails")
	}
}
