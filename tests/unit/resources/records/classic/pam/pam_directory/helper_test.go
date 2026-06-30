// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var hostnameOrIPAttrTypes = map[string]tftypes.Type{
	"hostname":            tftypes.String,
	"administrative_port": tftypes.Number,
}

var alternativeIPsSetType = tftypes.Set{ElementType: tftypes.String}

var shareElementAttrTypes = map[string]tftypes.Type{
	"can_share": tftypes.Bool,
	"can_edit":  tftypes.Bool,
}

var shareMapType = tftypes.Map{ElementType: tftypes.Object{AttributeTypes: shareElementAttrTypes}}

var pamDirectoryAttrTypes = map[string]tftypes.Type{
	"id":    tftypes.String,
	"title": tftypes.String,
	"hostname_or_ip": tftypes.Object{
		AttributeTypes: hostnameOrIPAttrTypes,
	},
	"use_ssl":         tftypes.Bool,
	"domain_name":     tftypes.String,
	"alternative_ips": alternativeIPsSetType,
	"directory_id":    tftypes.String,
	"directory_type":  tftypes.String,
	"user_match":      tftypes.String,
	"provider_group":  tftypes.String,
	"provider_region": tftypes.String,
	"notes":           tftypes.String,
	"folder_location": tftypes.String,
	"pam_settings":    tftypes.DynamicPseudoType,
	"share":           shareMapType,
}

func pamDirectoryObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: pamDirectoryAttrTypes}
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
	domainName interface{},
	alternativeIPs interface{},
	directoryId interface{},
	directoryType interface{},
	userMatch interface{},
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

	var altIPsVal tftypes.Value
	if alternativeIPs == nil {
		altIPsVal = tftypes.NewValue(alternativeIPsSetType, nil)
	} else {
		altIPsVal = tftypes.NewValue(alternativeIPsSetType, alternativeIPs)
	}

	return map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"title":           tftypes.NewValue(tftypes.String, title),
		"hostname_or_ip":  hostVal,
		"use_ssl":         tftypes.NewValue(tftypes.Bool, useSSL),
		"domain_name":     tftypes.NewValue(tftypes.String, domainName),
		"alternative_ips": altIPsVal,
		"directory_id":    tftypes.NewValue(tftypes.String, directoryId),
		"directory_type":  tftypes.NewValue(tftypes.String, directoryType),
		"user_match":      tftypes.NewValue(tftypes.String, userMatch),
		"provider_group":  tftypes.NewValue(tftypes.String, providerGroup),
		"provider_region": tftypes.NewValue(tftypes.String, providerRegion),
		"notes":           tftypes.NewValue(tftypes.String, notes),
		"folder_location": tftypes.NewValue(tftypes.String, folder),
		"pam_settings":    tftypes.NewValue(tftypes.DynamicPseudoType, nil),
		"share":           tftypes.NewValue(shareMapType, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema, pamDirectoryObjectType()
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *pamdirectory.PamDirectoryResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
