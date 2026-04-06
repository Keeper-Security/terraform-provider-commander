// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func recordObject(canShare, canEdit bool) types.Object {
	return types.ObjectValueMust(
		map[string]attr.Type{AttrCanShare: types.BoolType, AttrCanEdit: types.BoolType},
		map[string]attr.Value{
			AttrCanShare: types.BoolValue(canShare),
			AttrCanEdit:  types.BoolValue(canEdit),
		},
	)
}

func recordsMapWithKeys(t *testing.T, keysToPerm map[string]struct{ share, edit bool }) types.Map {
	t.Helper()
	elems := make(map[string]attr.Value, len(keysToPerm))
	for k, p := range keysToPerm {
		elems[k] = recordObject(p.share, p.edit)
	}
	m, diags := types.MapValue(RecordEntryMapElemType, elems)
	if diags.HasError() {
		t.Fatalf("MapValue: %v", diags)
	}
	return m
}

func TestRecordEntryMapKey_PrefersMatchingPriorUID(t *testing.T) {
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"uid-1": {true, false},
	})
	rec := utils.SharedFolderRecordEntry{RecordUID: "uid-1", RecordName: "My Record"}
	if got := recordEntryMapKey(rec, prior); got != "uid-1" {
		t.Fatalf("got %q want uid-1", got)
	}
}

func TestRecordEntryMapKey_PrefersMatchingPriorRecordName(t *testing.T) {
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"Custom Title": {false, true},
	})
	rec := utils.SharedFolderRecordEntry{RecordUID: "uid-99", RecordName: "Custom Title"}
	if got := recordEntryMapKey(rec, prior); got != "Custom Title" {
		t.Fatalf("got %q want Custom Title", got)
	}
}

func TestRecordEntryMapKey_NoPriorMatchUsesRecordName(t *testing.T) {
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{"other": {}})
	rec := utils.SharedFolderRecordEntry{RecordUID: "uid-1", RecordName: "Title A"}
	if got := recordEntryMapKey(rec, prior); got != "Title A" {
		t.Fatalf("got %q want Title A", got)
	}
}

func TestRecordEntryMapKey_NullPriorUsesRecordNameThenUID(t *testing.T) {
	rec := utils.SharedFolderRecordEntry{RecordUID: "uid-1", RecordName: "Title A"}
	if got := recordEntryMapKey(rec, types.MapNull(RecordEntryMapElemType)); got != "Title A" {
		t.Fatalf("got %q want Title A", got)
	}
	rec2 := utils.SharedFolderRecordEntry{RecordUID: "uid-only", RecordName: ""}
	if got := recordEntryMapKey(rec2, types.MapNull(RecordEntryMapElemType)); got != "uid-only" {
		t.Fatalf("got %q want uid-only", got)
	}
}

func TestBuildRecordsMapFromAPIResponse_NameOnlyEntryIncluded(t *testing.T) {
	entries := []utils.SharedFolderRecordEntry{
		{RecordUID: "", RecordName: "name-only", CanShare: true, CanEdit: false},
	}
	m, err := buildRecordsMapFromAPIResponse(entries, types.MapNull(RecordEntryMapElemType))
	if err != nil {
		t.Fatal(err)
	}
	if m.IsNull() {
		t.Fatal("expected non-null map")
	}
	els := m.Elements()
	if len(els) != 1 {
		t.Fatalf("len=%d want 1", len(els))
	}
	if _, ok := els["name-only"]; !ok {
		t.Fatalf("missing key name-only, keys=%v", els)
	}
}

func TestBuildRecordsMapFromAPIResponse_ReusesPriorKeyByUID(t *testing.T) {
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"legacy-uid-key": {false, false},
	})
	entries := []utils.SharedFolderRecordEntry{
		{RecordUID: "legacy-uid-key", RecordName: "New Title", CanShare: true, CanEdit: true},
	}
	m, err := buildRecordsMapFromAPIResponse(entries, prior)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.Elements()["legacy-uid-key"]; !ok {
		t.Fatalf("expected prior uid key preserved, got keys %v", m.Elements())
	}
}

func TestMapSharedFolderApiResponseToModel_RecordsMerge(t *testing.T) {
	api := &utils.SharedFolderResponse{
		SharedFolderUID: "sf-1",
		Path:            "/folder",
		ManageUsers:     false,
		ManageRecords:   false,
		CanShare:        false,
		CanEdit:         false,
		Records: []utils.SharedFolderRecordEntry{
			{RecordUID: "u1", RecordName: "n1", CanShare: true, CanEdit: false},
		},
		Users: nil,
	}
	state := &SharedFolderResourceModel{
		Records: recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
			"u1": {false, false},
		}),
	}
	if err := MapSharedFolderApiResponseToModel(api, state); err != nil {
		t.Fatal(err)
	}
	if state.Records.IsNull() {
		t.Fatal("records null")
	}
	if _, ok := state.Records.Elements()["u1"]; !ok {
		t.Fatalf("expected key u1, got %v", state.Records.Elements())
	}
}
