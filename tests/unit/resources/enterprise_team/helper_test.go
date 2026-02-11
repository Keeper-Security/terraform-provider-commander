// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterpriseTeamAttrTypes = map[string]tftypes.Type{
	"id":                       tftypes.String,
	"name":                     tftypes.String,
	"restrict_record_edit":     tftypes.Bool,
	"restrict_record_re_share": tftypes.Bool,
	"enable_privacy_screen":    tftypes.Bool,
	"users":                    tftypes.Set{ElementType: tftypes.String},
	"roles":                    tftypes.Set{ElementType: tftypes.String},
	"node":                     tftypes.String,
	"managed_company":          tftypes.String,
}

func enterpriseTeamObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterpriseTeamAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. Pass nil for null/optional attributes.
// users, roles: nil = null set; []interface{}{"a"} = set with elements.
func newPlanStateValues(id, name interface{}, restrictEdit, restrictShare, restrictView interface{}, users, roles interface{}, node, managedCompany interface{}) map[string]tftypes.Value {
	var usersVal, rolesVal tftypes.Value
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
	return map[string]tftypes.Value{
		"id":                       tftypes.NewValue(tftypes.String, id),
		"name":                     tftypes.NewValue(tftypes.String, name),
		"restrict_record_edit":     tftypes.NewValue(tftypes.Bool, restrictEdit),
		"restrict_record_re_share": tftypes.NewValue(tftypes.Bool, restrictShare),
		"enable_privacy_screen":    tftypes.NewValue(tftypes.Bool, restrictView),
		"users":                    usersVal,
		"roles":                    rolesVal,
		"node":                     tftypes.NewValue(tftypes.String, node),
		"managed_company":          tftypes.NewValue(tftypes.String, managedCompany),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := enterpriseTeamObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *enterpriseteam.EnterpriseTeamResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := enterpriseteam.NewEnterpriseTeamResource().(*enterpriseteam.EnterpriseTeamResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
