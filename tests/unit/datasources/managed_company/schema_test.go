// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"context"
	"testing"

	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/managed_company"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestManagedCompanyDataSource_Schema(t *testing.T) {
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["node"] == nil {
		t.Error("expected node attribute")
	}
	if resp.Schema.Attributes["plan"] == nil {
		t.Error("expected plan attribute")
	}
	if resp.Schema.Attributes["file_plan"] == nil {
		t.Error("expected file_plan attribute")
	}
}
