// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/manage_company"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestManageCompanyResource_ImportState_EmptyID(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestManageCompanyResource_ImportState_NoApiManager(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	req := resource.ImportStateRequest{ID: "Test Company"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is nil")
	}
}

func TestManageCompanyResource_ImportState_Success_ByName(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
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

func TestManageCompanyResource_ImportState_Success_ByID(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
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
