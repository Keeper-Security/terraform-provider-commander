// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPamDirectoryResource_Metadata(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_classic_pam_directory" {
		t.Errorf("expected TypeName commander_classic_pam_directory, got %s", resp.TypeName)
	}
}

func TestPamDirectoryResource_Configure_NilProviderData(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamDirectoryResource_Configure_InvalidProviderData(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamDirectoryResource_Configure_Success(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamDirectoryResource(t *testing.T) {
	res := pamdirectory.NewPamDirectoryResource()
	if res == nil {
		t.Fatal("NewPamDirectoryResource returned nil")
	}
	_, ok := res.(*pamdirectory.PamDirectoryResource)
	if !ok {
		t.Errorf("expected *PamDirectoryResource, got %T", res)
	}
}
