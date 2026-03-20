// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// enterpriseNodeAttrTypes is the tftypes map for enterprise_node resource schema.
var enterpriseNodeAttrTypes = map[string]tftypes.Type{
	"id":              tftypes.String,
	"name":            tftypes.String,
	"parent":          tftypes.String,
	"toggle_isolated": tftypes.Bool,
	"managed_company": tftypes.String,
}

func enterpriseNodeObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseNodeAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. Pass nil for null/optional attributes.
func newPlanStateValues(id, name, parent, toggleIsolated, managedCompany interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"name":            tftypes.NewValue(tftypes.String, name),
		"parent":          tftypes.NewValue(tftypes.String, parent),
		"toggle_isolated": tftypes.NewValue(tftypes.Bool, toggleIsolated),
		"managed_company": tftypes.NewValue(tftypes.String, managedCompany),
	}
}

// getSchema returns the resource schema from the enterprise_node resource.
func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := enterpriseNodeObjectType()
	return resp.Schema, objType
}

// newConfiguredResource returns an EnterpriseNodeResource configured with the given API server URL.
func newConfiguredResource(t *testing.T, server *httptest.Server) *enterprisenode.EnterpriseNodeResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := enterprisenode.NewEnterpriseNodeResource().(*enterprisenode.EnterpriseNodeResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

// startMockServer starts a mock Commander API server using shared helpers.
func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
