// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	managingNodesElemType = tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"privileges": tftypes.Set{ElementType: tftypes.String},
			"cascade":    tftypes.Bool,
		},
	}
	enterpriseRoleAttrTypes = map[string]tftypes.Type{
		"id":                   tftypes.String,
		"name":                 tftypes.String,
		"node":                 tftypes.String,
		"users":                tftypes.Set{ElementType: tftypes.String},
		"teams":                tftypes.Set{ElementType: tftypes.String},
		"managing_nodes":       tftypes.Map{ElementType: managingNodesElemType},
		"enforcement_policies": tftypes.Map{ElementType: tftypes.String},
		"managed_company":      tftypes.String,
	}
)

func enterpriseRoleObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseRoleAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. Pass nil for null/optional attributes.
// users, teams: nil = null set; []interface{}{"a"} = set with elements.
// managing_nodes, enforcement_policies: nil = null map.
//
//nolint:unparam // users is always nil at call sites but used in body for null set.
func newPlanStateValues(id, name, node interface{}, users, teams interface{}, managingNodes, enforcementPolicies interface{}, managedCompany interface{}) map[string]tftypes.Value {
	var usersVal, teamsVal tftypes.Value
	if users == nil {
		usersVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := users.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		usersVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
	}
	if teams == nil {
		teamsVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := teams.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		teamsVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
	}
	var managingNodesVal, enforcementPoliciesVal tftypes.Value
	if managingNodes == nil {
		managingNodesVal = tftypes.NewValue(tftypes.Map{ElementType: managingNodesElemType}, nil)
	} else {
		managingNodesVal = tftypes.NewValue(tftypes.Map{ElementType: managingNodesElemType}, managingNodes)
	}
	if enforcementPolicies == nil {
		enforcementPoliciesVal = tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil)
	} else {
		enforcementPoliciesVal = tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, enforcementPolicies)
	}
	return map[string]tftypes.Value{
		"id":                   tftypes.NewValue(tftypes.String, id),
		"name":                 tftypes.NewValue(tftypes.String, name),
		"node":                 tftypes.NewValue(tftypes.String, node),
		"users":                usersVal,
		"teams":                teamsVal,
		"managing_nodes":       managingNodesVal,
		"enforcement_policies": enforcementPoliciesVal,
		"managed_company":      tftypes.NewValue(tftypes.String, managedCompany),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := enterpriseRoleObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *enterpriserole.EnterpriseRoleResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := enterpriserole.NewEnterpriseRoleResource().(*enterpriserole.EnterpriseRoleResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
