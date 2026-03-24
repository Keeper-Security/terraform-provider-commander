// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEpmPolicyDataSource_MetadataAndNew(t *testing.T) {
	t.Parallel()
	d := epmpolicy.NewEpmPolicyDataSource().(*epmpolicy.EpmPolicyDataSource)
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "commander"}, &resp)
	if resp.TypeName != "commander_epm_policy" {
		t.Fatalf("got %s", resp.TypeName)
	}
}

func TestEpmPolicyDataSource_Configure(t *testing.T) {
	t.Parallel()
	d := epmpolicy.NewEpmPolicyDataSource().(*epmpolicy.EpmPolicyDataSource)
	var resp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: 1}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Configure_Success(t *testing.T) {
	t.Parallel()
	server := startMockServer(&helpers.CommandServer{}, func(string, int) (string, interface{}) { return "ok", nil })
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	if d == nil {
		t.Fatal("nil data source")
	}
}

func TestEpmPolicyDataSource_Schema(t *testing.T) {
	t.Parallel()
	d := epmpolicy.NewEpmPolicyDataSource().(*epmpolicy.EpmPolicyDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	if resp.Schema.Attributes == nil {
		t.Fatal("schema")
	}
}

func TestEpmPolicyDataSource_Read_Success(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", samplePolicyViewData()
		}
		return "ok", nil
	})
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	emptyState := tftypes.NewValue(objType, newConfigValues(nil, nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read: %v", resp.Diagnostics)
	}
}

func TestEpmPolicyDataSource_Read_ConfigGetError(t *testing.T) {
	t.Parallel()
	server := startMockServer(&helpers.CommandServer{}, func(string, int) (string, interface{}) { return "ok", nil })
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error on null config")
	}
}

func TestEpmPolicyDataSource_Read_EmptyPolicy(t *testing.T) {
	t.Parallel()
	server := startMockServer(&helpers.CommandServer{}, func(string, int) (string, interface{}) { return "ok", nil })
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_NoApiManager(t *testing.T) {
	t.Parallel()
	d := epmpolicy.NewEpmPolicyDataSource().(*epmpolicy.EpmPolicyDataSource)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_ExecuteError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm policy view", func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", samplePolicyViewData()
		}
		return "ok", nil
	})
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_EnterpriseDownError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "enterprise-down", nil)
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_EpmSyncDownError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer500OnSubstring(mock, "epm sync-down", nil)
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_MapPolicyViewError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", map[string]interface{}{
				"PolicyId":   "",
				"PolicyType": "Command",
				"Status":     "enforce",
			}
		}
		return "ok", nil
	})
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_UnmarshalError(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServer(mock, func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "epm policy view") {
			return "ok", map[string]interface{}{"PolicyId": true}
		}
		return "ok", nil
	})
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestEpmPolicyDataSource_Read_PolicyNotFound(t *testing.T) {
	t.Parallel()
	mock := &helpers.CommandServer{}
	server := startMockServerPolicyViewNotFound(mock)
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error when policy not found")
	}
}

func TestEpmPolicyDataSource_Read_ApiError(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"failed"}`))
	}))
	defer server.Close()
	d := newConfiguredDataSource(t, server)
	sch, objType := getSchema(t)
	configRaw := tftypes.NewValue(objType, newConfigValues("55", nil))
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	var resp datasource.ReadResponse
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}
