// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddCommand_IncludesSaasConfigurationFields(t *testing.T) {
	t.Parallel()

	data := commonrecordsaasconfiguration.SaasConfigurationModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("SaaS Config"),
			Notes: types.StringValue("rotation config"),
		},
		Custom: []commonrecordsutils.CustomFieldModel{
			{
				Type:  types.StringValue("text"),
				Label: types.StringValue("AppName"),
				Value: types.StringValue("MyApp"),
			},
		},
	}

	cmd := commonrecordsaasconfiguration.BuildAddCommand(utils.CmdRecordAdd, data)

	for _, want := range []string{
		"record-add",
		"--record-type saasConfiguration",
		"SaaS Config",
		"'c.text.AppName'='MyApp'",
		"MyApp",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("command %q missing %q", cmd, want)
		}
	}
}

func TestMapVaultRecordGetResponseToSaasConfigurationModel(t *testing.T) {
	t.Parallel()

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-saas-1",
		Title:     "SaaS Config",
		Type:      "saasConfiguration",
		Custom: []utils.VaultRecordFieldResponse{
			{Type: "text", Label: "AppName", Value: json.RawMessage(`["MyApp"]`)},
		},
	}

	var state commonrecordsaasconfiguration.SaasConfigurationModel
	commonrecordsaasconfiguration.MapVaultRecordGetResponseToSaasConfigurationModel(rec, types.StringNull(), &state)

	if state.Title.ValueString() != "SaaS Config" {
		t.Fatalf("title = %q", state.Title.ValueString())
	}
	if len(state.Custom) != 1 || state.Custom[0].Value.ValueString() != "MyApp" {
		t.Fatalf("custom = %+v", state.Custom)
	}
}

func TestUpdateHasMutations_CustomChanged(t *testing.T) {
	t.Parallel()

	plan := commonrecordsaasconfiguration.SaasConfigurationModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("SaaS Config"),
		},
		Custom: []commonrecordsutils.CustomFieldModel{
			{Type: types.StringValue("text"), Label: types.StringValue("AppName"), Value: types.StringValue("NewApp")},
		},
	}
	state := commonrecordsaasconfiguration.SaasConfigurationModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("SaaS Config"),
		},
		Custom: []commonrecordsutils.CustomFieldModel{
			{Type: types.StringValue("text"), Label: types.StringValue("AppName"), Value: types.StringValue("OldApp")},
		},
	}

	if !commonrecordsaasconfiguration.UpdateHasMutations(plan, state) {
		t.Fatal("expected custom change to be detected")
	}
}

func TestLinkToGateway_SendsExpectedCommand(t *testing.T) {
	t.Parallel()

	mock := &helpers.CommandServer{}
	var gotCmd string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		gotCmd = cmd
		return "ok", nil
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	if err := commonrecordsaasconfiguration.LinkToGateway(context.Background(), apiManager, "gw-1", "config-uid-1", "record-uid-1"); err != nil {
		t.Fatalf("LinkToGateway failed: %v", err)
	}

	for _, want := range []string{
		"pam action saas update",
		"--gateway 'gw-1'",
		"--configuration-uid 'config-uid-1'",
		"--config-record-uid 'record-uid-1'",
	} {
		if !strings.Contains(gotCmd, want) {
			t.Errorf("command %q missing %q", gotCmd, want)
		}
	}
}

func TestLinkToGateway_PropagatesError(t *testing.T) {
	t.Parallel()

	server := helpers.StartCommandServerWithResultHook(&helpers.CommandServer{}, nil, func(cmd string, _ int) (int, []byte) {
		return 500, []byte(`{"message":"gateway link failed"}`)
	})
	defer server.Close()

	apiManager := &api.ApiManager{
		ServiceModeUrl:    server.URL,
		ServiceModeApiKey: "test-key",
		HttpClient:        server.Client(),
	}

	if err := commonrecordsaasconfiguration.LinkToGateway(context.Background(), apiManager, "gw-1", "config-uid-1", "record-uid-1"); err == nil {
		t.Fatal("expected error to propagate")
	}
}
