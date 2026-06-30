// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPamMachineResource_Metadata(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_classic_pam_machine" {
		t.Errorf("expected TypeName commander_classic_pam_machine, got %s", resp.TypeName)
	}
}

func TestPamMachineResource_Configure_NilProviderData(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamMachineResource_Configure_InvalidProviderData(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamMachineResource_Configure_Success(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamMachineResource(t *testing.T) {
	res := pammachine.NewPamMachineResource()
	if res == nil {
		t.Fatal("NewPamMachineResource returned nil")
	}
	_, ok := res.(*pammachine.PamMachineResource)
	if !ok {
		t.Errorf("expected *PamMachineResource, got %T", res)
	}
}
