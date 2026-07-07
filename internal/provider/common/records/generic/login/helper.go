// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a login record.
func BuildAddCommand(cmd string, data LoginModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagLogin, data.Login)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagURL, data.WebsiteAddress)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeLogin, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable login fields.
func UpdateHasMutations(plan, state LoginModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Password.Equal(state.Password) ||
		!plan.WebsiteAddress.Equal(state.WebsiteAddress) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed login fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state LoginModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagLogin, plan.Login, state.Login)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)
	commonrecordsutils.AppendChangedStringField(&extra, FlagURL, plan.WebsiteAddress, state.WebsiteAddress)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToLoginModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToLoginModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *LoginModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Login = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeLogin)
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	m.WebsiteAddress = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeURL)

	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
