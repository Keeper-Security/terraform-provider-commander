// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPamRemoteBrowserResource_Metadata(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_classic_pam_remote_browser" {
		t.Errorf("expected TypeName commander_classic_pam_remote_browser, got %s", resp.TypeName)
	}
}

func TestPamRemoteBrowserResource_Configure_NilProviderData(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamRemoteBrowserResource_Configure_InvalidProviderData(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamRemoteBrowserResource_Configure_Success(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamRemoteBrowserResource(t *testing.T) {
	res := pamremotebrowser.NewPamRemoteBrowserResource()
	if res == nil {
		t.Fatal("NewPamRemoteBrowserResource returned nil")
	}
	_, ok := res.(*pamremotebrowser.PamRemoteBrowserResource)
	if !ok {
		t.Errorf("expected *PamRemoteBrowserResource, got %T", res)
	}
}
