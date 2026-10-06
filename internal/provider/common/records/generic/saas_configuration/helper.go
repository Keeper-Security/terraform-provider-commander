// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a saasConfiguration record.
func BuildAddCommand(cmd string, data SaasConfigurationModel) string {
	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(
		cmd,
		commonrecordsutils.RecordTypeSaasConfiguration,
		data.Title.ValueString(),
		data.FolderLocation,
		nil,
		custom,
		data.Notes,
	)
}

// UpdateHasMutations reports whether plan differs from state on updatable fields.
func UpdateHasMutations(plan, state SaasConfigurationModel) bool {
	if !plan.Title.Equal(state.Title) || !plan.Notes.Equal(state.Notes) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state SaasConfigurationModel) string {
	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(
		cmd,
		recordUID,
		plan.Title,
		state.Title,
		nil,
		customPlan,
		customState,
		plan.Notes,
		state.Notes,
	)
}

// MapVaultRecordGetResponseToSaasConfigurationModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToSaasConfigurationModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *SaasConfigurationModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}

// ErrOpListSaasPluginsForGateway formats ErrOpListSaasPlugins with the specific gateway name/UID.
func ErrOpListSaasPluginsForGateway(gateway string) string {
	return fmt.Sprintf("%s %q", ErrOpListSaasPlugins, gateway)
}

// LinkToGateway used to connect a just-created SaaS configuration record to the gateway and PAM configuration that will use it for rotation.
// So it performs all configuration complete.
func LinkToGateway(ctx context.Context, apiManager *api.ApiManager, gateway, configuration, recordUID string) error {
	command := fmt.Sprintf(
		"%s --gateway %s --configuration-uid %s --config-record-uid %s",
		utils.CmdPamActionSaaSUpdate,
		utils.QuoteShellSingle(gateway),
		utils.QuoteShellSingle(configuration),
		utils.QuoteShellSingle(recordUID),
	)
	_, err := apiManager.ExecuteCommand(ctx, command, ErrOpLinkSaasConfigToGateway)
	return err
}
