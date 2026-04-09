// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisepush_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterprisePushResource_Metadata(t *testing.T) {
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_push" {
		t.Errorf("expected TypeName commander_enterprise_push, got %s", resp.TypeName)
	}
}

func TestEnterprisePushResource_Configure_NilProviderData(t *testing.T) {
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterprisePushResource_Configure_InvalidProviderData(t *testing.T) {
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterprisePushResource_Configure_Success(t *testing.T) {
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterprisePushResource(t *testing.T) {
	res := enterprisepush.NewEnterprisePushResource()
	if res == nil {
		t.Fatal("NewEnterprisePushResource returned nil")
	}
	_, ok := res.(*enterprisepush.EnterprisePushResource)
	if !ok {
		t.Errorf("expected *EnterprisePushResource, got %T", res)
	}
}
