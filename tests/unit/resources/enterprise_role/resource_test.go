// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseRoleResource_Metadata(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_role" {
		t.Errorf("expected TypeName commander_enterprise_role, got %s", resp.TypeName)
	}
}

func TestEnterpriseRoleResource_Configure_NilProviderData(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseRoleResource_Configure_InvalidProviderData(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseRoleResource_Configure_Success(t *testing.T) {
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseRoleResource(t *testing.T) {
	res := enterpriserole.NewEnterpriseRoleResource()
	if res == nil {
		t.Fatal("NewEnterpriseRoleResource returned nil")
	}
	_, ok := res.(*enterpriserole.EnterpriseRoleResource)
	if !ok {
		t.Errorf("expected *EnterpriseRoleResource, got %T", res)
	}
}
