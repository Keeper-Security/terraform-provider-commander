// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"testing"

	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_team"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseTeamDataSource_Schema(t *testing.T) {
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["team"] == nil {
		t.Error("expected team attribute")
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
	if resp.Schema.Attributes["roles"] == nil {
		t.Error("expected roles attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
