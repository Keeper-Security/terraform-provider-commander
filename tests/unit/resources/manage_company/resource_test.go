// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/manage_company"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestManageCompanyResource_Metadata(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_manage_company" {
		t.Errorf("expected TypeName commander_manage_company, got %s", resp.TypeName)
	}
}

func TestManageCompanyResource_Configure_NilProviderData(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestManageCompanyResource_Configure_InvalidProviderData(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestManageCompanyResource_Configure_Success(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewManageCompanyResource(t *testing.T) {
	res := managecompany.NewManageCompanyResource()
	if res == nil {
		t.Fatal("NewManageCompanyResource returned nil")
	}
	_, ok := res.(*managecompany.ManageCompanyResource)
	if !ok {
		t.Errorf("expected *ManageCompanyResource, got %T", res)
	}
}
