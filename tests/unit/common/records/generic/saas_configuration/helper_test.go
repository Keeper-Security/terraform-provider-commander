// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"encoding/json"
	"strings"
	"testing"

	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
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
