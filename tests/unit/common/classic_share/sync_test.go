// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share_test

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// nestedShareMap builds a types.Map of object values matching the
// classic_share schema element type. Pass nil to get a null map.
func nestedShareMap(t *testing.T, entries map[string]struct{ canShare, canEdit bool }) types.Map {
	t.Helper()
	objType := types.ObjectType{AttrTypes: classic_share.SharePermissionsObjectType()}
	if entries == nil {
		return types.MapNull(objType)
	}
	elems := make(map[string]attr.Value, len(entries))
	for k, v := range entries {
		obj, diags := types.ObjectValue(classic_share.SharePermissionsObjectType(), map[string]attr.Value{
			classic_share.AttrCanShare: types.BoolValue(v.canShare),
			classic_share.AttrCanEdit:  types.BoolValue(v.canEdit),
		})
		if diags.HasError() {
			t.Fatalf("ObjectValue: %v", diags)
		}
		elems[k] = obj
	}
	m, diags := types.MapValue(objType, elems)
	if diags.HasError() {
		t.Fatalf("MapValue: %v", diags)
	}
	return m
}

func newApiManager(srv string) *api.ApiManager {
	return &api.ApiManager{
		ServiceModeUrl:    srv,
		ServiceModeApiKey: "test-key",
		HttpClient:        &http.Client{},
	}
}

func TestSyncSharePermissions_EmptyRecordUIDReturnsError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	am := newApiManager(server.URL)
	plan := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true},
	})
	state := types.MapNull(types.ObjectType{AttrTypes: classic_share.SharePermissionsObjectType()})

	if err := classic_share.SyncSharePermissions(context.Background(), am, "", plan, state); err == nil {
		t.Fatal("expected error for empty recordUID")
	}
	if mock.CommandCount() != 0 {
		t.Errorf("expected no commands issued, got %d", mock.CommandCount())
	}
}

func TestSyncSharePermissions_GrantsAllOnInitialCreate(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	am := newApiManager(server.URL)
	plan := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true, canEdit: false},
		"b@x.com": {canShare: true, canEdit: true},
	})
	state := types.MapNull(types.ObjectType{AttrTypes: classic_share.SharePermissionsObjectType()})

	if err := classic_share.SyncSharePermissions(context.Background(), am, "REC1", plan, state); err != nil {
		t.Fatalf("SyncSharePermissions: %v", err)
	}
	if got, want := mock.CommandCount(), 2; got != want {
		t.Errorf("issued %d commands, want %d", got, want)
	}
}

func TestSyncSharePermissions_RevokesRemovedAndGrantsChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	var seen []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		seen = append(seen, cmd)
		return "ok", nil
	})
	defer server.Close()

	am := newApiManager(server.URL)
	state := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true, canEdit: false}, // unchanged
		"b@x.com": {canShare: true, canEdit: false}, // changed
		"c@x.com": {canShare: true, canEdit: true},  // removed
	})
	plan := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true, canEdit: false},
		"b@x.com": {canShare: true, canEdit: true},
		"d@x.com": {canShare: false, canEdit: true}, // added
	})

	if err := classic_share.SyncSharePermissions(context.Background(), am, "REC1", plan, state); err != nil {
		t.Fatalf("SyncSharePermissions: %v", err)
	}

	if got, want := mock.CommandCount(), 3; got != want {
		t.Errorf("issued %d commands, want %d (seen=%v)", got, want, seen)
	}

	sort.Strings(seen)
	mustContain(t, seen, `share-record --email 'c@x.com' 'REC1' --action revoke`)
	mustContain(t, seen, `share-record --email 'b@x.com' 'REC1' --write`)
	mustContain(t, seen, `share-record --email 'd@x.com' 'REC1' --write`)

	for _, cmd := range seen {
		if strings.Contains(cmd, "'a@x.com'") {
			t.Errorf("did not expect any command for unchanged entry a@x.com, got: %s", cmd)
		}
	}
}

func TestSyncSharePermissions_NullPlanRevokesAll(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	am := newApiManager(server.URL)
	state := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true, canEdit: false},
		"b@x.com": {canShare: false, canEdit: true},
	})
	plan := types.MapNull(types.ObjectType{AttrTypes: classic_share.SharePermissionsObjectType()})

	if err := classic_share.SyncSharePermissions(context.Background(), am, "REC1", plan, state); err != nil {
		t.Fatalf("SyncSharePermissions: %v", err)
	}
	if got, want := mock.CommandCount(), 2; got != want {
		t.Errorf("issued %d commands, want %d", got, want)
	}
}

func TestSyncSharePermissions_DowngradeToViewerEmitsRevokeFlags(t *testing.T) {
	mock := &helpers.CommandServer{}
	var seen []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		seen = append(seen, cmd)
		return "ok", nil
	})
	defer server.Close()

	am := newApiManager(server.URL)
	state := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true, canEdit: true},
	})
	plan := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: false, canEdit: false},
	})

	if err := classic_share.SyncSharePermissions(context.Background(), am, "REC1", plan, state); err != nil {
		t.Fatalf("SyncSharePermissions: %v", err)
	}
	if got, want := mock.CommandCount(), 1; got != want {
		t.Fatalf("issued %d commands, want %d (seen=%v)", got, want, seen)
	}
	mustContain(t, seen, `share-record --email 'a@x.com' 'REC1' --action revoke --share --write`)
}

func TestSyncSharePermissions_GrantErrorAborts(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "--share", nil)
	defer server.Close()

	am := newApiManager(server.URL)
	plan := nestedShareMap(t, map[string]struct{ canShare, canEdit bool }{
		"a@x.com": {canShare: true, canEdit: false},
	})
	state := types.MapNull(types.ObjectType{AttrTypes: classic_share.SharePermissionsObjectType()})

	if err := classic_share.SyncSharePermissions(context.Background(), am, "REC1", plan, state); err == nil {
		t.Error("expected error when grant API call fails")
	}
}

func mustContain(t *testing.T, sorted []string, want string) {
	t.Helper()
	for _, s := range sorted {
		if s == want {
			return
		}
	}
	t.Errorf("expected to see command %q in %v", want, sorted)
}
