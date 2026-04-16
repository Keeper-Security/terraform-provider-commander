// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_records/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPamRemoteBrowserDataSource_Metadata(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_pam_remote_browser" {
		t.Errorf("expected TypeName commander_pam_remote_browser, got %s", resp.TypeName)
	}
}

func TestPamRemoteBrowserDataSource_Configure_NilProviderData(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamRemoteBrowserDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamRemoteBrowserDataSource_Configure_Success(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamRemoteBrowserDataSource(t *testing.T) {
	ds := pamremotebrowser.NewPamRemoteBrowserDataSource()
	if ds == nil {
		t.Fatal("NewPamRemoteBrowserDataSource returned nil")
	}
	_, ok := ds.(*pamremotebrowser.PamRemoteBrowserDataSource)
	if !ok {
		t.Errorf("expected *PamRemoteBrowserDataSource, got %T", ds)
	}
}
