// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush_test

import (
	"context"
	"strings"
	"testing"

	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseScimPushResource_Create_Success_Google(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "scim push") && strings.Contains(cmd, "google") {
			return "ok", nil
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "scim-1", "google", "record-uid-1", true))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseScimPushResource_Create_Success_AD(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) { return "ok", nil })
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "scim-2", "ad", "record-uid-2", false))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseScimPushResource_Create_Success_Record(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) { return "ok", nil })
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "scim-3", "record", "record-uid-3", true))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseScimPushResource_Create_InvalidSource_Error(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) { return "ok", nil })
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "scim-1", "invalid_source", "record-uid-1", true))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when source is invalid")
	}
	var found bool
	for _, d := range resp.Diagnostics.Errors() {
		if strings.Contains(d.Detail(), "Invalid source") || strings.Contains(d.Detail(), "google, ad, record") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected Invalid source message in diagnostics, got: %v", resp.Diagnostics)
	}
}

func TestEnterpriseScimPushResource_Create_NoApiManager(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(nil, "scim-1", "google", "record-uid-1", true))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: rawPlan}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}
