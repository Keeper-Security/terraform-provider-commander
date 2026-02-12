// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_user"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseUserDataSource_Metadata(t *testing.T) {
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_user" {
		t.Errorf("expected TypeName commander_enterprise_user, got %s", resp.TypeName)
	}
}

func TestEnterpriseUserDataSource_Configure_NilProviderData(t *testing.T) {
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseUserDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseUserDataSource_Configure_Success(t *testing.T) {
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseUserDataSource(t *testing.T) {
	ds := enterpriseuser.NewEnterpriseUserDataSource()
	if ds == nil {
		t.Fatal("NewEnterpriseUserDataSource returned nil")
	}
	_, ok := ds.(*enterpriseuser.EnterpriseUserDataSource)
	if !ok {
		t.Errorf("expected *EnterpriseUserDataSource, got %T", ds)
	}
}
