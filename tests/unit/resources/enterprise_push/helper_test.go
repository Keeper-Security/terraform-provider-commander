// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisepush_test

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var enterprisePushAttrTypes = map[string]tftypes.Type{
	"id":                  tftypes.String,
	"file_path":           tftypes.String,
	"file_content_sha256": tftypes.String,
	"email":               tftypes.Set{ElementType: tftypes.String},
	"team":                tftypes.Set{ElementType: tftypes.String},
	"managed_company":     tftypes.String,
}

func enterprisePushObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: enterprisePushAttrTypes}
}

// newPlanStateValues builds tftypes values for plan/state. email and team: nil = null set; []interface{}{"a","b"} = set with elements.
func newPlanStateValues(id, filePath, fileContentSha256 interface{}, emailSet, teamSet interface{}) map[string]tftypes.Value {
	var emailVal, teamVal tftypes.Value
	if emailSet == nil {
		emailVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := emailSet.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		emailVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
	}
	if teamSet == nil {
		teamVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := teamSet.([]interface{})
		vals := make([]tftypes.Value, len(elems))
		for i, e := range elems {
			vals[i] = tftypes.NewValue(tftypes.String, e)
		}
		teamVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
	}
	return map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, id),
		"file_path":           tftypes.NewValue(tftypes.String, filePath),
		"file_content_sha256": tftypes.NewValue(tftypes.String, fileContentSha256),
		"email":               emailVal,
		"team":                teamVal,
		"managed_company":     tftypes.NewValue(tftypes.String, nil),
	}
}

func getSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType := enterprisePushObjectType()
	return resp.Schema, objType
}

func newConfiguredResource(t *testing.T, server *httptest.Server) *enterprisepush.EnterprisePushResource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: am}, &resource.ConfigureResponse{})
	return r
}

func startMockServer(mock *helpers.CommandServer, responseForCommand func(cmd string, idx int) (message string, data interface{})) *httptest.Server {
	return helpers.StartCommandServer(mock, responseForCommand)
}

// defaultPushDataJSON is the JSON content used by createTempJSONFile for tests.
const defaultPushDataJSON = `{"records":[]}`

// createTempJSONFile creates a temp file with valid JSON content. The file lives in t.TempDir() (cleaned automatically).
func createTempJSONFile(t *testing.T) (path string) {
	t.Helper()
	dir := t.TempDir()
	path = filepath.Join(dir, "push-data.json")
	if err := os.WriteFile(path, []byte(defaultPushDataJSON), 0600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}
