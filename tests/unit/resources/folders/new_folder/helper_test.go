// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// newFolderAttrTypes mirrors the resource schema attribute types for tftypes-
// based plan/state construction in tests.
var newFolderAttrTypes = map[string]tftypes.Type{
	"id":              tftypes.String,
	"name":            tftypes.String,
	"folder_location": tftypes.String,
	"records":         tftypes.Set{ElementType: tftypes.String},
	"share":           tftypes.Map{ElementType: tftypes.String},
}

func newFolderObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: newFolderAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. Pass nil for null
// attributes; pass a map[string]tftypes.Value for share (or nil for null).
// folder_location and records are Computed and default to null.
func newPlanStateValues(id, name interface{}, share map[string]tftypes.Value) map[string]tftypes.Value {
	var shareVal tftypes.Value
	if share == nil {
		shareVal = tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil)
	} else {
		shareVal = tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, share)
	}
	return map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"name":            tftypes.NewValue(tftypes.String, name),
		"folder_location": tftypes.NewValue(tftypes.String, nil),
		"records":         tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"share":           shareVal,
	}
}

// shareMap is a tftypes.Value builder for the share attribute.
func shareMap(entries map[string]string) map[string]tftypes.Value {
	out := make(map[string]tftypes.Value, len(entries))
	for k, v := range entries {
		out[k] = tftypes.NewValue(tftypes.String, v)
	}
	return out
}

// getSchema returns the resource schema and the matching tftypes object.
func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema, newFolderObjectType()
}

// newConfiguredResource returns a configured NewFolderResource pointing at the
// supplied mock server.
func newConfiguredResource(t *testing.T, server *httptest.Server) *newfolder.NewFolderResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
	}
	r := newfolder.NewNewFolderResource().(*newfolder.NewFolderResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

// startMockServer is a thin shorthand for the shared helpers.StartCommandServer.
func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (string, interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
