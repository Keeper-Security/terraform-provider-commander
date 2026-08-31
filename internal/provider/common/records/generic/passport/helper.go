// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package passport

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a passport record.
func BuildAddCommand(cmd string, data PassportModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagAccountNumber, data.AccountNumber)
	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagBirthDate, data.BirthDate)
	commonrecordsutils.AppendOptionalTextField(&extra, FlagAddressRef, data.AddressRef)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagExpirationDate, data.ExpirationDate)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagDateIssued, data.DateIssued)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypePassport, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable passport fields.
func UpdateHasMutations(plan, state PassportModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.AccountNumber.Equal(state.AccountNumber) ||
		!plan.BirthDate.Equal(state.BirthDate) ||
		!plan.AddressRef.Equal(state.AddressRef) ||
		!plan.ExpirationDate.Equal(state.ExpirationDate) ||
		!plan.DateIssued.Equal(state.DateIssued) ||
		!plan.Password.Equal(state.Password) {
		return true
	}
	if !commonrecordsutils.NameEqual(plan.Name, state.Name) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed passport fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state PassportModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagAccountNumber, plan.AccountNumber, state.AccountNumber)

	namePlanJSON, namePlanErr := plan.Name.ToJSON()
	nameStateJSON, nameStateErr := state.Name.ToJSON()
	nameChanged := namePlanJSON != nameStateJSON || namePlanErr != nameStateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagName, namePlanJSON, nameStateJSON, nameChanged)

	commonrecordsutils.AppendChangedStringField(&extra, FlagBirthDate, plan.BirthDate, state.BirthDate)
	commonrecordsutils.AppendChangedStringField(&extra, FlagAddressRef, plan.AddressRef, state.AddressRef)
	commonrecordsutils.AppendChangedStringField(&extra, FlagExpirationDate, plan.ExpirationDate, state.ExpirationDate)
	commonrecordsutils.AppendChangedStringField(&extra, FlagDateIssued, plan.DateIssued, state.DateIssued)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToPassportModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPassportModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *PassportModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.AccountNumber = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeAccountNumber)
	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.BirthDate = commonrecordsutils.EpochMillisFieldUnlabeled(rec.Fields, commonrecordsutils.FieldTypeBirthDate)
	m.AddressRef = commonrecordsutils.FirstRefUID(rec.Fields, commonrecordsutils.FieldTypeAddressRef, "")
	m.ExpirationDate = commonrecordsutils.EpochMillisFieldUnlabeled(rec.Fields, commonrecordsutils.FieldTypeExpirationDate)
	m.DateIssued = commonrecordsutils.EpochMillisFieldUnlabeled(rec.Fields, commonrecordsutils.FieldTypeDate)
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
