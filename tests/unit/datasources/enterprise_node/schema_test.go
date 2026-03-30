// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"testing"

	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_node"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseNodesDataSource_Schema(t *testing.T) {
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["node"] == nil {
		t.Error("expected node attribute")
	}
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["parent"] == nil {
		t.Error("expected parent attribute")
	}
	if resp.Schema.Attributes["parent_id"] == nil {
		t.Error("expected parent_id attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
