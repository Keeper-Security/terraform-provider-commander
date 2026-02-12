// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseTeamResource_Metadata(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	req := resource.MetadataRequest{ProviderTypeName: "commander"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_team" {
		t.Errorf("expected TypeName commander_enterprise_team, got %s", resp.TypeName)
	}
}

func TestEnterpriseTeamResource_Configure_NilProviderData(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	req := resource.ConfigureRequest{ProviderData: nil}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseTeamResource_Configure_InvalidProviderData(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	req := resource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseTeamResource_Configure_Success(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := resource.ConfigureRequest{ProviderData: am}
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseTeamResource(t *testing.T) {
	res := enterpriseteam.NewEnterpriseTeamResource()
	if res == nil {
		t.Fatal("NewEnterpriseTeamResource returned nil")
	}
	_, ok := res.(*enterpriseteam.EnterpriseTeamResource)
	if !ok {
		t.Errorf("expected *EnterpriseTeamResource, got %T", res)
	}
}
