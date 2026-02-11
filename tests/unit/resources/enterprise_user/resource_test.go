// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_user"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseUserResource_Metadata(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_user" {
		t.Errorf("expected TypeName commander_enterprise_user, got %s", resp.TypeName)
	}
}

func TestEnterpriseUserResource_Configure_NilProviderData(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseUserResource_Configure_InvalidProviderData(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseUserResource_Configure_Success(t *testing.T) {
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseUserResource(t *testing.T) {
	res := enterpriseuser.NewEnterpriseUserResource()
	if res == nil {
		t.Fatal("NewEnterpriseUserResource returned nil")
	}
	_, ok := res.(*enterpriseuser.EnterpriseUserResource)
	if !ok {
		t.Errorf("expected *EnterpriseUserResource, got %T", res)
	}
}
