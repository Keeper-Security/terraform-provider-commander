// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim_test

import (
	"context"
	"testing"

	enterprisescim "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_scim"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseScimDataSource_Schema(t *testing.T) {
	d := enterprisescim.NewEnterpriseScimDataSource().(*enterprisescim.EnterpriseScimDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["scim"] == nil {
		t.Error("expected scim attribute")
	}
	if resp.Schema.Attributes["scim_id"] == nil {
		t.Error("expected scim_id attribute")
	}
	if resp.Schema.Attributes["scim_url"] == nil {
		t.Error("expected scim_url attribute")
	}
	if resp.Schema.Attributes["node_id"] == nil {
		t.Error("expected node_id attribute")
	}
	if resp.Schema.Attributes["node_name"] == nil {
		t.Error("expected node_name attribute")
	}
	if resp.Schema.Attributes["status"] == nil {
		t.Error("expected status attribute")
	}
	if resp.Schema.Attributes["prefix"] == nil {
		t.Error("expected prefix attribute")
	}
	if resp.Schema.Attributes["unique_groups"] == nil {
		t.Error("expected unique_groups attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
