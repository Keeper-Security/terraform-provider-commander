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

func TestEnterpriseNodeResource_Create_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-node --add") {
			return "Node ID: 123", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "TestNode", "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

// TestEnterpriseNodeResource_Create_NoParent covers addNodeBasicAttributes when Parent is null (no --parent in command).
func TestEnterpriseNodeResource_Create_NoParent(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-node --add") && !strings.Contains(cmd, "--parent") {
			return "Node ID: 111", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// Parent null -> no --parent in command
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "NoParentNode", nil, nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseNodeResource_Create_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Test", "Root", nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns 500")
	}
}

func TestEnterpriseNodeResource_Create_NoApiManager(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Test", "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

// TestEnterpriseNodeResource_Create_WithParentAndManagedCompany verifies create with parent and managed_company set.
// Parent is passed through as-is (no special "root" substitution).
func TestEnterpriseNodeResource_Create_WithParentAndManagedCompany(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-node --add") && strings.Contains(cmd, "--parent 'MyCompany'") {
			return "Node ID: 456", nil
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Child", "MyCompany", nil, "MyCompany"))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

// TestEnterpriseNodeResource_Create_ToggleIsolatedTrueRejected verifies that create returns an error
// when toggle_isolated is set to true (not supported on create).
func TestEnterpriseNodeResource_Create_ToggleIsolatedTrueRejected(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "IsolatedNode", "Root", true, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics when toggle_isolated is true on create")
	}
	var found bool
	for _, d := range resp.Diagnostics {
		if strings.Contains(d.Detail(), "toggle_isolated is not supported in create") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error message about toggle_isolated; got: %v", resp.Diagnostics)
	}
}

func TestEnterpriseNodeResource_Create_ExtractNodeIdFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-node --add") {
			return "No node ID in this response", nil // message doesn't contain parseable node ID
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Test", "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when node ID cannot be extracted from response")
	}
}

// TestEnterpriseNodeResource_Create_AlreadyExists reproduces the reported
// bug: Commander's `enterprise-node --add` on an already-existing name is a
// harmless no-op at the CLI level, but Service Mode's generic response
// parser surfaces the resulting warning as an HTTP 409. Create() must not
// hard-fail with the raw API error; it should look up the existing node and
// guide the practitioner toward `terraform import` instead.
func TestEnterpriseNodeResource_Create_AlreadyExists(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") {
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(4821),
					"name":        "Engineering",
					"parent_node": "Root",
					"parent_id":   float64(0),
					"isolated":    false,
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServerWithResultHook(mock, responseForCommand, func(cmd string, idx int) (int, []byte) {
		if strings.Contains(cmd, "enterprise-node --add") {
			return http.StatusConflict, []byte(`{"message":"Node 'Engineering' already exists: Skipping."}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "Engineering", "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected a diagnostics error when the node already exists")
	}
	found := false
	for _, d := range resp.Diagnostics.Errors() {
		if strings.Contains(d.Detail(), "terraform import commander_enterprise_node") && strings.Contains(d.Detail(), "4821") {
			found = true
		}
		if strings.Contains(d.Detail(), "status 409") {
			t.Errorf("raw API error text leaked into diagnostic instead of import guidance: %s", d.Detail())
		}
	}
	if !found {
		t.Errorf("expected import-guidance diagnostic naming the existing ID 4821, got: %v", resp.Diagnostics)
	}
}
