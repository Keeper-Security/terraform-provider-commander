// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_team"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseTeamDataSource_Metadata(t *testing.T) {
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_team" {
		t.Errorf("expected TypeName commander_enterprise_team, got %s", resp.TypeName)
	}
}

func TestEnterpriseTeamDataSource_Configure_NilProviderData(t *testing.T) {
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseTeamDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseTeamDataSource_Configure_Success(t *testing.T) {
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseTeamDataSource(t *testing.T) {
	ds := enterpriseteam.NewEnterpriseTeamDataSource()
	if ds == nil {
		t.Fatal("NewEnterpriseTeamDataSource returned nil")
	}
	_, ok := ds.(*enterpriseteam.EnterpriseTeamDataSource)
	if !ok {
		t.Errorf("expected *EnterpriseTeamDataSource, got %T", ds)
	}
}
