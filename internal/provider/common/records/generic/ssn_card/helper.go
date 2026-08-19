// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for an ssnCard (Identity Card) record.
func BuildAddCommand(cmd string, data SsnCardModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagAccountNumber, data.AccountNumber)
	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeSsnCard, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable ssnCard fields.
func UpdateHasMutations(plan, state SsnCardModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.AccountNumber.Equal(state.AccountNumber) {
		return true
	}
	if !commonrecordsutils.NameEqual(plan.Name, state.Name) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed ssnCard fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state SsnCardModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagAccountNumber, plan.AccountNumber, state.AccountNumber)

	namePlanJSON, namePlanErr := plan.Name.ToJSON()
	nameStateJSON, nameStateErr := state.Name.ToJSON()
	nameChanged := namePlanJSON != nameStateJSON || namePlanErr != nameStateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagName, namePlanJSON, nameStateJSON, nameChanged)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToSsnCardModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToSsnCardModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *SsnCardModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.AccountNumber = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeAccountNumber)
	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
