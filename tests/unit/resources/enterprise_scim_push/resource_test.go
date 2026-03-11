// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseScimPushResource_Metadata(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_scim_push" {
		t.Errorf("expected TypeName commander_enterprise_scim_push, got %s", resp.TypeName)
	}
}

func TestEnterpriseScimPushResource_Configure_NilProviderData(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseScimPushResource_Configure_InvalidProviderData(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseScimPushResource_Configure_Success(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseScimPushResource(t *testing.T) {
	res := enterprisescimpush.NewEnterpriseScimPushResource()
	if res == nil {
		t.Fatal("NewEnterpriseScimPushResource returned nil")
	}
	_, ok := res.(*enterprisescimpush.EnterpriseScimPushResource)
	if !ok {
		t.Errorf("expected *EnterpriseScimPushResource, got %T", res)
	}
}
