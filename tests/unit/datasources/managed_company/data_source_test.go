// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/managed_company"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestManagedCompanyDataSource_Metadata(t *testing.T) {
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_managed_company" {
		t.Errorf("expected TypeName commander_managed_company, got %s", resp.TypeName)
	}
}

func TestManagedCompanyDataSource_Configure_NilProviderData(t *testing.T) {
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestManagedCompanyDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestManagedCompanyDataSource_Configure_Success(t *testing.T) {
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewManagedCompanyDataSource(t *testing.T) {
	ds := managedcompany.NewManagedCompanyDataSource()
	if ds == nil {
		t.Fatal("NewManagedCompanyDataSource returned nil")
	}
	_, ok := ds.(*managedcompany.ManagedCompanyDataSource)
	if !ok {
		t.Errorf("expected *ManagedCompanyDataSource, got %T", ds)
	}
}
