// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseNodeResource_Update_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "UpdatedNode", "Root", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))

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

func TestEnterpriseNodeResource_Update_ManagedCompanyChangeError(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})

	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, "OtherMC"))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, "MC1"))

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

func TestEnterpriseNodeResource_Update_NoApiManager(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "Updated", "Root", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))

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

func TestEnterpriseNodeResource_Update_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "Updated", "Root", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))

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

func TestEnterpriseNodeResource_Update_ParentChangeWithManagedCompany(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		// Parent change with parent == managed company -> command uses --parent 'root'
		if strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--parent 'root'") {
			return "ok", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Change parent to same as managed company -> value becomes "root"
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "MyCompany", nil, "MyCompany"))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, "MyCompany"))

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

func TestEnterpriseNodeResource_Update_ParentChangeOnly(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		// Parent change without managed company -> --parent 'NewParent'
		if strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--parent 'NewParent'") {
			return "ok", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "NewParent", nil, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))

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

func TestEnterpriseNodeResource_Update_ToggleIsolatedChange(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Change toggle_isolated from false/null to true
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", true, nil))
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", false, nil))

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
