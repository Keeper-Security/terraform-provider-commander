// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_database"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPamDatabaseDataSource_Metadata(t *testing.T) {
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_classic_pam_database" {
		t.Errorf("expected TypeName commander_classic_pam_database, got %s", resp.TypeName)
	}
}

func TestPamDatabaseDataSource_Configure_NilProviderData(t *testing.T) {
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamDatabaseDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamDatabaseDataSource_Configure_Success(t *testing.T) {
	d := pamdatabase.NewPamDatabaseDataSource().(*pamdatabase.PamDatabaseDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamDatabaseDataSource(t *testing.T) {
	ds := pamdatabase.NewPamDatabaseDataSource()
	if ds == nil {
		t.Fatal("NewPamDatabaseDataSource returned nil")
	}
	_, ok := ds.(*pamdatabase.PamDatabaseDataSource)
	if !ok {
		t.Errorf("expected *PamDatabaseDataSource, got %T", ds)
	}
}
