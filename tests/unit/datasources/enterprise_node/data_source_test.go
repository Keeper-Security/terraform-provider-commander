// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_node"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseNodesDataSource_Metadata(t *testing.T) {
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	req := datasource.MetadataRequest{ProviderTypeName: "commander"}
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander_enterprise_node" {
		t.Errorf("expected TypeName commander_enterprise_node, got %s", resp.TypeName)
	}
}

func TestEnterpriseNodesDataSource_Configure_NilProviderData(t *testing.T) {
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	req := datasource.ConfigureRequest{ProviderData: nil}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestEnterpriseNodesDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	req := datasource.ConfigureRequest{ProviderData: "not-api-manager"}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when ProviderData is not *api.ApiManager")
	}
}

func TestEnterpriseNodesDataSource_Configure_Success(t *testing.T) {
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	am := &api.ApiManager{ServiceModeUrl: "http://test", ServiceModeApiKey: "key"}
	req := datasource.ConfigureRequest{ProviderData: am}
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestNewEnterpriseNodesDataSource(t *testing.T) {
	ds := enterprisenode.NewEnterpriseNodesDataSource()
	if ds == nil {
		t.Fatal("NewEnterpriseNodesDataSource returned nil")
	}
	_, ok := ds.(*enterprisenode.EnterpriseNodesDataSource)
	if !ok {
		t.Errorf("expected *EnterpriseNodesDataSource, got %T", ds)
	}
}
