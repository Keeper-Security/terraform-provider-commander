// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// reportedPluginListMessage is the exact `message` array shape from the real
// `pam action saas config --gateway xx --list` response, including the
// "Available SaaS Plugins" header line that must be skipped.
var reportedPluginListMessage = []string{
	"Available SaaS Plugins",
	"* AWS Cognito (Catalog) - Change a users password in AWS Cognito.",
	"* Cisco APIC (Catalog) - Change a user password in Cisco APIC.",
	"* Snowflake (Builtin) - For Snowflake, rotate the password for a user.",
	"* My Custom Plugin (Custom) - Some description.",
}

func TestFetchSaasPluginList_ParsesReportedResponseShape(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	// The message array contains embedded double quotes, which the simpler
	// StartCommandServer helper can't represent (it naively string-wraps
	// message), so build the raw response body via the result hook instead.
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, _ int) (int, []byte) {
		if !strings.Contains(cmd, "pam action saas config") || !strings.Contains(cmd, "--gateway") || !strings.Contains(cmd, "--list") {
			return 0, nil
		}
		if !strings.Contains(cmd, "'gw-1'") {
			t.Errorf("expected command to include --gateway 'gw-1', got: %s", cmd)
		}
		body, _ := json.Marshal(map[string]interface{}{
			"data":    nil,
			"status":  "success",
			"message": reportedPluginListMessage,
			"error":   "",
		})
		return http.StatusOK, body
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	entries, err := commonrecordsaasconfiguration.FetchSaasPluginList(context.Background(), apiManager, "gw-1")
	if err != nil {
		t.Fatalf("FetchSaasPluginList failed: %v", err)
	}

	// The header line ("Available SaaS Plugins") must not become an entry.
	if len(entries) != 4 {
		t.Fatalf("expected 4 plugin entries (header skipped), got %d: %+v", len(entries), entries)
	}

	want := []commonrecordsaasconfiguration.SaasPluginEntry{
		{Name: "AWS Cognito", IntegrationType: "Catalog", Description: "Change a users password in AWS Cognito."},
		{Name: "Cisco APIC", IntegrationType: "Catalog", Description: "Change a user password in Cisco APIC."},
		{Name: "Snowflake", IntegrationType: "Builtin", Description: "For Snowflake, rotate the password for a user."},
		{Name: "My Custom Plugin", IntegrationType: "Custom", Description: "Some description."},
	}
	for i, w := range want {
		if entries[i] != w {
			t.Errorf("entry %d = %+v, want %+v", i, entries[i], w)
		}
	}
}

func TestFetchSaasPluginList_NonListMessagePropagatesError(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		return "some plain non-list message", nil
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	// Note: this specific failure path (an unparseable message body) doesn't itself embed the
	// gateway name in the returned error - the gateway still reaches the practitioner via the
	// diagnostic Summary at the call site (see TestErrOpListSaasPluginsForGateway).
	if _, err := commonrecordsaasconfiguration.FetchSaasPluginList(context.Background(), apiManager, "gw-1"); err == nil {
		t.Fatal("expected an error when the response message isn't list-shaped")
	}
}

func TestErrOpListSaasPluginsForGateway(t *testing.T) {
	t.Parallel()

	got := commonrecordsaasconfiguration.ErrOpListSaasPluginsForGateway("gw-1")
	if !strings.Contains(got, commonrecordsaasconfiguration.ErrOpListSaasPlugins) || !strings.Contains(got, "gw-1") {
		t.Errorf("expected summary to contain both the base message and the gateway, got: %q", got)
	}
}

func TestFetchSaasPluginList_GatewayOfflineReturnsClearError(t *testing.T) {
	t.Parallel()

	// Exact reported response shape: still "status":"success" with no real plugin
	// entries, just an offline notice followed by the usual header line.
	offlineMessage := []string{
		"This Gateway currently is not online.",
		"Did not get router response.",
		"Available SaaS Plugins",
	}

	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, nil, func(cmd string, _ int) (int, []byte) {
		if !strings.Contains(cmd, "pam action saas config") || !strings.Contains(cmd, "--list") {
			return 0, nil
		}
		body, _ := json.Marshal(map[string]interface{}{
			"data":    nil,
			"status":  "success",
			"message": offlineMessage,
			"error":   "",
		})
		return http.StatusOK, body
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	entries, err := commonrecordsaasconfiguration.FetchSaasPluginList(context.Background(), apiManager, "gw-1")
	if err == nil {
		t.Fatal("expected an error when the gateway is offline, got none")
	}
	if !errors.Is(err, commonrecordsaasconfiguration.ErrGatewayOffline) {
		t.Errorf("expected errors.Is(err, ErrGatewayOffline) to be true, got: %v", err)
	}
	if entries != nil {
		t.Errorf("expected no entries when the gateway is offline, got: %+v", entries)
	}
}

func TestValidateSaasType(t *testing.T) {
	t.Parallel()

	entries := []commonrecordsaasconfiguration.SaasPluginEntry{
		{Name: "Snowflake", IntegrationType: "Builtin", Description: "desc"},
		{Name: "Okta", IntegrationType: "Builtin", Description: "desc"},
	}

	if err := commonrecordsaasconfiguration.ValidateSaasType("Snowflake", entries); err != nil {
		t.Errorf("expected Snowflake to be valid, got error: %v", err)
	}

	// Empty saasType is not this function's concern (the required-field validator handles that).
	if err := commonrecordsaasconfiguration.ValidateSaasType("", entries); err != nil {
		t.Errorf("expected empty saasType to be ignored, got error: %v", err)
	}

	err := commonrecordsaasconfiguration.ValidateSaasType("NotAPlugin", entries)
	if err == nil {
		t.Fatal("expected an error for an unsupported SaaS type")
	}
	if !strings.Contains(err.Error(), "NotAPlugin") || !strings.Contains(err.Error(), "Snowflake") || !strings.Contains(err.Error(), "Okta") {
		t.Errorf("expected error to name the invalid value and list supported values, got: %v", err)
	}
}

func TestSaasTypeFromCustom(t *testing.T) {
	t.Parallel()

	custom := []commonrecordsutils.CustomFieldModel{
		{Type: types.StringValue("text"), Label: types.StringValue("AppName"), Value: types.StringValue("MyApp")},
		{Type: types.StringValue("text"), Label: types.StringValue("SaaS Type"), Value: types.StringValue("Snowflake")},
	}

	if got := commonrecordsaasconfiguration.SaasTypeFromCustom(custom); got != "Snowflake" {
		t.Errorf("SaasTypeFromCustom = %q, want Snowflake", got)
	}

	if got := commonrecordsaasconfiguration.SaasTypeFromCustom(nil); got != "" {
		t.Errorf("SaasTypeFromCustom(nil) = %q, want empty", got)
	}
}
