// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_team"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterpriseTeamDataSourceAttrTypes = map[string]tftypes.Type{
	"team":            tftypes.String,
	"id":              tftypes.String,
	"name":            tftypes.String,
	"users":           tftypes.Set{ElementType: tftypes.String},
	"roles":           tftypes.Set{ElementType: tftypes.String},
	"managed_company": tftypes.String,
}

func enterpriseTeamDataSourceObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseTeamDataSourceAttrTypes}
}

// newConfigValues builds tftypes values for datasource config. team is required; managed_company optional; computed fields null in config.
//
//nolint:unparam // managedCompany is always nil at call sites but used for managed_company attribute.
func newConfigValues(team, managedCompany interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"team":            tftypes.NewValue(tftypes.String, team),
		"id":              tftypes.NewValue(tftypes.String, nil),
		"name":            tftypes.NewValue(tftypes.String, nil),
		"users":           tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"roles":           tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"managed_company": tftypes.NewValue(tftypes.String, managedCompany),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	objType := enterpriseTeamDataSourceObjectType()
	return resp.Schema, objType
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *enterpriseteam.EnterpriseTeamDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := enterpriseteam.NewEnterpriseTeamDataSource().(*enterpriseteam.EnterpriseTeamDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
