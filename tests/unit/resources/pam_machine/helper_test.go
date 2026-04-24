// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_records/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var hostnameOrIPAttrTypes = map[string]tftypes.Type{
	"hostname":            tftypes.String,
	"administrative_port": tftypes.Number,
}

var pamMachineAttrTypes = map[string]tftypes.Type{
	"id":    tftypes.String,
	"title": tftypes.String,
	"hostname_or_ip": tftypes.Object{
		AttributeTypes: hostnameOrIPAttrTypes,
	},
	"operating_system": tftypes.String,
	"instance_name":    tftypes.String,
	"instance_id":      tftypes.String,
	"provider_group":   tftypes.String,
	"provider_region":  tftypes.String,
	"notes":            tftypes.String,
	"folder":           tftypes.String,
	"pam_settings":     tftypes.DynamicPseudoType,
}

func pamMachineObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: pamMachineAttrTypes}
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
	operatingSystem, instanceName, instanceId interface{},
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
		"id":               tftypes.NewValue(tftypes.String, id),
		"title":            tftypes.NewValue(tftypes.String, title),
		"hostname_or_ip":   hostVal,
		"operating_system": tftypes.NewValue(tftypes.String, operatingSystem),
		"instance_name":    tftypes.NewValue(tftypes.String, instanceName),
		"instance_id":      tftypes.NewValue(tftypes.String, instanceId),
		"provider_group":   tftypes.NewValue(tftypes.String, providerGroup),
		"provider_region":  tftypes.NewValue(tftypes.String, providerRegion),
		"notes":            tftypes.NewValue(tftypes.String, notes),
		"folder":           tftypes.NewValue(tftypes.String, folder),
		"pam_settings":     tftypes.NewValue(tftypes.DynamicPseudoType, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema, pamMachineObjectType()
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *pammachine.PamMachineResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
