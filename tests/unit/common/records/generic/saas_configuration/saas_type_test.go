// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// reportedPluginListData is the exact `data` array shape from the real
// `pam action saas config --gateway xx --list --format json` response. Fields callers don't
// use (author/email/installed/readme/used_by_count/using) are included to confirm they're
// tolerated, not just the fields FetchSaasPluginList reads.
var reportedPluginListData = []map[string]interface{}{
	{"author": "Keeper Security", "email": "pam@keepersecurity.com", "installed": false, "name": "AWS Cognito", "readme": "https://example.com/aws_cognito", "summary": "Change a users password in AWS Cognito.", "type": "catalog", "used_by_count": 0, "using": false},
	{"author": "Keeper Security", "email": "pam@keepersecurity.com", "installed": false, "name": "Cisco APIC", "readme": "https://example.com/cisco_apic", "summary": "Change a user password in Cisco APIC.", "type": "catalog", "used_by_count": 0, "using": false},
	{"author": "Keeper Security", "email": "pam@keepersecuirty.com", "installed": true, "name": "Snowflake", "readme": "", "summary": "For Snowflake, rotate the password for a user.", "type": "builtin", "used_by_count": 0, "using": false},
	{"author": "Keeper Security", "email": "pam@keepersecurity.com", "installed": true, "name": "My Custom Plugin", "readme": "", "summary": "Some description.", "type": "custom", "used_by_count": 0, "using": false},
}

func TestFetchSaasPluginList_ParsesReportedResponseShape(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	var gotCommand string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		gotCommand = cmd
		return "", reportedPluginListData
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

	if !strings.Contains(gotCommand, "pam action saas config") || !strings.Contains(gotCommand, "--list") || !strings.Contains(gotCommand, "--format json") {
		t.Errorf("expected command to list plugins as JSON, got: %s", gotCommand)
	}
	if !strings.Contains(gotCommand, "'gw-1'") {
		t.Errorf("expected command to include --gateway 'gw-1', got: %s", gotCommand)
	}

	if len(entries) != 4 {
		t.Fatalf("expected 4 plugin entries, got %d: %+v", len(entries), entries)
	}

	want := []commonrecordsaasconfiguration.SaasPluginEntry{
		{Name: "AWS Cognito", IntegrationType: "catalog", Description: "Change a users password in AWS Cognito."},
		{Name: "Cisco APIC", IntegrationType: "catalog", Description: "Change a user password in Cisco APIC."},
		{Name: "Snowflake", IntegrationType: "builtin", Description: "For Snowflake, rotate the password for a user."},
		{Name: "My Custom Plugin", IntegrationType: "custom", Description: "Some description."},
	}
	for i, w := range want {
		if entries[i] != w {
			t.Errorf("entry %d = %+v, want %+v", i, entries[i], w)
		}
	}
}

func TestFetchSaasPluginList_NullDataReturnsNoEntries(t *testing.T) {
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

	entries, err := commonrecordsaasconfiguration.FetchSaasPluginList(context.Background(), apiManager, "gw-1")
	if err != nil {
		t.Fatalf("FetchSaasPluginList failed: %v", err)
	}
	if entries != nil {
		t.Errorf("expected no entries for a null data response, got: %+v", entries)
	}
}

func TestFetchSaasPluginList_MalformedDataPropagatesError(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(_ string, _ int) (string, interface{}) {
		// A string instead of the expected array of plugin objects.
		return "", "not-an-array"
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	// Note: this specific failure path (a malformed data payload) doesn't itself embed the
	// gateway name in the returned error - the gateway still reaches the practitioner via the
	// diagnostic Summary at the call site (see TestErrOpListSaasPluginsForGateway).
	if _, err := commonrecordsaasconfiguration.FetchSaasPluginList(context.Background(), apiManager, "gw-1"); err == nil {
		t.Fatal("expected an error when the response data isn't shaped as a list of plugins")
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
	// entries, just an offline notice instead of the usual list.
	const offlineMessage = "This Gateway currently is not online. Did not get router response."

	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServer(mock, func(_ string, _ int) (string, interface{}) {
		return offlineMessage, nil
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
	// Supported values must render as a SaaS Type / Plugin Type / Description table, not a
	// flat comma-separated list.
	if !strings.Contains(err.Error(), "SAAS TYPE") || !strings.Contains(err.Error(), "PLUGIN TYPE") || !strings.Contains(err.Error(), "DESCRIPTION") {
		t.Errorf("expected error to render supported values as a table with headers, got: %v", err)
	}
	if !strings.Contains(err.Error(), "desc") {
		t.Errorf("expected error table to include each plugin's description, got: %v", err)
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
