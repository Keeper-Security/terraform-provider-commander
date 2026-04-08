// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_user"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterpriseUserAttrTypes = map[string]tftypes.Type{
	"id":              tftypes.String,
	"email":           tftypes.String,
	"name":            tftypes.String,
	"job_title":       tftypes.String,
	"roles":           tftypes.Set{ElementType: tftypes.String},
	"teams":           tftypes.Set{ElementType: tftypes.String},
	"node":            tftypes.String,
	"managed_company": tftypes.String,
	"status":          tftypes.String,
}

func enterpriseUserObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseUserAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. Pass nil for null/optional attributes.
// roles and teams: nil = null set; []interface{}{"a","b"} = set with elements.
//
//nolint:unparam // roles is always nil at call sites but used in body for null set.
func newPlanStateValues(id, email, name, jobTitle, roles, teams, node, managedCompany, status interface{}) map[string]tftypes.Value {
	var rolesVal, teamsVal tftypes.Value
	if roles == nil {
		rolesVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := roles.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		rolesVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
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
	return map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"email":           tftypes.NewValue(tftypes.String, email),
		"name":            tftypes.NewValue(tftypes.String, name),
		"job_title":       tftypes.NewValue(tftypes.String, jobTitle),
		"roles":           rolesVal,
		"teams":           teamsVal,
		"node":            tftypes.NewValue(tftypes.String, node),
		"managed_company": tftypes.NewValue(tftypes.String, managedCompany),
		"status":          tftypes.NewValue(tftypes.String, status),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := enterpriseUserObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *enterpriseuser.EnterpriseUserResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := enterpriseuser.NewEnterpriseUserResource().(*enterpriseuser.EnterpriseUserResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
