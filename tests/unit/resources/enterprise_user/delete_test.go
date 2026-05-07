// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_user"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseUserResource_Delete_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", nil, nil, nil, nil, "Root", nil, "Invited"))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Delete_NoApiManager(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", nil, nil, nil, nil, "Root", nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseUserResource_Delete_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", nil, nil, nil, nil, "Root", nil, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

func TestEnterpriseUserResource_Delete_StateGetError(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, _ := getSchema(t)
	wrongTypes := map[string]tftypes.Type{
		"id": tftypes.Number, "email": tftypes.String, "name": tftypes.String, "job_title": tftypes.String,
		"roles": tftypes.Set{ElementType: tftypes.String}, "teams": tftypes.Set{ElementType: tftypes.String},
		"node": tftypes.String, "managed_company": tftypes.String, "status": tftypes.String,
	}
	rawState := tftypes.NewValue(tftypes.Object{AttributeTypes: wrongTypes}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.Number, 123), "email": tftypes.NewValue(tftypes.String, "u@ex.com"),
		"name": tftypes.NewValue(tftypes.String, nil), "job_title": tftypes.NewValue(tftypes.String, nil),
		"roles": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"teams": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"node":  tftypes.NewValue(tftypes.String, "Root"), "managed_company": tftypes.NewValue(tftypes.String, nil),
		"status": tftypes.NewValue(tftypes.String, nil),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when state get fails")
	}
}
