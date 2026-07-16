// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolderds_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newfolderds "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/new_folder"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestNewFolderDataSource_Metadata(t *testing.T) {
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_new_folder" {
		t.Errorf("expected TypeName commander_new_folder, got %s", resp.TypeName)
	}
}

func TestNewFolderDataSource_Configure_SetsApiManager(t *testing.T) {
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://example", ServiceModeApiKey: "k"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure produced diagnostics: %v", resp.Diagnostics)
	}
	if err := d.EnsureApiManager(); err != nil {
		t.Errorf("EnsureApiManager after Configure: %v", err)
	}
}

func TestNewFolderDataSource_Configure_NilProviderData_NoOp(t *testing.T) {
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("Configure with nil ProviderData should be a no-op; got diagnostics: %v", resp.Diagnostics)
	}
	if err := d.EnsureApiManager(); err == nil {
		t.Error("EnsureApiManager should fail when ApiManager was never set")
	}
}

func TestNewFolderDataSource_Configure_WrongProviderData_Errors(t *testing.T) {
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "not-an-api-manager"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected Configure to fail when ProviderData is the wrong type")
	}
}
