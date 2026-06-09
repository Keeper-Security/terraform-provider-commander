// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_folder_test

import (
	"context"
	"testing"

	commonnewfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/new_folder"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
)

func TestBuildNewFolderGetCommand(t *testing.T) {
	got := commonnewfolder.BuildNewFolderGetCommand("Cuuc9aK6VuATH49ewBf0zg")
	want := `nsf-get "Cuuc9aK6VuATH49ewBf0zg" --format json`
	if got != want {
		t.Errorf("BuildNewFolderGetCommand() = %q, want %q", got, want)
	}
}

func TestMapResponseToModel(t *testing.T) {
	var m commonnewfolder.Model
	apiData := &commonnewfolder.NewFolderGetResponse{
		FolderUID: "Cuuc9aK6VuATH49ewBf0zg",
		Name:      "Engineering",
	}
	if err := commonnewfolder.MapResponseToModel(context.Background(), apiData, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	if m.Id.ValueString() != apiData.FolderUID {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), apiData.FolderUID)
	}
	if m.Name.ValueString() != apiData.Name {
		t.Errorf("name = %q, want %q", m.Name.ValueString(), apiData.Name)
	}
}

func TestMapResponseToModel_NilApi(t *testing.T) {
	var m commonnewfolder.Model
	if err := commonnewfolder.MapResponseToModel(context.Background(), nil, &m); err == nil {
		t.Error("expected error for nil API response")
	}
}

func TestNewFolderGetResponse_EmbedsShareFragment(t *testing.T) {
	// Compile-time check that NewFolderGetResponse can carry user_permissions
	// via the embedded new_share.ShareResponseFragment.
	apiData := &commonnewfolder.NewFolderGetResponse{
		FolderUID: "F1",
		Name:      "n",
		ShareResponseFragment: new_share.ShareResponseFragment{
			UserPermissions: []new_share.UserPermissionEntry{
				{Accessor: "a@x.com", Role: "viewer"},
				{Accessor: "owner@x.com", Role: "owner"},
			},
		},
	}
	if len(apiData.UserPermissions) != 2 {
		t.Errorf("expected 2 user_permissions, got %d", len(apiData.UserPermissions))
	}
}
