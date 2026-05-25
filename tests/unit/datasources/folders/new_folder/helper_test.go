// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolderds_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	newfolderds "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/new_folder"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// newFolderDSAttrTypes mirrors the data source schema attribute types for
// tftypes-based config construction in tests.
var newFolderDSAttrTypes = map[string]tftypes.Type{
	"new_folder": tftypes.String,
	"id":         tftypes.String,
	"name":       tftypes.String,
	"share":      tftypes.Map{ElementType: tftypes.String},
}

func newFolderDSObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: newFolderDSAttrTypes}
}

// newConfigValues builds a tftypes value bag for the data source config. The
// only attribute the user actually supplies is `new_folder`; id, name, and
// share are Computed and therefore null in the config phase.
func newConfigValues(lookup interface{}) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"new_folder": tftypes.NewValue(tftypes.String, lookup),
		"id":         tftypes.NewValue(tftypes.String, nil),
		"name":       tftypes.NewValue(tftypes.String, nil),
		"share":      tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
	}
}

// getDSSchema returns the data source schema and the matching tftypes object.
func getDSSchema(t *testing.T) (dschema.Schema, tftypes.Object) {
	t.Helper()
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema, newFolderDSObjectType()
}

// newConfiguredDataSource returns a configured NewFolderDataSource pointing at
// the supplied mock server.
func newConfiguredDataSource(t *testing.T, server *httptest.Server) *newfolderds.NewFolderDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
	}
	d := newfolderds.NewNewFolderDataSource().(*newfolderds.NewFolderDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

// startMockServer is a thin shorthand for the shared helpers.StartCommandServer.
func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (string, interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}
