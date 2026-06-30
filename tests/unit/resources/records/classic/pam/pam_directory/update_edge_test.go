// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUpdate_UseSSLClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, true, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_DomainNameClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, "old-domain.com", nil, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_DirectoryTypeClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, "active_directory", nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_AlternativeIPsChanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	oldIPs := []tftypes.Value{tftypes.NewValue(tftypes.String, "10.0.0.1")}
	newIPs := []tftypes.Value{
		tftypes.NewValue(tftypes.String, "10.0.0.2"),
		tftypes.NewValue(tftypes.String, "10.0.0.3"),
	}
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, oldIPs, nil, nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, newIPs, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_DirectoryIdClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, "old-dir-id", nil, nil, nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_UserMatchClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, "OU=Users", nil, nil, nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_ProviderFieldsClearedToNull(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, "old-group", "us-west-1", nil, nil,
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}

func TestUpdate_MoveRecordFails(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "mv", nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, "old-folder",
	))
	rawPlan := tftypes.NewValue(objType, newPlanStateValues(
		"uid-1", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, "new-folder",
	))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: rawState},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when mv command fails")
	}
}
