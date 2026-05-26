// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// importedState mirrors the resource attributes so tests can decode the
// post-import state via State.Get.
type importedState struct {
	Id    string            `tfsdk:"id"`
	Name  *string           `tfsdk:"name"`
	Share map[string]string `tfsdk:"share"`
}

func TestNewFolderResource_ImportState_EmptyID(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
	if mock.CommandCount() != 0 {
		t.Errorf("expected no commands during import; got %d", mock.CommandCount())
	}
}

func TestNewFolderResource_ImportState_WhitespaceOnlyID(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	req := resource.ImportStateRequest{ID: "   \t  "}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only import ID (should be trimmed and rejected)")
	}
	if mock.CommandCount() != 0 {
		t.Errorf("expected no commands during import; got %d", mock.CommandCount())
	}
}

func TestNewFolderResource_ImportState_NoApiManager(t *testing.T) {
	// Build an unconfigured resource directly. Configure is never called, so
	// BaseResource.ApiManager stays nil and EnsureApiManager fails first.
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)

	req := resource.ImportStateRequest{ID: "Cuuc9aK6VuATH49ewBf0zg"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is nil")
	}
}

func TestNewFolderResource_ImportState_Success_SetsIdAndNullDefaults(t *testing.T) {
	// ImportState does not call the API; it just seeds Id (verbatim, after
	// trim) and leaves Name and Share as null. The subsequent automatic Read
	// hydrates the rest from nsf-get.
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil))

	req := resource.ImportStateRequest{ID: "  Cuuc9aK6VuATH49ewBf0zg  "}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyRaw}}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}

	// Sanity: ImportState makes no API calls; Read will, on the next step.
	if mock.CommandCount() != 0 {
		t.Errorf("expected no API calls during ImportState; got %d", mock.CommandCount())
	}

	var got importedState
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.Id != "Cuuc9aK6VuATH49ewBf0zg" {
		t.Errorf("imported id = %q (want trimmed UID %q)", got.Id, "Cuuc9aK6VuATH49ewBf0zg")
	}
	if got.Name != nil {
		t.Errorf("name should be null after import (Read populates it), got %q", *got.Name)
	}
	if got.Share != nil {
		t.Errorf("share should be null after import (Read populates it), got %v", got.Share)
	}
}

func TestNewFolderResource_ImportState_NullValueType_MatchesSchema(t *testing.T) {
	// Regression guard: the null we set for `share` must use the schema's
	// element type (types.StringType via new_share.ShareEntryAttrType).
	// Mis-matched element types would surface as a state conversion error.
	if got := types.MapNull(new_share.ShareEntryAttrType); !got.IsNull() {
		t.Errorf("MapNull(ShareEntryAttrType) should report IsNull()=true, got %v", got)
	}
	if et := types.MapNull(new_share.ShareEntryAttrType).ElementType(context.Background()); et == nil || et.String() != types.StringType.String() {
		t.Errorf("MapNull element type = %v, want %v", et, types.StringType)
	}
}
