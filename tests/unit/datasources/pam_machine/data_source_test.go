// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_records/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPamMachineDataSource_Metadata(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_pam_machine" {
		t.Errorf("expected TypeName commander_pam_machine, got %s", resp.TypeName)
	}
}

func TestPamMachineDataSource_Configure_NilProviderData(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamMachineDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamMachineDataSource_Configure_Success(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamMachineDataSource(t *testing.T) {
	ds := pammachine.NewPamMachineDataSource()
	if ds == nil {
		t.Fatal("NewPamMachineDataSource returned nil")
	}
	_, ok := ds.(*pammachine.PamMachineDataSource)
	if !ok {
		t.Errorf("expected *PamMachineDataSource, got %T", ds)
	}
}
