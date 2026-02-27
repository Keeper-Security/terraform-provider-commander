// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/managed_company"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestManagedCompanyResource_Metadata(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_managed_company" {
		t.Errorf("expected TypeName commander_managed_company, got %s", resp.TypeName)
	}
}

func TestManagedCompanyResource_Configure_NilProviderData(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestManagedCompanyResource_Configure_InvalidProviderData(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestManagedCompanyResource_Configure_Success(t *testing.T) {
	r := managedcompany.NewManagedCompanyResource().(*managedcompany.ManagedCompanyResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewManagedCompanyResource(t *testing.T) {
	res := managedcompany.NewManagedCompanyResource()
	if res == nil {
		t.Fatal("NewManagedCompanyResource returned nil")
	}
	_, ok := res.(*managedcompany.ManagedCompanyResource)
	if !ok {
		t.Errorf("expected *ManagedCompanyResource, got %T", res)
	}
}
