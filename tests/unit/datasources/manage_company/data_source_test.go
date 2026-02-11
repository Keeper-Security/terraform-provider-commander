// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/manage_company"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestManageCompanyDataSource_Metadata(t *testing.T) {
	d := managecompany.NewManageCompanyDataSource().(*managecompany.ManageCompanyDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_manage_company" {
		t.Errorf("expected TypeName commander_manage_company, got %s", resp.TypeName)
	}
}

func TestManageCompanyDataSource_Configure_NilProviderData(t *testing.T) {
	d := managecompany.NewManageCompanyDataSource().(*managecompany.ManageCompanyDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestManageCompanyDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := managecompany.NewManageCompanyDataSource().(*managecompany.ManageCompanyDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestManageCompanyDataSource_Configure_Success(t *testing.T) {
	d := managecompany.NewManageCompanyDataSource().(*managecompany.ManageCompanyDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewManageCompanyDataSource(t *testing.T) {
	ds := managecompany.NewManageCompanyDataSource()
	if ds == nil {
		t.Fatal("NewManageCompanyDataSource returned nil")
	}
	_, ok := ds.(*managecompany.ManageCompanyDataSource)
	if !ok {
		t.Errorf("expected *ManageCompanyDataSource, got %T", ds)
	}
}
