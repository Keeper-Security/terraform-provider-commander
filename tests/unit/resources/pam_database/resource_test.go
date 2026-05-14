// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_records/pam_database"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPamDatabaseResource_Metadata(t *testing.T) {
	r := pamdatabase.NewPamDatabaseResource().(*pamdatabase.PamDatabaseResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_pam_database" {
		t.Errorf("expected TypeName commander_pam_database, got %s", resp.TypeName)
	}
}

func TestPamDatabaseResource_Configure_NilProviderData(t *testing.T) {
	r := pamdatabase.NewPamDatabaseResource().(*pamdatabase.PamDatabaseResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamDatabaseResource_Configure_InvalidProviderData(t *testing.T) {
	r := pamdatabase.NewPamDatabaseResource().(*pamdatabase.PamDatabaseResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamDatabaseResource_Configure_Success(t *testing.T) {
	r := pamdatabase.NewPamDatabaseResource().(*pamdatabase.PamDatabaseResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewPamDatabaseResource(t *testing.T) {
	res := pamdatabase.NewPamDatabaseResource()
	if res == nil {
		t.Fatal("NewPamDatabaseResource returned nil")
	}
	_, ok := res.(*pamdatabase.PamDatabaseResource)
	if !ok {
		t.Errorf("expected *PamDatabaseResource, got %T", res)
	}
}
