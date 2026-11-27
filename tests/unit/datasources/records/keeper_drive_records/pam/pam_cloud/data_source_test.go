// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamcloud_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_cloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPamCloudDataSource_Metadata(t *testing.T) {
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "commander"}, &resp)
	if resp.TypeName != "commander_new_pam_cloud" {
		t.Errorf("expected TypeName commander_new_pam_cloud, got %s", resp.TypeName)
	}
}

func TestPamCloudDataSource_Configure_NilProviderData(t *testing.T) {
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestPamCloudDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "not-api-manager"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestPamCloudDataSource_Configure_Success(t *testing.T) {
	d := newpamcloud.NewPamCloudDataSource().(*newpamcloud.PamCloudDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}
