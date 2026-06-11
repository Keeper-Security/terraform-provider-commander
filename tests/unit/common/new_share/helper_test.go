// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share_test

import (
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildGrantCommand_Folder(t *testing.T) {
	got := new_share.BuildGrantCommand(new_share.CmdNsfShareFolder, "FOLDER_UID_1", "user@example.com", "viewer")
	want := `nsf-share-folder "FOLDER_UID_1" --email='user@example.com' --action=grant --role='viewer'`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildGrantCommand_Record(t *testing.T) {
	got := new_share.BuildGrantCommand(new_share.CmdNsfShareRecord, "REC_UID_1", "alice@example.com", "full-manager")
	want := `nsf-share-record "REC_UID_1" --email='alice@example.com' --action=grant --role='full-manager'`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildGrantCommand_EscapesEmailWithApostrophe(t *testing.T) {
	got := new_share.BuildGrantCommand(new_share.CmdNsfShareFolder, "F1", "o'brien@example.com", "viewer")
	want := `nsf-share-folder "F1" --email='o''brien@example.com' --action=grant --role='viewer'`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildRevokeCommand_FolderUsesRemove(t *testing.T) {
	got := new_share.BuildRevokeCommand(new_share.CmdNsfShareFolder, "FOLDER_UID_1", "user@example.com")
	want := `nsf-share-folder "FOLDER_UID_1" --email='user@example.com' --action=remove`
	if got != want {
		t.Errorf("BuildRevokeCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildRevokeCommand_RecordUsesRevoke(t *testing.T) {
	got := new_share.BuildRevokeCommand(new_share.CmdNsfShareRecord, "REC_UID_1", "user@example.com")
	want := `nsf-share-record "REC_UID_1" --email='user@example.com' --action=revoke`
	if got != want {
		t.Errorf("BuildRevokeCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestMapResponseToModel_DropsOwnerAndEmpty(t *testing.T) {
	perms := []new_share.UserPermissionEntry{
		{Accessor: "owner@example.com", Role: "owner"},
		{Accessor: "viewer1@example.com", Role: "viewer"},
		{Accessor: "", Role: "viewer"},
		{Accessor: "no-role@example.com", Role: ""},
		{Accessor: "manager@example.com", Role: "full-manager"},
		{Accessor: "uppercase-owner@example.com", Role: "OWNER"},
	}
	var m new_share.ShareModel
	if err := new_share.MapResponseToModel(perms, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	if m.Share.IsNull() || m.Share.IsUnknown() {
		t.Fatalf("expected non-null share map, got null/unknown")
	}
	got := map[string]string{}
	for k, v := range m.Share.Elements() {
		s, ok := v.(types.String)
		if !ok {
			t.Fatalf("expected types.String, got %T", v)
		}
		got[k] = s.ValueString()
	}
	want := map[string]string{
		"viewer1@example.com": "viewer",
		"manager@example.com": "full-manager",
	}
	if len(got) != len(want) {
		t.Errorf("got %d entries, want %d (got=%v want=%v)", len(got), len(want), got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("share[%q] = %q, want %q", k, got[k], v)
		}
	}
}

func TestMapResponseToModel_EmptyResponseProducesNullMap(t *testing.T) {
	var m new_share.ShareModel
	if err := new_share.MapResponseToModel(nil, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	if !m.Share.IsNull() {
		t.Errorf("expected null share map (schema rejects {}), got %v", m.Share)
	}
}

func TestMapResponseToModel_OnlyOwnerProducesNullMap(t *testing.T) {
	perms := []new_share.UserPermissionEntry{
		{Accessor: "owner@example.com", Role: new_share.RoleOwner},
	}
	var m new_share.ShareModel
	if err := new_share.MapResponseToModel(perms, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	if !m.Share.IsNull() {
		t.Errorf("expected null share map when only the owner row is present, got %v", m.Share)
	}
}

func TestMapResponseToModel_NilModel(t *testing.T) {
	if err := new_share.MapResponseToModel(nil, nil); err == nil {
		t.Error("expected error for nil model")
	}
}

func TestAllowedRoles_ExcludesOwner(t *testing.T) {
	for _, r := range new_share.AllowedRoles {
		if r == new_share.RoleOwner {
			t.Errorf("AllowedRoles should not contain RoleOwner (got %q)", r)
		}
	}
	want := []string{
		new_share.RoleViewer,
		new_share.RoleShareManager,
		new_share.RoleContentManager,
		new_share.RoleContentShareManager,
		new_share.RoleFullManager,
	}
	if len(new_share.AllowedRoles) != len(want) {
		t.Errorf("AllowedRoles length = %d, want %d", len(new_share.AllowedRoles), len(want))
	}
	have := map[string]bool{}
	for _, r := range new_share.AllowedRoles {
		have[r] = true
	}
	for _, r := range want {
		if !have[r] {
			t.Errorf("AllowedRoles missing %q", r)
		}
	}
}
