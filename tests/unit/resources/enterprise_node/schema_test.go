// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"testing"

	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestEnterpriseNodeResource_Schema(t *testing.T) {
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["parent"] == nil {
		t.Error("expected parent attribute")
	}
	toggleAttr := resp.Schema.Attributes["toggle_isolated"]
	if toggleAttr == nil {
		t.Error("expected toggle_isolated attribute")
	}
	if b, ok := toggleAttr.(schema.BoolAttribute); ok && b.Default == nil {
		t.Error("expected toggle_isolated to have a default (false)")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
