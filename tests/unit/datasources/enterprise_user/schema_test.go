// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"testing"

	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_user"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnterpriseUserDataSource_Schema(t *testing.T) {
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["user"] == nil {
		t.Error("expected user attribute")
	}
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["email"] == nil {
		t.Error("expected email attribute")
	}
	if resp.Schema.Attributes["job_title"] == nil {
		t.Error("expected job_title attribute")
	}
	if resp.Schema.Attributes["roles"] == nil {
		t.Error("expected roles attribute")
	}
	if resp.Schema.Attributes["teams"] == nil {
		t.Error("expected teams attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
