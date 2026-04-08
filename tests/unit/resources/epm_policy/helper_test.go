// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var epmPolicyAttrTypes = map[string]tftypes.Type{
	"id":                  tftypes.String,
	"managed_company":     tftypes.String,
	"policy_name":         tftypes.String,
	"policy_type":         tftypes.String,
	"status":              tftypes.String,
	"control":             tftypes.Set{ElementType: tftypes.String},
	"user_groups":         tftypes.Set{ElementType: tftypes.String},
	"machine_collections": tftypes.Set{ElementType: tftypes.String},
	"applications":        tftypes.Set{ElementType: tftypes.String},
	"day_filter":          tftypes.Set{ElementType: tftypes.String},
	"time_filter":         tftypes.Set{ElementType: tftypes.String},
	"date_filter":         tftypes.Set{ElementType: tftypes.String},
}

func epmPolicyObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: epmPolicyAttrTypes}
}

func nullStringSet() tftypes.Value {
	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
}

func stringSet(elems ...interface{}) tftypes.Value {
	if len(elems) == 0 {
		return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{})
	}
	vals := make([]tftypes.Value, len(elems))
	for i, e := range elems {
		vals[i] = tftypes.NewValue(tftypes.String, e)
	}
	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
}

// newPlanValues builds tftypes for resource plan/state. Pass nil for unknown id; null sets via empty slice vs nil map.
func newPlanValues(id, managedCompany, policyName, policyType, status interface{},
	control, userGroups, machineCollections, applications, dayFilter, timeFilter, dateFilter tftypes.Value,
) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, id),
		"managed_company":     tftypes.NewValue(tftypes.String, managedCompany),
		"policy_name":         tftypes.NewValue(tftypes.String, policyName),
		"policy_type":         tftypes.NewValue(tftypes.String, policyType),
		"status":              tftypes.NewValue(tftypes.String, status),
		"control":             control,
		"user_groups":         userGroups,
		"machine_collections": machineCollections,
		"applications":        applications,
		"day_filter":          dayFilter,
		"time_filter":         timeFilter,
		"date_filter":         dateFilter,
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema, epmPolicyObjectType()
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *epmpolicy.EpmPolicyResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

// startMockServerPolicyViewNotFound returns HTTP 500 with "not found" on poll for epm policy view only.
func startMockServerPolicyViewNotFound(mock *helpers.CommandServer) *httptest.Server {
	return helpers.StartCommandServerWithResultHook(mock,
		func(cmd string, idx int) (string, interface{}) {
			if strings.Contains(cmd, "epm policy view") {
				return "ok", samplePolicyViewData()
			}
			return "ok", nil
		},
		func(cmd string, idx int) (int, []byte) {
			if strings.Contains(cmd, "epm policy view") {
				return http.StatusInternalServerError, []byte(`{"message":"policy not found"}`)
			}
			return 0, nil
		},
	)
}

func samplePolicyViewData() map[string]interface{} {
	return map[string]interface{}{
		"PolicyId":         "55",
		"PolicyName":       "Synced",
		"PolicyType":       "CommandLine",
		"Status":           "enforce",
		"UserCheck":        []string{"u1"},
		"MachineCheck":     []string{"m1"},
		"ApplicationCheck": []string{},
		"DayCheck":         []int{1},
		"DateCheck":        []interface{}{},
		"TimeCheck":        []interface{}{},
		"Actions": map[string]interface{}{
			"OnSuccess": map[string]interface{}{"Controls": []string{"NOTIFY"}},
		},
	}
}

// httptestNewServerOK returns a mock Commander API server; caller must Close().
func httptestNewServerOK(t *testing.T) *httptest.Server {
	t.Helper()
	return startMockServer(&helpers.CommandServer{}, func(string, int) (string, interface{}) { return "ok", nil })
}
