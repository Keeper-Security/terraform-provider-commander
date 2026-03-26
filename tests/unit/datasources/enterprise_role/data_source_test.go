// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_role"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseRoleDataSource_Metadata(t *testing.T) {
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_role" {
		t.Errorf("expected TypeName commander_enterprise_role, got %s", resp.TypeName)
	}
}

func TestEnterpriseRoleDataSource_Configure_NilProviderData(t *testing.T) {
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseRoleDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseRoleDataSource_Configure_Success(t *testing.T) {
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseRoleDataSource(t *testing.T) {
	ds := enterpriserole.NewEnterpriseRoleDataSource()
	if ds == nil {
		t.Fatal("NewEnterpriseRoleDataSource returned nil")
	}
	_, ok := ds.(*enterpriserole.EnterpriseRoleDataSource)
	if !ok {
		t.Errorf("expected *EnterpriseRoleDataSource, got %T", ds)
	}
}
