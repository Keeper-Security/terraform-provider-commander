// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"context"
	"strings"
	"testing"

	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEpmPolicyResource_MetadataAndNew(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "commander"}, &resp)
	if resp.TypeName != "commander_epm_policy" {
		t.Fatalf("got %s", resp.TypeName)
	}
}

func TestEpmPolicyResource_Configure(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var resp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: nil}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "bad"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error for bad provider data")
	}
}

func TestEpmPolicyResource_SchemaAndConfigValidators(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var sresp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &sresp)
	if sresp.Schema.Attributes == nil {
		t.Fatal("schema")
	}
	vals := r.ConfigValidators(context.Background())
	if len(vals) == 0 {
		t.Fatal("config validators")
	}
}

func TestEpmPolicyResource_Create(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy add") {
			return "ok", map[string]interface{}{"policy_id": "999"}
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	planRaw := tftypes.NewValue(objType, newPlanValues(
		nil, nil, "New Policy", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	emptyState := tftypes.NewValue(objType, newPlanValues(
		nil, nil, nil, nil, nil,
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: planRaw}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_Read(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", samplePolicyViewData()
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		stringSet("Monday"), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_Read_NoId(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error without id")
	}
}

func TestEpmPolicyResource_Read_NoApiManager(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"1", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Read_NotFoundRemoves(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServerPolicyViewNotFound(mock)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics when policy missing (removed from state): %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_Read_MapPolicyViewError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", map[string]interface{}{
				"PolicyId":   "",
				"PolicyType": "Command",
				"Status":     "enforce",
			}
		}
		return "ok", nil
	})
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error when API returns empty policy id")
	}
}

func TestEpmPolicyResource_Read_UnmarshalError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", map[string]interface{}{"PolicyId": true, "PolicyType": "Command", "Status": "enforce"}
		}
		return "ok", nil
	})
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error on bad API JSON")
	}
}

func TestEpmPolicyResource_Create_ExecuteError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm policy add", nil)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	planRaw := tftypes.NewValue(objType, newPlanValues(
		nil, nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: planRaw}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error when API fails")
	}
}

func TestEpmPolicyResource_Read_EnterpriseDownError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "enterprise-down", nil)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Read_EpmSyncDownError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm sync-down", nil)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Read_ExecuteError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm policy view", nil)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Update_ExecuteError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm policy edit", nil)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	v := newPlanValues(
		"55", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	)
	raw := tftypes.NewValue(objType, v)
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: raw},
		State: tfsdk.State{Schema: sch, Raw: raw},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Delete_StateGetError(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	badRaw := tftypes.NewValue(objType, nil)
	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: badRaw}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error on null state")
	}
}

func TestEpmPolicyResource_Delete_ExecuteError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm policy delete", nil)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Update_ImmutableManagedCompany(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", "mc-a", "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	planRaw := tftypes.NewValue(objType, newPlanValues(
		"55", "mc-b", "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: planRaw},
		State: tfsdk.State{Schema: sch, Raw: stateRaw},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want immutable managed_company error")
	}
}

func TestEpmPolicyResource_ImportState_InvalidID(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	empty := tftypes.NewValue(objType, nil)
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: empty}}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "   "}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error for empty import id")
	}
}

func TestEpmPolicyResource_ImportState_EmptyResourcePart(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	empty := tftypes.NewValue(objType, nil)
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: empty}}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "managed,"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error when resource id part empty after comma")
	}
}

func TestEpmPolicyResource_Create_PlanGetError(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error on null plan")
	}
}

func TestEpmPolicyResource_Update_PlanGetError(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	bad := tftypes.NewValue(objType, nil)
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: bad},
		State: tfsdk.State{Schema: sch, Raw: bad},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Read_StateGetError(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	var resp resource.ReadResponse
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Create_NoApiManager(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	sch, objType := getSchema(t)
	planRaw := tftypes.NewValue(objType, newPlanValues(
		nil, nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: planRaw}}
	var resp resource.CreateResponse
	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Delete_NoApiManager(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Update_NoApiManager(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	sch, objType := getSchema(t)
	v := newPlanValues(
		"55", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	)
	raw := tftypes.NewValue(objType, v)
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: raw},
		State: tfsdk.State{Schema: sch, Raw: raw},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Update(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy edit") {
			return "ok", nil
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	base := newPlanValues(
		"55", nil, "Old", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	)
	planVals := newPlanValues(
		"55", nil, "New", "command", "enforce",
		stringSet("notify"), stringSet("u1"), stringSet("m1"), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	)
	stateRaw := tftypes.NewValue(objType, base)
	planRaw := tftypes.NewValue(objType, planVals)
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: planRaw},
		State: tfsdk.State{Schema: sch, Raw: stateRaw},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_Update_ImmutablePolicyType(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	planRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "P", "elevation", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: planRaw},
		State: tfsdk.State{Schema: sch, Raw: stateRaw},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want immutable error")
	}
}

func TestEpmPolicyResource_Update_EmptyStateId(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	v := newPlanValues(
		"", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	)
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: tftypes.NewValue(objType, v)},
		State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, v)},
	}
	var resp resource.UpdateResponse
	r.Update(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyResource_Delete(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(string, int) (string, interface{}) { return "ok", nil })
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	stateRaw := tftypes.NewValue(objType, newPlanValues(
		"55", nil, "P", "command", "off",
		nullStringSet(), nullStringSet(), nullStringSet(), nullStringSet(),
		nullStringSet(), nullStringSet(), nullStringSet(),
	))
	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: stateRaw}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_ImportState(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	empty := tftypes.NewValue(objType, nil)
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: empty}}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "77"}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Import: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_ImportState_WithManagedCompany(t *testing.T) {
	t.Parallel()
	server := httptestNewServerOK(t)
	defer server.Close()
	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	empty := tftypes.NewValue(objType, nil)
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: empty}}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "Acme Corp,88"}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Import: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyResource_ImportState_NoApiManager(t *testing.T) {
	t.Parallel()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var resp resource.ImportStateResponse
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "1"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}
