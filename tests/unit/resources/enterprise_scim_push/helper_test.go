// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterpriseScimPushAttrTypes = map[string]tftypes.Type{
	"id":              tftypes.String,
	"scim_id":         tftypes.String,
	"source":          tftypes.String,
	"record":          tftypes.String,
	"auto_approve":    tftypes.Bool,
	"managed_company": tftypes.String,
}

func enterpriseScimPushObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseScimPushAttrTypes}
}

// newPlanStateValues: id, scim_id, source, record, auto_approve, managed_company. Use nil for null.
func newPlanStateValues(id, scimId, source, record interface{}, autoApprove interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"scim_id":         tftypes.NewValue(tftypes.String, scimId),
		"source":          tftypes.NewValue(tftypes.String, source),
		"record":          tftypes.NewValue(tftypes.String, record),
		"auto_approve":    tftypes.NewValue(tftypes.Bool, autoApprove),
		"managed_company": tftypes.NewValue(tftypes.String, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := enterpriseScimPushObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *enterprisescimpush.EnterpriseScimPushResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
