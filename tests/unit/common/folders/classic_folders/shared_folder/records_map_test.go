// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"testing"

	commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func recordObject(canShare, canEdit bool) types.Object {
	return types.ObjectValueMust(
		map[string]attr.Type{commonsharedfolder.AttrCanShare: types.BoolType, commonsharedfolder.AttrCanEdit: types.BoolType},
		map[string]attr.Value{
			commonsharedfolder.AttrCanShare: types.BoolValue(canShare),
			commonsharedfolder.AttrCanEdit:  types.BoolValue(canEdit),
		},
	)
}

func recordsMapWithKeys(t *testing.T, keysToPerm map[string]struct{ share, edit bool }) types.Map {
	t.Helper()
	elems := make(map[string]attr.Value, len(keysToPerm))
	for k, p := range keysToPerm {
		elems[k] = recordObject(p.share, p.edit)
	}
	m, diags := types.MapValue(commonsharedfolder.RecordEntryMapElemType, elems)
	if diags.HasError() {
		t.Fatalf("MapValue: %v", diags)
	}
	return m
}

func nullPriorUsers() types.Map {
	return types.MapNull(commonsharedfolder.UserEntryMapElemType)
}

func mapRecordsViaAPI(t *testing.T, api *commonsharedfolder.SharedFolderResponse, priorRecords types.Map) *commonsharedfolder.Model {
	t.Helper()
	state := &commonsharedfolder.Model{}
	if err := commonsharedfolder.MapResponseToModel(api, state, nullPriorUsers(), priorRecords); err != nil {
		t.Fatal(err)
	}
	return state
}

func TestMapResponseToModel_Records_PrefersMatchingPriorUID(t *testing.T) {
	t.Parallel()
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"uid-1": {true, false},
	})
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "uid-1", RecordName: "My Record", CanShare: false, CanEdit: false},
		},
	}
	state := mapRecordsViaAPI(t, api, prior)
	if _, ok := state.Records.Elements()["uid-1"]; !ok {
		t.Fatalf("want key uid-1, got %v", state.Records.Elements())
	}
}

func TestMapResponseToModel_Records_PrefersMatchingPriorRecordName(t *testing.T) {
	t.Parallel()
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"Custom Title": {false, true},
	})
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "uid-99", RecordName: "Custom Title", CanShare: false, CanEdit: false},
		},
	}
	state := mapRecordsViaAPI(t, api, prior)
	if _, ok := state.Records.Elements()["Custom Title"]; !ok {
		t.Fatalf("want key Custom Title, got %v", state.Records.Elements())
	}
}

func TestMapResponseToModel_Records_NoPriorMatchUsesRecordName(t *testing.T) {
	t.Parallel()
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{"other": {}})
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "uid-1", RecordName: "Title A", CanShare: false, CanEdit: false},
		},
	}
	state := mapRecordsViaAPI(t, api, prior)
	if _, ok := state.Records.Elements()["Title A"]; !ok {
		t.Fatalf("want key Title A, got %v", state.Records.Elements())
	}
}

func TestMapResponseToModel_Records_NullPriorUsesRecordName(t *testing.T) {
	t.Parallel()
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "uid-1", RecordName: "Title A", CanShare: false, CanEdit: false},
		},
	}
	state := mapRecordsViaAPI(t, api, types.MapNull(commonsharedfolder.RecordEntryMapElemType))
	if _, ok := state.Records.Elements()["Title A"]; !ok {
		t.Fatalf("want key Title A, got %v", state.Records.Elements())
	}
}

func TestMapResponseToModel_Records_NullPriorUsesUIDWhenNameEmpty(t *testing.T) {
	t.Parallel()
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "uid-only", RecordName: "", CanShare: false, CanEdit: false},
		},
	}
	state := mapRecordsViaAPI(t, api, types.MapNull(commonsharedfolder.RecordEntryMapElemType))
	if _, ok := state.Records.Elements()["uid-only"]; !ok {
		t.Fatalf("want key uid-only, got %v", state.Records.Elements())
	}
}

func TestMapResponseToModel_Records_NameOnlyEntryIncluded(t *testing.T) {
	t.Parallel()
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "", RecordName: "name-only", CanShare: true, CanEdit: false},
		},
	}
	state := mapRecordsViaAPI(t, api, types.MapNull(commonsharedfolder.RecordEntryMapElemType))
	els := state.Records.Elements()
	if len(els) != 1 {
		t.Fatalf("len=%d want 1", len(els))
	}
	if _, ok := els["name-only"]; !ok {
		t.Fatalf("missing key name-only, keys=%v", els)
	}
}

func TestMapResponseToModel_Records_ReusesPriorKeyByUID(t *testing.T) {
	t.Parallel()
	prior := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"legacy-uid-key": {false, false},
	})
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID: "sf-1", Path: "/f",
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "legacy-uid-key", RecordName: "New Title", CanShare: true, CanEdit: true},
		},
	}
	state := mapRecordsViaAPI(t, api, prior)
	if _, ok := state.Records.Elements()["legacy-uid-key"]; !ok {
		t.Fatalf("want key legacy-uid-key, got %v", state.Records.Elements())
	}
}

func TestMapResponseToModel_RecordsMerge(t *testing.T) {
	t.Parallel()
	api := &commonsharedfolder.SharedFolderResponse{
		FolderUID:     "sf-1",
		Path:          "/folder",
		ManageUsers:   false,
		ManageRecords: false,
		CanShare:      false,
		CanEdit:       false,
		Records: []commonsharedfolder.SharedFolderRecordEntry{
			{RecordUID: "u1", RecordName: "n1", CanShare: true, CanEdit: false},
		},
		Users: nil,
	}
	priorRecords := recordsMapWithKeys(t, map[string]struct{ share, edit bool }{
		"u1": {false, false},
	})
	state := &commonsharedfolder.Model{}
	if err := commonsharedfolder.MapResponseToModel(api, state, nullPriorUsers(), priorRecords); err != nil {
		t.Fatal(err)
	}
	if state.Records.IsNull() {
		t.Fatal("records null")
	}
	if _, ok := state.Records.Elements()["u1"]; !ok {
		t.Fatalf("expected key u1, got %v", state.Records.Elements())
	}
}
