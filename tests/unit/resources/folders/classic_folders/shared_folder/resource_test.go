// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	classicsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSharedFolderResource_Metadata(t *testing.T) {
	r := classicsharedfolder.NewClassicSharedFolderResource().(*classicsharedfolder.ClassicSharedFolderResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_shared_folder" {
		t.Errorf("expected TypeName commander_shared_folder, got %s", resp.TypeName)
	}
}

func TestSharedFolderResource_Configure_NilProviderData(t *testing.T) {
	r := classicsharedfolder.NewClassicSharedFolderResource().(*classicsharedfolder.ClassicSharedFolderResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestSharedFolderResource_Configure_InvalidProviderData(t *testing.T) {
	r := classicsharedfolder.NewClassicSharedFolderResource().(*classicsharedfolder.ClassicSharedFolderResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestSharedFolderResource_Configure_Success(t *testing.T) {
	r := classicsharedfolder.NewClassicSharedFolderResource().(*classicsharedfolder.ClassicSharedFolderResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewClassicSharedFolderResource(t *testing.T) {
	res := classicsharedfolder.NewClassicSharedFolderResource()
	if res == nil {
		t.Fatal("NewClassicSharedFolderResource returned nil")
	}
	_, ok := res.(*classicsharedfolder.ClassicSharedFolderResource)
	if !ok {
		t.Errorf("expected *ClassicSharedFolderResource, got %T", res)
	}
}
