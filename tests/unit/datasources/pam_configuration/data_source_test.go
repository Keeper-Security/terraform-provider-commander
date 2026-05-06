// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPamConfigurationDataSource_Metadata(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_pam_configuration" {
		t.Errorf("expected TypeName commander_pam_configuration, got %s", resp.TypeName)
	}
}

func TestPamConfigurationDataSource_Configure_NilProviderData(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamConfigurationDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamConfigurationDataSource_Configure_Success(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamConfigurationDataSource(t *testing.T) {
	ds := pamconfiguration.NewPamConfigurationDataSource()
	if ds == nil {
		t.Fatal("NewPamConfigurationDataSource returned nil")
	}
	_, ok := ds.(*pamconfiguration.PamConfigurationDataSource)
	if !ok {
		t.Errorf("expected *PamConfigurationDataSource, got %T", ds)
	}
}
