// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/managed_company"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestManagedCompanyResource_ImportState_EmptyID(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestManagedCompanyResource_ImportState_NoApiManager(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	req := resource.ImportStateRequest{ID: "Test Company"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is nil")
	}
}

func TestManagedCompanyResource_ImportState_Success_ByName(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil))
	req := resource.ImportStateRequest{ID: "Test Company"}
	var resp resource.ImportStateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: emptyRaw}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}

func TestManagedCompanyResource_ImportState_Success_ByID(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil, nil, nil))
	req := resource.ImportStateRequest{ID: "1169425105420462"}
	var resp resource.ImportStateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: emptyRaw}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}
