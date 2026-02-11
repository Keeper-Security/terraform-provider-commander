// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/manage_company"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var manageCompanyAttrTypes = map[string]tftypes.Type{
	"id":        tftypes.String,
	"name":      tftypes.String,
	"node":      tftypes.String,
	"seats":     tftypes.Number,
	"plan":      tftypes.String,
	"file_plan": tftypes.String,
	"add_ons":   tftypes.Set{ElementType: tftypes.String},
}

func manageCompanyObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: manageCompanyAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. Pass nil for null/optional attributes.
// add_ons: nil = null set; []interface{}{"secrets_manager"} = set with elements.
//
//nolint:unparam // addOns is always nil at call sites but used in body for null set.
func newPlanStateValues(id, name, node interface{}, seats interface{}, plan, filePlan interface{}, addOns interface{}) map[string]tftypes.Value {
	var addOnsVal tftypes.Value
	if addOns == nil {
		addOnsVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := addOns.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		addOnsVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
	}
	return map[string]tftypes.Value{
		"id":        tftypes.NewValue(tftypes.String, id),
		"name":      tftypes.NewValue(tftypes.String, name),
		"node":      tftypes.NewValue(tftypes.String, node),
		"seats":     tftypes.NewValue(tftypes.Number, seats),
		"plan":      tftypes.NewValue(tftypes.String, plan),
		"file_plan": tftypes.NewValue(tftypes.String, filePlan),
		"add_ons":   addOnsVal,
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := manageCompanyObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *managecompany.ManageCompanyResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := managecompany.NewManageCompanyResource().(*managecompany.ManageCompanyResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
