// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share_test

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringMap is a tiny helper to build a types.Map[string,string] in tests.
func stringMap(t *testing.T, entries map[string]string) types.Map {
	t.Helper()
	if entries == nil {
		return types.MapNull(types.StringType)
	}
	elems := make(map[string]attr.Value, len(entries))
	for k, v := range entries {
		elems[k] = types.StringValue(v)
	}
	m, diags := types.MapValue(types.StringType, elems)
	if diags.HasError() {
		t.Fatalf("types.MapValue: %v", diags)
	}
	return m
}

// newApiManager wires an ApiManager pointing at the given mock server.
func newApiManager(srv string) *api.ApiManager {
	return &api.ApiManager{
		ServiceModeUrl:    srv,
		ServiceModeApiKey: "test-key",
		HttpClient:        &http.Client{},
	}
}

func TestSyncSharePermissions_EmptyIdReturnsError(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	am := newApiManager(server.URL)
	plan := stringMap(t, map[string]string{"a@x.com": "viewer"})
	state := types.MapNull(types.StringType)

	err := new_share.SyncSharePermissions(context.Background(), am, new_share.CmdShareFolder, "", plan, state)
	if err == nil {
		t.Fatal("expected error for empty id")
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
	plan := stringMap(t, map[string]string{
		"a@x.com": "viewer",
		"b@x.com": "share-manager",
	})
	state := types.MapNull(types.StringType)

	if err := new_share.SyncSharePermissions(context.Background(), am, new_share.CmdShareFolder, "FID", plan, state); err != nil {
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
	// state: a=viewer, b=viewer, c=share-manager
	// plan : a=viewer (unchanged), b=full-manager (changed), d=content-manager (added); c removed
	state := stringMap(t, map[string]string{
		"a@x.com": "viewer",
		"b@x.com": "viewer",
		"c@x.com": "share-manager",
	})
	plan := stringMap(t, map[string]string{
		"a@x.com": "viewer",
		"b@x.com": "full-manager",
		"d@x.com": "content-manager",
	})

	if err := new_share.SyncSharePermissions(context.Background(), am, new_share.CmdShareRecord, "REC1", plan, state); err != nil {
		t.Fatalf("SyncSharePermissions: %v", err)
	}

	// Expected: 1 revoke for c, 2 grants for b (changed) and d (added). a is unchanged.
	if got, want := mock.CommandCount(), 3; got != want {
		t.Errorf("issued %d commands, want %d (seen=%v)", got, want, seen)
	}

	sort.Strings(seen)
	mustContain(t, seen, `nsf-share-record "REC1" --email='c@x.com' --action=revoke`)
	mustContain(t, seen, `nsf-share-record "REC1" --email='b@x.com' --action=grant --role='full-manager'`)
	mustContain(t, seen, `nsf-share-record "REC1" --email='d@x.com' --action=grant --role='content-manager'`)

	for _, cmd := range seen {
		if strings.Contains(cmd, "a@x.com") {
			t.Errorf("did not expect any command for unchanged entry a@x.com, got: %s", cmd)
		}
	}
}

func TestSyncSharePermissions_NullPlanRevokesAll(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	am := newApiManager(server.URL)
	state := stringMap(t, map[string]string{
		"a@x.com": "viewer",
		"b@x.com": "share-manager",
	})
	plan := types.MapNull(types.StringType)

	if err := new_share.SyncSharePermissions(context.Background(), am, new_share.CmdShareFolder, "F1", plan, state); err != nil {
		t.Fatalf("SyncSharePermissions: %v", err)
	}
	if got, want := mock.CommandCount(), 2; got != want {
		t.Errorf("issued %d commands, want %d", got, want)
	}
}

func TestSyncSharePermissions_GrantErrorAborts(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "--action=grant", nil)
	defer server.Close()

	am := newApiManager(server.URL)
	plan := stringMap(t, map[string]string{"a@x.com": "viewer"})
	state := types.MapNull(types.StringType)

	if err := new_share.SyncSharePermissions(context.Background(), am, new_share.CmdShareFolder, "F1", plan, state); err == nil {
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
