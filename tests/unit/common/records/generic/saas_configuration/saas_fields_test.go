// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// reportedPluginInfoMessage is the exact `message` array shape from the real
// `pam action saas config --gateway xx --plugin REST --info` response.
var reportedPluginInfoMessage = []string{
	"REST",
	"Type: builtin",
	"Author: Keeper Security (pam@keepersecuirty.com)",
	"Summary: Generic REST implementation that calls your custom webservice to rotate a user's password.",
	"Fields",
	"* Required: REST Url - URL endpoint.",
	"* Required: REST Token - A Bearer token.",
	"* Optional: REST Method - HTTP method. Either 'POST' or 'PUT'",
}

func newApiManagerForPluginInfo(t *testing.T, gateway, plugin string, message []string) *api.ApiManager {
	t.Helper()
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, _ int) (int, []byte) {
		if !strings.Contains(cmd, "pam action saas config") || !strings.Contains(cmd, "--info") {
			return 0, nil
		}
		if !strings.Contains(cmd, "--gateway '"+gateway+"'") {
			t.Errorf("expected command to include --gateway '%s', got: %s", gateway, cmd)
		}
		if !strings.Contains(cmd, "--plugin '"+plugin+"'") {
			t.Errorf("expected command to include --plugin '%s', got: %s", plugin, cmd)
		}
		body, _ := json.Marshal(map[string]interface{}{
			"data":    nil,
			"status":  "success",
			"message": message,
			"error":   "",
		})
		return http.StatusOK, body
	})
	t.Cleanup(server.Close)

	return &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}
}

func TestFetchSaasFieldSpecs_ParsesReportedResponseShape(t *testing.T) {
	t.Parallel()

	apiManager := newApiManagerForPluginInfo(t, "gw-1", "REST", reportedPluginInfoMessage)

	specs, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(context.Background(), apiManager, "gw-1", "REST")
	if err != nil {
		t.Fatalf("FetchSaasFieldSpecs failed: %v", err)
	}

	want := []commonrecordsaasconfiguration.SaasFieldSpec{
		{Label: "REST Url", Required: true, Description: "URL endpoint."},
		{Label: "REST Token", Required: true, Description: "A Bearer token."},
		{Label: "REST Method", Required: false, Description: "HTTP method. Either 'POST' or 'PUT'"},
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

func TestFetchSaasFieldSpecs_NoFieldsSectionReturnsEmpty(t *testing.T) {
	t.Parallel()

	apiManager := newApiManagerForPluginInfo(t, "gw-1", "NoFieldsPlugin", []string{
		"NoFieldsPlugin",
		"Type: builtin",
		"Summary: A plugin with no configurable fields.",
	})

	specs, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(context.Background(), apiManager, "gw-1", "NoFieldsPlugin")
	if err != nil {
		t.Fatalf("FetchSaasFieldSpecs failed: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("expected no field specs, got %+v", specs)
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
		{Label: "REST Url", Required: true},
		{Label: "REST Token", Required: true},
		{Label: "REST Method", Required: false},
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
		{Label: "REST Url", Required: true},
		{Label: "REST Method", Required: false},
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
		{Label: "REST Url", Required: true},
		{Label: "REST Token", Required: true},
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
		{Label: "REST Url", Required: true},
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
		{Label: "REST Url", Required: true},
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
		{Label: "REST Url", Required: true},
		{Label: "REST Method", Required: false},
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
		{Label: "REST Token", Required: true},
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
		{Label: "REST Url", Required: true},
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
