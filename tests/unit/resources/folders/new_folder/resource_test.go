// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestNewFolderResource_Metadata(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_new_folder" {
		t.Errorf("expected TypeName commander_new_folder, got %s", resp.TypeName)
	}
}

func TestNewFolderResource_Configure_NilProviderData(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestNewFolderResource_Configure_InvalidProviderData(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestNewFolderResource_Configure_Success(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewNewFolderResource(t *testing.T) {
	res := newfolder.NewNewFolderResource()
	if res == nil {
		t.Fatal("NewNewFolderResource returned nil")
	}
	_, ok := res.(*newfolder.NewFolderResource)
	if !ok {
		t.Errorf("expected *NewFolderResource, got %T", res)
	}
}
