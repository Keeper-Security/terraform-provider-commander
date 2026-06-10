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

func TestCollectFolderSharePermissions_NilApiData(t *testing.T) {
	if got := commonnewfolder.CollectFolderSharePermissions(nil); got != nil {
		t.Errorf("expected nil for nil apiData, got %v", got)
	}
}

func TestCollectFolderSharePermissions_EmptyArrays(t *testing.T) {
	got := commonnewfolder.CollectFolderSharePermissions(&commonnewfolder.NewFolderGetResponse{})
	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d (%v)", len(got), got)
	}
}

func TestCollectFolderSharePermissions_UsersOnly(t *testing.T) {
	apiData := &commonnewfolder.NewFolderGetResponse{
		ShareResponseFragment: new_share.ShareResponseFragment{
			UserPermissions: []new_share.UserPermissionEntry{
				{Accessor: "a@x.com", AccessorType: "AT_USER", Role: "viewer"},
				{Accessor: "b@x.com", AccessorType: "AT_USER", Role: "full-manager"},
			},
		},
	}
	got := commonnewfolder.CollectFolderSharePermissions(apiData)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d (%v)", len(got), got)
	}
	if got[0].Accessor != "a@x.com" || got[1].Accessor != "b@x.com" {
		t.Errorf("unexpected order/contents: %v", got)
	}
}

func TestCollectFolderSharePermissions_TeamsOnly(t *testing.T) {
	apiData := &commonnewfolder.NewFolderGetResponse{
		ShareResponseFragment: new_share.ShareResponseFragment{
			TeamPermissions: []new_share.UserPermissionEntry{
				{Accessor: "Metron", AccessorType: "AT_TEAM", Role: "viewer"},
			},
		},
	}
	got := commonnewfolder.CollectFolderSharePermissions(apiData)
	if len(got) != 1 || got[0].Accessor != "Metron" {
		t.Errorf("expected [Metron], got %v", got)
	}
}

func TestCollectFolderSharePermissions_MixedFiltersApplication(t *testing.T) {
	apiData := &commonnewfolder.NewFolderGetResponse{
		ShareResponseFragment: new_share.ShareResponseFragment{
			UserPermissions: []new_share.UserPermissionEntry{
				{Accessor: "a@x.com", AccessorType: "AT_USER", Role: "viewer"},
				{Accessor: "app-1", AccessorType: commonnewfolder.AccessorTypeApplication, Role: "viewer"},
				{Accessor: "b@x.com", AccessorType: "AT_USER", Role: "share-manager"},
			},
			TeamPermissions: []new_share.UserPermissionEntry{
				{Accessor: "Metron", AccessorType: "AT_TEAM", Role: "viewer"},
				{Accessor: "app-team", AccessorType: commonnewfolder.AccessorTypeApplication, Role: "viewer"},
			},
		},
	}
	got := commonnewfolder.CollectFolderSharePermissions(apiData)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries after AT_APPLICATION filter, got %d (%v)", len(got), got)
	}
	accessors := map[string]bool{}
	for _, e := range got {
		accessors[e.Accessor] = true
		if e.AccessorType == commonnewfolder.AccessorTypeApplication {
			t.Errorf("AT_APPLICATION entry %q leaked through filter", e.Accessor)
		}
	}
	for _, want := range []string{"a@x.com", "b@x.com", "Metron"} {
		if !accessors[want] {
			t.Errorf("expected %q in result, got %v", want, got)
		}
	}
}

func TestCollectFolderSharePermissions_CaseSensitiveFilter(t *testing.T) {
	apiData := &commonnewfolder.NewFolderGetResponse{
		ShareResponseFragment: new_share.ShareResponseFragment{
			UserPermissions: []new_share.UserPermissionEntry{
				{Accessor: "lower", AccessorType: "at_application", Role: "viewer"},
				{Accessor: "mixed", AccessorType: "At_Application", Role: "viewer"},
				{Accessor: "exact", AccessorType: commonnewfolder.AccessorTypeApplication, Role: "viewer"},
			},
		},
	}
	got := commonnewfolder.CollectFolderSharePermissions(apiData)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries (strict equality keeps non-canonical case), got %d (%v)", len(got), got)
	}
	for _, e := range got {
		if e.Accessor == "exact" {
			t.Errorf("strict-equal AT_APPLICATION entry %q should have been filtered", e.Accessor)
		}
	}
}

func TestCollectFolderSharePermissions_UsersBeforeTeams(t *testing.T) {
	apiData := &commonnewfolder.NewFolderGetResponse{
		ShareResponseFragment: new_share.ShareResponseFragment{
			UserPermissions: []new_share.UserPermissionEntry{
				{Accessor: "u1", AccessorType: "AT_USER", Role: "viewer"},
			},
			TeamPermissions: []new_share.UserPermissionEntry{
				{Accessor: "t1", AccessorType: "AT_TEAM", Role: "viewer"},
			},
		},
	}
	got := commonnewfolder.CollectFolderSharePermissions(apiData)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Accessor != "u1" || got[1].Accessor != "t1" {
		t.Errorf("expected users before teams; got %+v", got)
	}
}
