// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseTeamResource_Update_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "UpdatedTeam", nil, nil, nil, nil, nil, "Root", nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

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

func TestEnterpriseTeamResource_Update_ManagedCompanyChangeError(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})

	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", "OtherMC"))
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", "MC1"))

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

func TestEnterpriseTeamResource_Update_NoApiManager(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Updated", nil, nil, nil, nil, nil, "Root", nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

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

func TestEnterpriseTeamResource_Update_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Updated", nil, nil, nil, nil, nil, "Root", nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

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
