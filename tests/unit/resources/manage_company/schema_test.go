// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"testing"

	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/manage_company"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestManageCompanyResource_Schema(t *testing.T) {
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["node"] == nil {
		t.Error("expected node attribute")
	}
	if resp.Schema.Attributes["seats"] == nil {
		t.Error("expected seats attribute")
	}
	if resp.Schema.Attributes["plan"] == nil {
		t.Error("expected plan attribute")
	}
	if resp.Schema.Attributes["file_plan"] == nil {
		t.Error("expected file_plan attribute")
	}
	if resp.Schema.Attributes["add_ons"] == nil {
		t.Error("expected add_ons attribute")
	}
}
