// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_node"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterpriseNodeDataSourceAttrTypes = map[string]tftypes.Type{
	"node":            tftypes.String,
	"id":              tftypes.String,
	"name":            tftypes.String,
	"parent":          tftypes.String,
	"parent_id":       tftypes.String,
	"managed_company": tftypes.String,
}

func enterpriseNodeDataSourceObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseNodeDataSourceAttrTypes}
}

// newConfigValues builds tftypes values for datasource config. node is required; managed_company optional; computed fields null in config.
//
//nolint:unparam // managedCompany is always nil at call sites but used for managed_company attribute.
func newConfigValues(node, managedCompany interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"node":            tftypes.NewValue(tftypes.String, node),
		"id":              tftypes.NewValue(tftypes.String, nil),
		"name":            tftypes.NewValue(tftypes.String, nil),
		"parent":          tftypes.NewValue(tftypes.String, nil),
		"parent_id":       tftypes.NewValue(tftypes.String, nil),
		"managed_company": tftypes.NewValue(tftypes.String, managedCompany),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	objType := enterpriseNodeDataSourceObjectType()
	return resp.Schema, objType
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *enterprisenode.EnterpriseNodesDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := enterprisenode.NewEnterpriseNodesDataSource().(*enterprisenode.EnterpriseNodesDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
