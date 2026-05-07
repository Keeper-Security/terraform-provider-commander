// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisescim "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_scim"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseScimDataSource_Metadata(t *testing.T) {
	d := enterprisescim.NewEnterpriseScimDataSource().(*enterprisescim.EnterpriseScimDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_scim" {
		t.Errorf("expected TypeName commander_enterprise_scim, got %s", resp.TypeName)
	}
}

func TestEnterpriseScimDataSource_Configure_NilProviderData(t *testing.T) {
	d := enterprisescim.NewEnterpriseScimDataSource().(*enterprisescim.EnterpriseScimDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseScimDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := enterprisescim.NewEnterpriseScimDataSource().(*enterprisescim.EnterpriseScimDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseScimDataSource_Configure_Success(t *testing.T) {
	d := enterprisescim.NewEnterpriseScimDataSource().(*enterprisescim.EnterpriseScimDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseScimDataSource(t *testing.T) {
	ds := enterprisescim.NewEnterpriseScimDataSource()
	if ds == nil {
		t.Fatal("NewEnterpriseScimDataSource returned nil")
	}
	_, ok := ds.(*enterprisescim.EnterpriseScimDataSource)
	if !ok {
		t.Errorf("expected *EnterpriseScimDataSource, got %T", ds)
	}
}
