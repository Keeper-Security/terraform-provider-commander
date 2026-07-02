// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"testing"

	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestImportState_Success(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	sch, objType := getSchema(t)
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.ImportStateRequest{ID: "uid-import-123"}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.ImportState(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}
}

func TestImportState_EmptyID(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	sch, objType := getSchema(t)
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.ImportStateRequest{ID: ""}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty import ID")
	}
}

func TestImportState_WhitespaceID(t *testing.T) {
	r := pamremotebrowser.NewPamRemoteBrowserResource().(*pamremotebrowser.PamRemoteBrowserResource)
	sch, objType := getSchema(t)
	emptyState := tftypes.NewValue(objType, newPlanStateValues(
		nil, nil, nil, nil, nil, nil,
	))

	req := resource.ImportStateRequest{ID: "   "}
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.ImportState(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only import ID")
	}
}
