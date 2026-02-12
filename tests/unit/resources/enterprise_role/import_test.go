// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseRoleResource_ImportState_EmptyID(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestEnterpriseRoleResource_ImportState_ManagedCompanyEmptyRole(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	req := resource.ImportStateRequest{ID: "Company,"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when role part is empty in managed company format")
	}
}

func TestEnterpriseRoleResource_ImportState_NoApiManager(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	req := resource.ImportStateRequest{ID: "Admin"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is nil")
	}
}

func TestEnterpriseRoleResource_ImportState_Success_RoleOnly(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil))
	req := resource.ImportStateRequest{ID: "Admin"}
	var resp resource.ImportStateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: emptyRaw}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseRoleResource_ImportState_Success_WithManagedCompany(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil, nil))
	req := resource.ImportStateRequest{ID: "Test Company,Admin"}
	var resp resource.ImportStateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: emptyRaw}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}
