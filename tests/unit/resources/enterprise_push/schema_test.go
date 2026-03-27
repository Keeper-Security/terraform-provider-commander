// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush_test

import (
	"context"
	"testing"

	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnterprisePushResource_Schema(t *testing.T) {
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["file_path"] == nil {
		t.Error("expected file_path attribute")
	}
	if resp.Schema.Attributes["file_content_sha256"] == nil {
		t.Error("expected file_content_sha256 attribute")
	}
	if resp.Schema.Attributes["email"] == nil {
		t.Error("expected email attribute")
	}
	if resp.Schema.Attributes["team"] == nil {
		t.Error("expected team attribute")
	}
	if resp.Schema.Attributes["managed_company"] == nil {
		t.Error("expected managed_company attribute")
	}
}
