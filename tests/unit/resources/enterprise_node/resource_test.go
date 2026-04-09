// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseNodeResource_Metadata(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_node" {
		t.Errorf("expected TypeName commander_enterprise_node, got %s", resp.TypeName)
	}
}

func TestEnterpriseNodeResource_Configure_NilProviderData(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseNodeResource_Configure_InvalidProviderData(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseNodeResource_Configure_Success(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseNodeResource(t *testing.T) {
	res := enterprisenode.NewEnterpriseNodeResource()
	if res == nil {
		t.Fatal("NewEnterpriseNodeResource returned nil")
	}
	_, ok := res.(*enterprisenode.EnterpriseNodeResource)
	if !ok {
		t.Errorf("expected *EnterpriseNodeResource, got %T", res)
	}
}
