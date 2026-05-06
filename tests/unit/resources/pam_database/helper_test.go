// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_records/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var hostnameOrIPAttrTypes = map[string]tftypes.Type{
	"hostname":            tftypes.String,
	"administrative_port": tftypes.Number,
}

var pamDatabaseAttrTypes = map[string]tftypes.Type{
	"id":    tftypes.String,
	"title": tftypes.String,
	"hostname_or_ip": tftypes.Object{
		AttributeTypes: hostnameOrIPAttrTypes,
	},
	"use_ssl":         tftypes.Bool,
	"database_id":     tftypes.String,
	"database_type":   tftypes.String,
	"provider_group":  tftypes.String,
	"provider_region": tftypes.String,
	"notes":           tftypes.String,
	"folder":          tftypes.String,
	"pam_settings":    tftypes.DynamicPseudoType,
}

func pamDatabaseObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: pamDatabaseAttrTypes}
}

func newHostnameValues(hostname interface{}, port interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"hostname":            tftypes.NewValue(tftypes.String, hostname),
		"administrative_port": tftypes.NewValue(tftypes.Number, port),
	}
}

func newPlanStateValues(
	id, title interface{},
	hostnameOrIP interface{},
	useSSL interface{},
	databaseId interface{},
	databaseType interface{},
	providerGroup, providerRegion interface{},
	notes, folder interface{},
) map[string]tftypes.Value {
	hostObjType := tftypes.Object{AttributeTypes: hostnameOrIPAttrTypes}
	var hostVal tftypes.Value
	if hostnameOrIP == nil {
		hostVal = tftypes.NewValue(hostObjType, nil)
	} else {
		hostVal = tftypes.NewValue(hostObjType, hostnameOrIP)
	}

	return map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"title":           tftypes.NewValue(tftypes.String, title),
		"hostname_or_ip":  hostVal,
		"use_ssl":         tftypes.NewValue(tftypes.Bool, useSSL),
		"database_id":     tftypes.NewValue(tftypes.String, databaseId),
		"database_type":   tftypes.NewValue(tftypes.String, databaseType),
		"provider_group":  tftypes.NewValue(tftypes.String, providerGroup),
		"provider_region": tftypes.NewValue(tftypes.String, providerRegion),
		"notes":           tftypes.NewValue(tftypes.String, notes),
		"folder":          tftypes.NewValue(tftypes.String, folder),
		"pam_settings":    tftypes.NewValue(tftypes.DynamicPseudoType, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := pamdatabase.NewPamDatabaseResource().(*pamdatabase.PamDatabaseResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema, pamDatabaseObjectType()
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *pamdatabase.PamDatabaseResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := pamdatabase.NewPamDatabaseResource().(*pamdatabase.PamDatabaseResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

func vaultRecordGetJSON(uid, title, notes string, hostname, port string) interface{} {
	return map[string]interface{}{
		"record_uid": uid,
		"type":       "pamDatabase",
		"title":      title,
		"notes":      notes,
		"fields": []map[string]interface{}{
			{
				"type":  "pamHostname",
				"label": "pamHostname",
				"value": json.RawMessage(`[{"hostName":"` + hostname + `","port":"` + port + `"}]`),
			},
			{
				"type":  "checkbox",
				"label": "useSSL",
				"value": json.RawMessage(`[true]`),
			},
			{
				"type":  "text",
				"label": "databaseId",
				"value": json.RawMessage(`["test_id"]`),
			},
			{
				"type":  "databaseType",
				"value": json.RawMessage(`["mongodb"]`),
			},
		},
	}
}
