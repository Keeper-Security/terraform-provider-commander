// Copyright (c) HashiCorp, Inc.
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

func TestEnterpriseUserResource_Update_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "Updated Name", "Lead", nil, nil, "Root", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "Dev", "Engineer", nil, nil, "Root", nil, "Invited"))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseUserResource_Update_ManagedCompanyChangeError(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key", IsMspAccount: false}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})

	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "u@ex.com", nil, nil, nil, nil, "Root", "OtherMC", nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "u@ex.com", nil, nil, nil, nil, "Root", "MC1", nil))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when managed_company is changed")
	}
}

func TestEnterpriseUserResource_Update_NoApiManager(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "New", nil, nil, nil, "Root", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "Old", nil, nil, nil, "Root", nil, nil))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseUserResource_Update_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "New", nil, nil, nil, "Root", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "user@example.com", "Old", nil, nil, nil, "Root", nil, nil))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}
