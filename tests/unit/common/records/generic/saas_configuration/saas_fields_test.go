// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// reportedPluginInfoData is the exact `data` object shape from the real
// `pam action saas config --gateway xx --plugin REST --info --format json` response. Fields
// callers don't use (name/type/author/email/installed/readme/summary/
// allows_remote_management) are included to confirm they're tolerated.
var reportedPluginInfoData = map[string]interface{}{
	"name":                     "REST",
	"type":                     "builtin",
	"author":                   "Keeper Security",
	"email":                    "pam@keepersecuirty.com",
	"installed":                true,
	"readme":                   "",
	"summary":                  "Generic REST implementation that calls your custom webservice to rotate a user's password.",
	"allows_remote_management": false,
	"required_fields": []map[string]interface{}{
		{"label": "REST Url", "type": "text", "required": true, "description": "URL endpoint."},
		{"label": "REST Token", "type": "secret", "required": true, "description": "A Bearer token."},
	},
	"optional_fields": []map[string]interface{}{
		{"label": "REST Method", "type": "text", "required": false, "description": "HTTP method. Either 'POST' or 'PUT'"},
	},
}

func newApiManagerForPluginInfo(t *testing.T, gateway, plugin string, data interface{}) *api.ApiManager {
	t.Helper()
	mock := &helpers.CommandServer{}
	var gotCommand string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		gotCommand = cmd
		return "", data
	})
	t.Cleanup(server.Close)
	t.Cleanup(func() {
		if !strings.Contains(gotCommand, "pam action saas config") || !strings.Contains(gotCommand, "--info") || !strings.Contains(gotCommand, "--format json") {
			t.Errorf("expected command to fetch plugin info as JSON, got: %s", gotCommand)
		}
		if !strings.Contains(gotCommand, "--gateway '"+gateway+"'") {
			t.Errorf("expected command to include --gateway '%s', got: %s", gateway, gotCommand)
		}
		if !strings.Contains(gotCommand, "--plugin '"+plugin+"'") {
			t.Errorf("expected command to include --plugin '%s', got: %s", plugin, gotCommand)
		}
	})

	return &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
}

func TestFetchSaasFieldSpecs_ParsesReportedResponseShape(t *testing.T) {
	t.Parallel()

	apiManager := newApiManagerForPluginInfo(t, "gw-1", "REST", reportedPluginInfoData)

	specs, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(context.Background(), apiManager, "gw-1", "REST")
	if err != nil {
		t.Fatalf("FetchSaasFieldSpecs failed: %v", err)
	}

	want := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url", Type: "text"},
		{Label: "REST Token", Type: "secret"},
		{Label: "REST Method", Type: "text"},
	}
	if len(specs) != len(want) {
		t.Fatalf("expected %d field specs, got %d: %+v", len(want), len(specs), specs)
	}
	for i, w := range want {
		if specs[i] != w {
			t.Errorf("spec %d = %+v, want %+v", i, specs[i], w)
		}
	}
}

