// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPamDirectoryDataSource_Metadata(t *testing.T) {
	d := pamdirectory.NewPamDirectoryDataSource().(*pamdirectory.PamDirectoryDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_classic_pam_directory" {
		t.Errorf("expected TypeName commander_classic_pam_directory, got %s", resp.TypeName)
	}
}

func TestPamDirectoryDataSource_Configure_NilProviderData(t *testing.T) {
	d := pamdirectory.NewPamDirectoryDataSource().(*pamdirectory.PamDirectoryDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamDirectoryDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := pamdirectory.NewPamDirectoryDataSource().(*pamdirectory.PamDirectoryDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamDirectoryDataSource_Configure_Success(t *testing.T) {
	d := pamdirectory.NewPamDirectoryDataSource().(*pamdirectory.PamDirectoryDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamDirectoryDataSource(t *testing.T) {
	ds := pamdirectory.NewPamDirectoryDataSource()
	if ds == nil {
		t.Fatal("NewPamDirectoryDataSource returned nil")
	}
	_, ok := ds.(*pamdirectory.PamDirectoryDataSource)
	if !ok {
		t.Errorf("expected *PamDirectoryDataSource, got %T", ds)
	}
}
