// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_user"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterpriseUserDataSourceAttrTypes = map[string]tftypes.Type{
	"user":            tftypes.String,
	"id":              tftypes.String,
	"name":            tftypes.String,
	"email":           tftypes.String,
	"job_title":       tftypes.String,
	"roles":           tftypes.Set{ElementType: tftypes.String},
	"teams":           tftypes.Set{ElementType: tftypes.String},
	"managed_company": tftypes.String,
}

func enterpriseUserDataSourceObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseUserDataSourceAttrTypes}
}

// newConfigValues builds tftypes values for datasource config. user is required; managed_company optional; computed fields null in config.
//
//nolint:unparam // managedCompany is always nil at call sites but used for managed_company attribute.
func newConfigValues(user, managedCompany interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"user":            tftypes.NewValue(tftypes.String, user),
		"id":              tftypes.NewValue(tftypes.String, nil),
		"name":            tftypes.NewValue(tftypes.String, nil),
		"email":           tftypes.NewValue(tftypes.String, nil),
		"job_title":       tftypes.NewValue(tftypes.String, nil),
		"roles":           tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"teams":           tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"managed_company": tftypes.NewValue(tftypes.String, managedCompany),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	objType := enterpriseUserDataSourceObjectType()
	return resp.Schema, objType
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *enterpriseuser.EnterpriseUserDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := enterpriseuser.NewEnterpriseUserDataSource().(*enterpriseuser.EnterpriseUserDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