func TestFetchSaasFieldSpecs_NoFieldsReturnsEmpty(t *testing.T) {
	t.Parallel()

	apiManager := newApiManagerForPluginInfo(t, "gw-1", "NoFieldsPlugin", map[string]interface{}{
		"name":    "NoFieldsPlugin",
		"type":    "builtin",
		"summary": "A plugin with no configurable fields.",
	})

	specs, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(context.Background(), apiManager, "gw-1", "NoFieldsPlugin")
	if err != nil {
		t.Fatalf("FetchSaasFieldSpecs failed: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("expected no field specs, got %+v", specs)
	}
}

func TestFetchSaasFieldSpecs_NullDataReturnsNoSpecs(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(_ string, _ int) (string, interface{}) {
		return "", nil
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	specs, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(context.Background(), apiManager, "gw-1", "REST")
	if err != nil {
		t.Fatalf("FetchSaasFieldSpecs failed: %v", err)
	}
	if specs != nil {
		t.Errorf("expected no specs for a null data response, got: %+v", specs)
	}
}

func TestFetchSaasFieldSpecs_MalformedDataPropagatesError(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(_ string, _ int) (string, interface{}) {
		// A string instead of the expected plugin info object.
		return "", "not-an-object"
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	if _, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(context.Background(), apiManager, "gw-1", "REST"); err == nil {
		t.Fatal("expected an error when the response data isn't shaped as a plugin info object")
	}
}

func textField(label, value string) commonrecordsutils.CustomFieldModel {
	return commonrecordsutils.CustomFieldModel{
		Type:  types.StringValue("text"),
		Label: types.StringValue(label),
		Value: types.StringValue(value),
	}
}

func TestValidateSaasFields_AcceptsAllRequiredAndValidOptional(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
		{Label: "REST Token"},
		{Label: "REST Method"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		textField("REST Url", "https://example.com"),
		textField("REST Token", "s3cr3t"),
		textField("REST Method", "POST"),
	}

	if err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs); err != nil {
		t.Errorf("expected valid, got error: %v", err)
	}
}

func TestValidateSaasFields_RejectsOptionalFieldOmittedEntirely(t *testing.T) {
	t.Parallel()

	// By product decision, every field the plugin reports (Required or Optional) is
	// mandatory - omitting one reported "Optional" is now an error too.
	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
		{Label: "REST Method"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		textField("REST Url", "https://example.com"),
	}

	err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs)
	if err == nil {
		t.Fatal("expected error when an optional-per-the-plugin field is omitted entirely")
	}
	if !strings.Contains(err.Error(), "REST Method") {
		t.Errorf("expected error to name the missing field, got: %v", err)
	}
}

func TestValidateSaasFields_RejectsMissingRequiredField(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
		{Label: "REST Token"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		textField("REST Url", "https://example.com"),
	}

	err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs)
	if err == nil {
		t.Fatal("expected error when a required field is missing")
	}
	if !strings.Contains(err.Error(), "REST Token") {
		t.Errorf("expected error to name the missing field, got: %v", err)
	}
}

func TestValidateSaasFields_ErrorNamesPluginAndShowsConfigExample(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
	}

	err := commonrecordsaasconfiguration.ValidateSaasFields("REST", nil, specs)
	if err == nil {
		t.Fatal("expected an error")
	}

	msg := err.Error()
	for _, want := range []string{
		`SaaS Type (Plugin) "REST"`,
		`"REST Url" is missing`,
		`type = "text"`,
		`label = "REST Url"`,
		`value = "<your value>"`,
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message missing %q, got: %s", want, msg)
		}
	}
}

func TestValidateSaasFields_RejectsEmptyRequiredValue(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		textField("REST Url", "   "),
	}

	if err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs); err == nil {
		t.Fatal("expected error when a required field's value is blank/whitespace-only")
	}
}

func TestValidateSaasFields_RejectsEmptyOptionalValueWhenPresent(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
		{Label: "REST Method"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		textField("REST Url", "https://example.com"),
		textField("REST Method", ""),
	}

	err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs)
	if err == nil {
		t.Fatal("expected error when an optional field is present but blank")
	}
	if !strings.Contains(err.Error(), "REST Method") {
		t.Errorf("expected error to name the blank optional field, got: %v", err)
	}
}

func TestValidateSaasFields_RejectsWrongFieldType(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Token"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		{Type: types.StringValue("secret"), Label: types.StringValue("REST Token"), Value: types.StringValue("s3cr3t")},
	}

	err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs)
	if err == nil {
		t.Fatal("expected error when a plugin field is provided with a non-text type")
	}
	if !strings.Contains(err.Error(), "REST Token") {
		t.Errorf("expected error to name the mistyped field, got: %v", err)
	}
}

func TestValidateSaasFields_IgnoresUnrelatedCustomFields(t *testing.T) {
	t.Parallel()

	specs := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url"},
	}
	custom := []commonrecordsutils.CustomFieldModel{
		textField("SaaS Type", "REST"),
		textField("Active", "true"),
		textField("REST Url", "https://example.com"),
	}

	if err := commonrecordsaasconfiguration.ValidateSaasFields("REST", custom, specs); err != nil {
		t.Errorf("expected the universal SaaS Type/Active fields to be ignored, got error: %v", err)
	}
}
