// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/managed_company"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestManagedCompanyResource_Delete_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("1169425105420462", "Test Company", "Root", nil, "business", nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.DeleteResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", resp.Diagnostics)
	}
}

func TestManagedCompanyResource_Delete_NoApiManager(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("1169425105420462", "Test Company", "Root", nil, "business", nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestManagedCompanyResource_Delete_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("1169425105420462", "Test Company", "Root", nil, "business", nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

func TestManagedCompanyResource_Delete_StateGetError(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, _ := getSchema(t)
	wrongTypes := map[string]tftypes.Type{
		"id": tftypes.Number, "name": tftypes.String, "node": tftypes.String,
		"seats": tftypes.Number, "plan": tftypes.String, "file_plan": tftypes.String,
		"add_ons": tftypes.Set{ElementType: tftypes.String},
	}
	rawState := tftypes.NewValue(tftypes.Object{AttributeTypes: wrongTypes}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.Number, 1169425105420462), "name": tftypes.NewValue(tftypes.String, "Test"),
		"node": tftypes.NewValue(tftypes.String, "Root"), "seats": tftypes.NewValue(tftypes.Number, nil),
		"plan": tftypes.NewValue(tftypes.String, "business"), "file_plan": tftypes.NewValue(tftypes.String, nil),
		"add_ons": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when state get fails")
	}
}
