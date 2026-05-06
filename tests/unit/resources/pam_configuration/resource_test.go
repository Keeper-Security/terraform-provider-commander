// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPamConfigurationResource_Metadata(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_pam_configuration" {
		t.Errorf("expected TypeName commander_pam_configuration, got %s", resp.TypeName)
	}
}

func TestPamConfigurationResource_Configure_NilProviderData(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamConfigurationResource_Configure_InvalidProviderData(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamConfigurationResource_Configure_Success(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamConfigurationResource(t *testing.T) {
	res := pamconfiguration.NewPamConfigurationResource()
	if res == nil {
		t.Fatal("NewPamConfigurationResource returned nil")
	}
	_, ok := res.(*pamconfiguration.PamConfigurationResource)
	if !ok {
		t.Errorf("expected *PamConfigurationResource, got %T", res)
	}
}
