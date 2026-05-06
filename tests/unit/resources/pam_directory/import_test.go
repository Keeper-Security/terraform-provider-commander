// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"testing"

	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_records/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestImportState_Success(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	sch, objType := getSchema(t)
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.ImportStateRequest{ID: "uid-import-123"}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}

func TestImportState_EmptyID(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	sch, objType := getSchema(t)
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.ImportStateRequest{ID: ""}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestImportState_WhitespaceID(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	sch, objType := getSchema(t)
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.ImportStateRequest{ID: "   "}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only import ID")
	}
}
