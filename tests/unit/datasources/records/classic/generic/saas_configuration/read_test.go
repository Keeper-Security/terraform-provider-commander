// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	saasconfigurationds "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/saas_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func getDSSchema(t *testing.T) dschema.Schema {
	t.Helper()
	d := saasconfigurationds.NewSaasConfigurationDataSource().(*saasconfigurationds.SaasConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

func newConfiguredDataSource(t *testing.T, server *httptest.Server) *saasconfigurationds.SaasConfigurationDataSource {
	t.Helper()
	am := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "key",
		HttpClient:        server.Client(),
		IsMspAccount:      false,
	}
	d := saasconfigurationds.NewSaasConfigurationDataSource().(*saasconfigurationds.SaasConfigurationDataSource)
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: am}, &datasource.ConfigureResponse{})
	return d
}

func dsVaultRecordJSON(uid, title string) interface{} {
	return map[string]interface{}{
		"record_uid": uid,
		"type":       "saasConfiguration",
		"title":      title,
		"custom": []map[string]interface{}{
			{"type": "text", "label": "SaaS Type", "value": []string{"Okta"}},
			{"type": "text", "label": "AppName", "value": []string{"Example App"}},
		},
	}
}

func newDSConfigRaw(t *testing.T, sch dschema.Schema, saasConfiguration string) tftypes.Value {
	t.Helper()
	tfType := sch.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", tfType)
	}

	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		if name == "saas_configuration" {
			vals[name] = tftypes.NewValue(tftypes.String, saasConfiguration)
		} else {
			vals[name] = tftypes.NewValue(attrType, nil)
		}
	}
	return tftypes.NewValue(objType, vals)
}

func newDSEmptyState(t *testing.T, sch dschema.Schema) tftypes.Value {
	t.Helper()
	tfType := sch.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", tfType)
	}

	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(objType, vals)
}

func TestRead_DS_Success(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.Contains(cmd, "get") && strings.Contains(cmd, "--format json") {
			return "ok", dsVaultRecordJSON("uid-saas-1", "SaaS Config")
		}
		return "ok", nil
	})
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-saas-1")
	emptyState := newDSEmptyState(t, sch)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}

func TestRead_DS_EmptyLookup(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, nil)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch := getDSSchema(t)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: newDSConfigRaw(t, sch, "")}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when saas_configuration is empty")
	}
}

func TestRead_DS_WrongRecordType(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.Contains(cmd, "get") && strings.Contains(cmd, "--format json") {
			return "ok", map[string]interface{}{
				"record_uid": "uid-1",
				"type":       "login",
				"title":      "Not SaaS Configuration",
			}
		}
		return "ok", nil
	})
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch := getDSSchema(t)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: newDSConfigRaw(t, sch, "uid-1")}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when record type is not saasConfiguration")
	}
}

func TestRead_DS_NoApiManager(t *testing.T) {
	d := saasconfigurationds.NewSaasConfigurationDataSource().(*saasconfigurationds.SaasConfigurationDataSource)
	sch := getDSSchema(t)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: newDSConfigRaw(t, sch, "uid-1")}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when apiManager is not configured")
	}
}
