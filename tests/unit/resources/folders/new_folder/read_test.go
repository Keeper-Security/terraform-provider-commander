// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestNewFolderResource_Read_Success_ShareDeclared(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.HasPrefix(cmd, "nsf-get") {
			return "ok", map[string]interface{}{
				"nested_share_folder_uid": "FID-1",
				"name":                    "Engineering",
				"user_permissions": []interface{}{
					map[string]interface{}{"accessor": "kapil@metronlabs.io", "role": "owner"},
					map[string]interface{}{"accessor": "anant@metronlabs.com", "role": "content-manager"},
				},
			}
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("FID-1", "Engineering", shareMap(map[string]string{
		"anant@metronlabs.com": "viewer", // stale role; Read should refresh to content-manager
	})))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}

	// Inspect the persisted state.
	type state struct {
		Id    string            `tfsdk:"id"`
		Name  string            `tfsdk:"name"`
		Share map[string]string `tfsdk:"share"`
	}
	var got state
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.Id != "FID-1" {
		t.Errorf("id = %q, want %q", got.Id, "FID-1")
	}
	if got.Name != "Engineering" {
		t.Errorf("name = %q, want %q", got.Name, "Engineering")
	}
	if len(got.Share) != 1 {
		t.Errorf("expected exactly 1 share entry (owner filtered out), got %d (%v)", len(got.Share), got.Share)
	}
	if got.Share["anant@metronlabs.com"] != "content-manager" {
		t.Errorf("share role drift not reconciled; got %q, want %q", got.Share["anant@metronlabs.com"], "content-manager")
	}
	if _, hasOwner := got.Share["kapil@metronlabs.io"]; hasOwner {
		t.Errorf("owner entry should not appear in share state")
	}
}

func TestNewFolderResource_Read_Success_ShareNotDeclared_LeavesShareNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.HasPrefix(cmd, "nsf-get") {
			return "ok", map[string]interface{}{
				"nested_share_folder_uid": "FID-2",
				"name":                    "NoShareFolder",
				"user_permissions": []interface{}{
					map[string]interface{}{"accessor": "external@example.com", "role": "viewer"},
				},
			}
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	// share is null in prior state (user never declared it).
	rawState := tftypes.NewValue(objType, newPlanStateValues("FID-2", "NoShareFolder", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}

	type state struct {
		Id    string            `tfsdk:"id"`
		Name  string            `tfsdk:"name"`
		Share map[string]string `tfsdk:"share"`
	}
	var got state
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.Share != nil {
		t.Errorf("expected share to remain null (user did not declare it), got: %v", got.Share)
	}
}

func TestNewFolderResource_Read_NotFound_RemovesResource(t *testing.T) {
	// Only fail nsf-get (with a 500 body that contains "not found" so the api
	// package translates it into api.ErrResourceNotFound). Let sync-down pass.
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, _ int) (int, []byte) {
		if strings.HasPrefix(cmd, "nsf-get") {
			return http.StatusInternalServerError, []byte(`{"message":"folder not found"}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("FID-MISSING", "Name", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should not add a diagnostic for not-found (state should be removed): %v", resp.Diagnostics)
	}
	// resp.State.Raw is null when RemoveResource has been called.
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected state to be removed (Raw should be null), got: %v", resp.State.Raw)
	}
}

func TestNewFolderResource_Read_EmptyIdInState(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("", "Name", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when state.Id is empty")
	}
}

func TestNewFolderResource_Read_NoApiManager(t *testing.T) {
	// Build an unconfigured resource directly. Configure is never called, so
	// BaseResource.ApiManager stays nil and EnsureApiManager fails.
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)

	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("FID-1", "Name", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestNewFolderResource_Read_ApiErrorPropagates(t *testing.T) {
	// 500 without "not found" in the body -> generic API error, not RemoveResource.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"backend unavailable"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("FID-1", "Name", nil))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for non-NotFound API error")
	}
	if resp.State.Raw.IsNull() {
		t.Error("did not expect state to be removed on a generic API error")
	}
}
