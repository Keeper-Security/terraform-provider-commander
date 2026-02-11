// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEnterpriseTeamResource_Delete_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, nil)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", resp.Diagnostics)
	}
}

func TestEnterpriseTeamResource_Delete_NoApiManager(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}

func TestEnterpriseTeamResource_Delete_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed"}`))
	}))
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("team-uid-123", "Engineering", nil, nil, nil, nil, nil, "Root", nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when API returns error")
	}
}

func TestEnterpriseTeamResource_Delete_StateGetError(t *testing.T) {
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &api.ApiManager{}}, &resource.ConfigureResponse{})
	sch, _ := getSchema(t)
	wrongTypes := map[string]tftypes.Type{
		"id": tftypes.Number, "name": tftypes.String,
		"restrict_record_edit": tftypes.Bool, "restrict_record_re_share": tftypes.Bool, "enable_privacy_screen": tftypes.Bool,
		"users": tftypes.Set{ElementType: tftypes.String}, "roles": tftypes.Set{ElementType: tftypes.String},
		"node": tftypes.String, "managed_company": tftypes.String,
	}
	rawState := tftypes.NewValue(tftypes.Object{AttributeTypes: wrongTypes}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.Number, 123), "name": tftypes.NewValue(tftypes.String, "Engineering"),
		"restrict_record_edit": tftypes.NewValue(tftypes.Bool, nil), "restrict_record_re_share": tftypes.NewValue(tftypes.Bool, nil),
		"enable_privacy_screen": tftypes.NewValue(tftypes.Bool, nil),
		"users":                 tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"roles":                 tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"node":                  tftypes.NewValue(tftypes.String, "Root"), "managed_company": tftypes.NewValue(tftypes.String, nil),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when state get fails")
	}
}
