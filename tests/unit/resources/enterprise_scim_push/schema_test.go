// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush_test

import (
	"context"
	"testing"

	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterpriseScimPushResource_Schema(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["scim_id"] == nil {
		t.Error("expected scim_id attribute")
	}
	if resp.Schema.Attributes["source"] == nil {
		t.Error("expected source attribute")
	}
	if resp.Schema.Attributes["record"] == nil {
		t.Error("expected record attribute")
	}
	if resp.Schema.Attributes["auto_approve"] == nil {
		t.Error("expected auto_approve attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
