// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"testing"

	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseTeamResource_Schema(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["restrict_record_edit"] == nil {
		t.Error("expected restrict_record_edit attribute")
	}
	if resp.Schema.Attributes["restrict_record_re_share"] == nil {
		t.Error("expected restrict_record_re_share attribute")
	}
	if resp.Schema.Attributes["enable_privacy_screen"] == nil {
		t.Error("expected enable_privacy_screen attribute")
	}
	if resp.Schema.Attributes["users"] == nil {
		t.Error("expected users attribute")
	}
	if resp.Schema.Attributes["roles"] == nil {
		t.Error("expected roles attribute")
	}
	if resp.Schema.Attributes["node"] == nil {
		t.Error("expected node attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
