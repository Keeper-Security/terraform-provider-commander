// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_role"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	managingNodesElemType = tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"privileges": tftypes.Set{ElementType: tftypes.String},
			"cascade":    tftypes.Bool,
		},
	}
	enterpriseRoleDataSourceAttrTypes = map[string]tftypes.Type{
		"role":                 tftypes.String,
		"id":                   tftypes.String,
		"name":                 tftypes.String,
		"users":                tftypes.Set{ElementType: tftypes.String},
		"teams":                tftypes.Set{ElementType: tftypes.String},
		"managing_nodes":       tftypes.Map{ElementType: managingNodesElemType},
		"enforcement_policies": tftypes.Map{ElementType: tftypes.String},
		"managed_company":      tftypes.String,
	}
)

func enterpriseRoleDataSourceObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseRoleDataSourceAttrTypes}
}

// newConfigValues builds tftypes values for datasource config. role is required; managed_company optional; computed fields null in config.
//
//nolint:unparam // managedCompany is always nil at call sites but used for managed_company attribute.
func newConfigValues(role, managedCompany interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"role":                 tftypes.NewValue(tftypes.String, role),
		"id":                   tftypes.NewValue(tftypes.String, nil),
		"name":                 tftypes.NewValue(tftypes.String, nil),
		"users":                tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"teams":                tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"managing_nodes":       tftypes.NewValue(tftypes.Map{ElementType: managingNodesElemType}, nil),
		"enforcement_policies": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		"managed_company":      tftypes.NewValue(tftypes.String, managedCompany),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	objType := enterpriseRoleDataSourceObjectType()
	return resp.Schema, objType
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *enterpriserole.EnterpriseRoleDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := enterpriserole.NewEnterpriseRoleDataSource().(*enterpriserole.EnterpriseRoleDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
