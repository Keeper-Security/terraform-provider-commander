// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"testing"

	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_role"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseRoleDataSource_Schema(t *testing.T) {
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["role"] == nil {
		t.Error("expected role attribute")
	}
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["users"] == nil {
		t.Error("expected users attribute")
	}
	if resp.Schema.Attributes["teams"] == nil {
		t.Error("expected teams attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
