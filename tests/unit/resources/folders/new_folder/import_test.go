// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})

	req := resource.ImportStateRequest{ID: ""}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestNewFolderResource_ImportState_WhitespaceOnlyID(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})

	req := resource.ImportStateRequest{ID: "   \t  "}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only import ID (should be trimmed and rejected)")
	}
}

func TestNewFolderResource_ImportState_NoApiManager(t *testing.T) {
	// Build an unconfigured resource directly. Configure is never called, so
	// BaseResource.ApiManager stays nil and EnsureApiManager fails before any
	// import ID validation.
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)

	req := resource.ImportStateRequest{ID: "E6laPVJ1T3-sWchJCRaWOg"}
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is nil")
	}
}

func TestNewFolderResource_ImportState_Success_ByUID(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})

	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil))

	req := resource.ImportStateRequest{ID: "E6laPVJ1T3-sWchJCRaWOg"}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyRaw}}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}

	var got importedState
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.Id != "E6laPVJ1T3-sWchJCRaWOg" {
		t.Errorf("imported id = %q, want %q", got.Id, "E6laPVJ1T3-sWchJCRaWOg")
	}
	if got.Name != nil {
		t.Errorf("name should be null after import (Read refreshes it), got %q", *got.Name)
	}
	if got.Share != nil {
		t.Errorf("share should be null after import (Optional-only semantics), got %v", got.Share)
	}
}

func TestNewFolderResource_ImportState_Success_ByName_TrimsWhitespace(t *testing.T) {
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})

	sch, objType := getSchema(t)
	emptyRaw := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil))

	req := resource.ImportStateRequest{ID: "  Engineering  "}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyRaw}}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}

	var got importedState
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if got.Id != "Engineering" {
		t.Errorf("imported id should be trimmed; got %q, want %q", got.Id, "Engineering")
	}
}
