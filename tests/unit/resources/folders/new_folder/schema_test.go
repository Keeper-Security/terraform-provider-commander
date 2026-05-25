// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder_test

import (
	"context"
	"testing"

	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestNewFolderResource_Schema(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
}
