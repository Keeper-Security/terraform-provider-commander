// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"context"
	"testing"

	sharedfolderds "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestSharedFolderDataSource_Schema(t *testing.T) {
	d := sharedfolderds.NewSharedFolderDataSource().(*sharedfolderds.SharedFolderDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)
	for _, attr := range []string{"shared_folder", "id", "name", "user_permissions", "record_permissions", "records", "users"} {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected %s attribute", attr)
		}
	}
}
