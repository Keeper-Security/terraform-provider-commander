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
	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var epmPolicyDSAttrTypes = map[string]tftypes.Type{
	"policy":                         tftypes.String,
	"managed_company":                tftypes.String,
	"id":                             tftypes.String,
	"policy_name":                    tftypes.String,
	"policy_type":                    tftypes.String,
	"status":                         tftypes.String,
	"message":                        tftypes.String,
	"require_policy_acknowledgement": tftypes.Bool,
	"control":                        tftypes.Set{ElementType: tftypes.String},
	"user_groups":                    tftypes.Set{ElementType: tftypes.String},
	"machine_collections":            tftypes.Set{ElementType: tftypes.String},
	"applications":                   tftypes.Set{ElementType: tftypes.String},
	"day_filter":                     tftypes.Set{ElementType: tftypes.String},
	"time_filter":                    tftypes.Set{ElementType: tftypes.String},
	"date_filter":                    tftypes.Set{ElementType: tftypes.String},
}

func epmPolicyDSObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: epmPolicyDSAttrTypes}
}

// newConfigValues builds tftypes values for datasource config. policy is required; computed fields null in config.
func newConfigValues(policy interface{}) map[string]tftypes.Value {
	nullSet := tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	return map[string]tftypes.Value{
		"policy":                         tftypes.NewValue(tftypes.String, policy),
		"managed_company":                tftypes.NewValue(tftypes.String, nil),
		"id":                             tftypes.NewValue(tftypes.String, nil),
		"policy_name":                    tftypes.NewValue(tftypes.String, nil),
		"policy_type":                    tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, nil),
		"message":                        tftypes.NewValue(tftypes.String, nil),
		"require_policy_acknowledgement": tftypes.NewValue(tftypes.Bool, nil),
		"control":                        nullSet,
		"user_groups":                    nullSet,
		"machine_collections":            nullSet,
		"applications":                   nullSet,
		"day_filter":                     nullSet,
		"time_filter":                    nullSet,
		"date_filter":                    nullSet,
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	d := epmpolicy.NewEpmPolicyDataSource().(*epmpolicy.EpmPolicyDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema, epmPolicyDSObjectType()
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *epmpolicy.EpmPolicyDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := epmpolicy.NewEpmPolicyDataSource().(*epmpolicy.EpmPolicyDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

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
