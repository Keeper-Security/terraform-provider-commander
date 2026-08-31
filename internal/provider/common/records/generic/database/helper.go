// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a databaseCredentials record.
func BuildAddCommand(cmd string, data DatabaseModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagLogin, data.Login)
	commonrecordsutils.AppendOptionalHostFlatAdd(&extra, FlagHost, data.Hostname, data.Port)
	commonrecordsutils.AppendOptionalTextField(&extra, FlagTextType, data.Type)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(
		cmd,
		commonrecordsutils.RecordTypeDatabaseCredentials,
		data.Title.ValueString(),
		data.FolderLocation,
		extra,
		custom,
		data.Notes,
	)
}

// UpdateHasMutations reports whether plan differs from state on updatable database fields.
func UpdateHasMutations(plan, state DatabaseModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Hostname.Equal(state.Hostname) ||
		!plan.Port.Equal(state.Port) ||
		!plan.Type.Equal(state.Type) ||
		!plan.Password.Equal(state.Password) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed database fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state DatabaseModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagLogin, plan.Login, state.Login)
	commonrecordsutils.AppendChangedHostFlatUpdate(&extra, FlagHost, plan.Hostname, plan.Port, state.Hostname, state.Port)
	commonrecordsutils.AppendChangedStringField(&extra, FlagTextType, plan.Type, state.Type)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(
		cmd,
		recordUID,
		plan.Title,
		state.Title,
		extra,
		customPlan,
		customState,
		plan.Notes,
		state.Notes,
	)
}

// MapVaultRecordGetResponseToDatabaseModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToDatabaseModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *DatabaseModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Login = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeLogin)
	m.Hostname, m.Port = commonrecordsutils.FlatHostFromFields(rec.Fields)
	m.Type = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, TextTypeLabel)
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
