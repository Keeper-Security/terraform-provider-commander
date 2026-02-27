// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/managed_company"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var manageCompanyDataSourceAttrTypes = map[string]tftypes.Type{
	"managed_company": tftypes.String,
	"id":              tftypes.String,
	"name":            tftypes.String,
	"node":            tftypes.String,
	"plan":            tftypes.String,
	"file_plan":       tftypes.String,
}

func manageCompanyDataSourceObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: manageCompanyDataSourceAttrTypes}
}

// newConfigValues builds tftypes values for datasource config. managed_company is required; computed fields null in config.
func newConfigValues(managedCompany interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"managed_company": tftypes.NewValue(tftypes.String, managedCompany),
		"id":              tftypes.NewValue(tftypes.String, nil),
		"name":            tftypes.NewValue(tftypes.String, nil),
		"node":            tftypes.NewValue(tftypes.String, nil),
		"plan":            tftypes.NewValue(tftypes.String, nil),
		"file_plan":       tftypes.NewValue(tftypes.String, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	objType := manageCompanyDataSourceObjectType()
	return resp.Schema, objType
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *managedcompany.ManagedCompanyDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := managedcompany.NewManagedCompanyDataSource().(*managedcompany.ManagedCompanyDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
