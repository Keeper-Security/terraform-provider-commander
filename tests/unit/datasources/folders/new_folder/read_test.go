// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolderds_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	newfolderds "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/new_folder"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// dsState mirrors the data source attributes so tests can decode the final
// state via State.Get.
type dsState struct {
	NewFolder string            `tfsdk:"new_folder"`
	Id        string            `tfsdk:"id"`
	Name      string            `tfsdk:"name"`
	Share     map[string]string `tfsdk:"share"`
}

func TestNewFolderDataSource_Read_Success_PopulatesIdNameShare(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.HasPrefix(cmd, "nsf-get") {
			return "ok", map[string]interface{}{
				"nested_share_folder_uid": "FID-DS-1",
				"name":                    "Engineering",
				"user_permissions": []interface{}{
					map[string]interface{}{"accessor": "kapil@metronlabs.io", "role": "owner"},
					map[string]interface{}{"accessor": "anant@metronlabs.com", "role": "content-manager"},
					map[string]interface{}{"accessor": "viewer@example.com", "role": "viewer"},
				},
			}
		}
		return "ok", nil
	})
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	rawCfg := tftypes.NewValue(objType, newConfigValues("Engineering"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: rawCfg}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawCfg}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}

	var got dsState
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.NewFolder != "Engineering" {
		t.Errorf("new_folder = %q, want %q", got.NewFolder, "Engineering")
	}
	if got.Id != "FID-DS-1" {
		t.Errorf("id = %q, want %q", got.Id, "FID-DS-1")
	}
	if got.Name != "Engineering" {
		t.Errorf("name = %q, want %q", got.Name, "Engineering")
	}
	if len(got.Share) != 2 {
		t.Errorf("expected 2 share entries (owner filtered out), got %d (%v)", len(got.Share), got.Share)
	}
	if got.Share["anant@metronlabs.com"] != "content-manager" {
		t.Errorf("share[anant] = %q, want content-manager", got.Share["anant@metronlabs.com"])
	}
	if got.Share["viewer@example.com"] != "viewer" {
		t.Errorf("share[viewer] = %q, want viewer", got.Share["viewer@example.com"])
	}
	if _, hasOwner := got.Share["kapil@metronlabs.io"]; hasOwner {
		t.Errorf("owner entry should not appear in share state")
	}
}

func TestNewFolderDataSource_Read_Success_EmptyUserPermissions_EmptyShare(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.HasPrefix(cmd, "nsf-get") {
			return "ok", map[string]interface{}{
				"nested_share_folder_uid": "FID-DS-2",
				"name":                    "EmptyFolder",
				"user_permissions":        []interface{}{},
			}
		}
		return "ok", nil
	})
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	rawCfg := tftypes.NewValue(objType, newConfigValues("EmptyFolder"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: rawCfg}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawCfg}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}

	var got dsState
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.Share == nil {
		t.Error("share should be an empty (non-null) map, not nil")
	}
	if len(got.Share) != 0 {
		t.Errorf("expected empty share, got %v", got.Share)
	}
}

func TestNewFolderDataSource_Read_EmptyLookup_Errors(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	rawCfg := tftypes.NewValue(objType, newConfigValues(""))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: rawCfg}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawCfg}}
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when new_folder lookup value is empty")
	}
}

func TestNewFolderDataSource_Read_NotFound_AddsDiagnostic(t *testing.T) {
	// Sync-down succeeds; nsf-get returns 500 with "not found" body so the api
	// package translates it to api.ErrResourceNotFound. A data source should
	// surface this as a diagnostic error (unlike a resource which removes
	// itself from state).
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, _ int) (int, []byte) {
		if strings.HasPrefix(cmd, "nsf-get") {
			return http.StatusInternalServerError, []byte(`{"message":"folder not found"}`)
		}
		return 0, nil
	})
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	rawCfg := tftypes.NewValue(objType, newConfigValues("Missing"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: rawCfg}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawCfg}}
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when nsf-get returns not-found")
	}
}

func TestNewFolderDataSource_Read_NoApiManager_Errors(t *testing.T) {
	// Build an unconfigured data source directly. Configure is never called,
	// so BaseDataSource.ApiManager stays nil and EnsureApiManager fails.
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)

	sch, objType := getDSSchema(t)
	rawCfg := tftypes.NewValue(objType, newConfigValues("Engineering"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: rawCfg}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawCfg}}
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestNewFolderDataSource_Read_ApiErrorPropagates(t *testing.T) {
	// 500 without "not found" -> generic API error surfaces as a diagnostic.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"backend unavailable"}`))
	}))
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getDSSchema(t)
	rawCfg := tftypes.NewValue(objType, newConfigValues("Engineering"))

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: rawCfg}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawCfg}}
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for non-NotFound API error")
	}
}
