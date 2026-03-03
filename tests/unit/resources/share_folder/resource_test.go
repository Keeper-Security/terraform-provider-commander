// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	sharefolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/share_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestShareFolderResource_Metadata(t *testing.T) {
	r := sharefolder.NewShareFolderResource().(*sharefolder.ShareFolderResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_share_folder" {
		t.Errorf("expected TypeName commander_share_folder, got %s", resp.TypeName)
	}
}

func TestShareFolderResource_Configure_NilProviderData(t *testing.T) {
	r := sharefolder.NewShareFolderResource().(*sharefolder.ShareFolderResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestShareFolderResource_Configure_InvalidProviderData(t *testing.T) {
	r := sharefolder.NewShareFolderResource().(*sharefolder.ShareFolderResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestShareFolderResource_Configure_Success(t *testing.T) {
	r := sharefolder.NewShareFolderResource().(*sharefolder.ShareFolderResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewShareFolderResource(t *testing.T) {
	res := sharefolder.NewShareFolderResource()
	if res == nil {
		t.Fatal("NewShareFolderResource returned nil")
	}
	_, ok := res.(*sharefolder.ShareFolderResource)
	if !ok {
		t.Errorf("expected *ShareFolderResource, got %T", res)
	}
}
