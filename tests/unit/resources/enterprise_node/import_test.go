// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseNodeResource_ImportState_EmptyID(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestEnterpriseNodeResource_ImportState_ManagedCompanyEmptyNode(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	req := resource.ImportStateRequest{ID: "Company,"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when node part is empty in managed company format")
	}
}

func TestEnterpriseNodeResource_ImportState_NoApiManager(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	req := resource.ImportStateRequest{ID: "Root"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is nil")
	}
}

func TestEnterpriseNodeResource_ImportState_Success_NodeOnly(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))
	req := resource.ImportStateRequest{ID: "Root"}
	var resp resource.ImportStateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: emptyRaw}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseNodeResource_ImportState_Success_WithManagedCompany(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))
	req := resource.ImportStateRequest{ID: "Test Company,Root"}
	var resp resource.ImportStateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: emptyRaw}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}
