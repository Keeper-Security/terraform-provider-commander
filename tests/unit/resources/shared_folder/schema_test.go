// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"context"
	"testing"

	sharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSharedFolderResource_Schema(t *testing.T) {
	r := sharedfolder.NewSharedFolderResource().(*sharedfolder.SharedFolderResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["folder_location"] == nil {
		t.Error("expected folder_location attribute")
	}
	if resp.Schema.Attributes["user_permissions"] == nil {
		t.Error("expected user_permissions attribute")
	}
	if resp.Schema.Attributes["record_permissions"] == nil {
		t.Error("expected record_permissions attribute")
	}
	if resp.Schema.Attributes["records"] == nil {
		t.Error("expected records attribute")
	}
	if resp.Schema.Attributes["users"] == nil {
		t.Error("expected users attribute")
	}
}
