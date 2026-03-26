// Copyright (c) Keeper Security, Inc.
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

func TestEnterpriseNodeResource_Read_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info -n -v --format json --node") {
			return "ok", []map[string]interface{}{
				{"node_id": float64(123), "name": "TestNode", "parent_node": "Root", "parent_id": float64(0)},
			}
		}
		if strings.Contains(cmd, "enterprise-info -n --format json -v -q") {
			return "ok", []map[string]interface{}{
				{"node_id": float64(123), "name": "TestNode", "parent_node": "Root", "parent_id": float64(0)},
				{"node_id": float64(1), "name": "Root", "parent_node": "", "parent_id": float64(0)},
			}
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseNodeResource_Read_NotFound(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info -n -v --format json --node") {
			return "ok", []interface{}{} // empty = node not found
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("999", "Node", "Root", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read (not found) should not add error: %v", resp.Diagnostics)
	}
}

func TestEnterpriseNodeResource_Read_NoApiManager(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseNodeResource_Read_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

// TestEnterpriseNodeResource_Read_RestoreParentFormatError triggers RestoreUserInputFormatForNode
// to fail by returning invalid data for the node-list command (enterprise-info -n --format json -v -q).
func TestEnterpriseNodeResource_Read_RestoreParentFormatError(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info -n -v --format json --node") {
			return "ok", []map[string]interface{}{
				{"node_id": float64(123), "name": "TestNode", "parent_node": "Root", "parent_id": float64(0), "isolated": false},
			}
		}
		if strings.Contains(cmd, "enterprise-info -n --format json -v -q") {
			// Return non-array data so ParseNodesResponse fails and RestoreUserInputFormatForNode returns error
			return "ok", "not-an-array"
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("123", "Node", "Root", nil, nil))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when RestoreUserInputFormatForNode fails")
	}
}

// TestEnterpriseNodeResource_Read_StateGetError passes state with wrong type so State.Get adds diagnostics.
func TestEnterpriseNodeResource_Read_StateGetError(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, _ := getSchema(t)
	wrongTypes := map[string]tftypes.Type{
		"id": tftypes.Number, "name": tftypes.String, "parent": tftypes.String,
		"toggle_isolated": tftypes.Bool, "managed_company": tftypes.String,
	}
	rawState := tftypes.NewValue(tftypes.Object{AttributeTypes: wrongTypes}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.Number, 123), "name": tftypes.NewValue(tftypes.String, "Node"),
		"parent": tftypes.NewValue(tftypes.String, "Root"), "toggle_isolated": tftypes.NewValue(tftypes.Bool, nil),
		"managed_company": tftypes.NewValue(tftypes.String, nil),
	})

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when state get fails")
	}
}
