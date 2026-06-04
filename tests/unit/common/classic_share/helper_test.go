// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share_test

import (
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildGrantCommand_GrantsBothFlags(t *testing.T) {
	got := classic_share.BuildGrantCommand("REC_UID_1", "user@example.com", true, true)
	want := `share-record --email 'user@example.com' 'REC_UID_1' --share --write`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildGrantCommand_GrantsShareOnly(t *testing.T) {
	got := classic_share.BuildGrantCommand("REC_UID_1", "user@example.com", true, false)
	want := `share-record --email 'user@example.com' 'REC_UID_1' --share`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildGrantCommand_GrantsEditOnly(t *testing.T) {
	got := classic_share.BuildGrantCommand("REC_UID_1", "user@example.com", false, true)
	want := `share-record --email 'user@example.com' 'REC_UID_1' --write`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildGrantCommand_DowngradeToViewerStripsShareAndWrite(t *testing.T) {
	got := classic_share.BuildGrantCommand("REC_UID_1", "user@example.com", false, false)
	want := `share-record --email 'user@example.com' 'REC_UID_1' --action revoke --share --write`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildGrantCommand_EscapesEmailWithApostrophe(t *testing.T) {
	got := classic_share.BuildGrantCommand("REC", "o'brien@example.com", true, false)
	want := `share-record --email 'o''brien@example.com' 'REC' --share`
	if got != want {
		t.Errorf("BuildGrantCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestBuildRevokeCommand(t *testing.T) {
	got := classic_share.BuildRevokeCommand("REC_UID_1", "user@example.com")
	want := `share-record --email 'user@example.com' 'REC_UID_1' --action revoke`
	if got != want {
		t.Errorf("BuildRevokeCommand =\n  %q\nwant\n  %q", got, want)
	}
}

func TestMapResponseToModel_PopulatesShare(t *testing.T) {
	perms := []classic_share.UserPermissionEntry{
		{Username: "alice@example.com", Shareable: true, Editable: false},
		{Username: "bob@example.com", Shareable: false, Editable: true},
		{Username: "carol@example.com", Shareable: true, Editable: true},
	}
	var m classic_share.ShareModel
	if err := classic_share.MapResponseToModel(perms, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	if m.Share.IsNull() || m.Share.IsUnknown() {
		t.Fatalf("expected non-null share map, got null/unknown")
	}

	want := map[string]struct {
		canShare, canEdit bool
	}{
		"alice@example.com": {true, false},
		"bob@example.com":   {false, true},
		"carol@example.com": {true, true},
	}
	got := unwrapShareMap(t, m.Share)
	if len(got) != len(want) {
		t.Errorf("got %d entries, want %d (got=%v want=%v)", len(got), len(want), got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("share[%q] = %+v, want %+v", k, got[k], v)
		}
	}
}

func TestMapResponseToModel_DropsEmptyUsername(t *testing.T) {
	perms := []classic_share.UserPermissionEntry{
		{Username: "", Shareable: true, Editable: true},
		{Username: "   ", Shareable: false, Editable: true},
		{Username: "alice@example.com", Shareable: true, Editable: false},
	}
	var m classic_share.ShareModel
	if err := classic_share.MapResponseToModel(perms, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	got := unwrapShareMap(t, m.Share)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d (%v)", len(got), got)
	}
	if _, ok := got["alice@example.com"]; !ok {
		t.Errorf("expected only alice in map, got %v", got)
	}
}

func TestMapResponseToModel_EmptyResponseProducesEmptyMap(t *testing.T) {
	var m classic_share.ShareModel
	if err := classic_share.MapResponseToModel(nil, &m); err != nil {
		t.Fatalf("MapResponseToModel: %v", err)
	}
	if m.Share.IsNull() {
		t.Error("expected non-null empty share map, got null")
	}
	if len(m.Share.Elements()) != 0 {
		t.Errorf("expected 0 elements, got %d", len(m.Share.Elements()))
	}
}

func TestMapResponseToModel_NilModel(t *testing.T) {
	if err := classic_share.MapResponseToModel(nil, nil); err == nil {
		t.Error("expected error for nil model")
	}
}

func TestSharePermissionsObjectType_ShapeStable(t *testing.T) {
	got := classic_share.SharePermissionsObjectType()
	if len(got) != 2 {
		t.Errorf("expected 2 attribute types, got %d", len(got))
	}
	if !got[classic_share.AttrCanShare].Equal(types.BoolType) {
		t.Errorf("expected %q to be types.BoolType, got %T", classic_share.AttrCanShare, got[classic_share.AttrCanShare])
	}
	if !got[classic_share.AttrCanEdit].Equal(types.BoolType) {
		t.Errorf("expected %q to be types.BoolType, got %T", classic_share.AttrCanEdit, got[classic_share.AttrCanEdit])
	}
}

// unwrapShareMap extracts the {can_share, can_edit} object map back into a
// plain Go map keyed by email.
func unwrapShareMap(t *testing.T, m types.Map) map[string]struct{ canShare, canEdit bool } {
	t.Helper()
	got := map[string]struct{ canShare, canEdit bool }{}
	for k, v := range m.Elements() {
		obj, ok := v.(types.Object)
		if !ok {
			t.Fatalf("expected types.Object for key %q, got %T", k, v)
		}
		attrs := obj.Attributes()
		got[k] = struct{ canShare, canEdit bool }{
			canShare: attrBoolValue(t, attrs, classic_share.AttrCanShare),
			canEdit:  attrBoolValue(t, attrs, classic_share.AttrCanEdit),
		}
	}
	return got
}

func attrBoolValue(t *testing.T, attrs map[string]attr.Value, key string) bool {
	t.Helper()
	v, ok := attrs[key].(types.Bool)
	if !ok {
		t.Fatalf("expected types.Bool at %q, got %T", key, attrs[key])
	}
	return v.ValueBool()
}
