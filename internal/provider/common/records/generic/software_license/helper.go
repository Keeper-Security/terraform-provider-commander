// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a softwareLicense record.
func BuildAddCommand(cmd string, data SoftwareLicenseModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagLicenseNumber, data.SoftwareLicenseKey)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagExpirationDate, data.ExpirationDate)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagDateActive, data.DateActive)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(
		cmd,
		commonrecordsutils.RecordTypeSoftwareLicense,
		data.Title.ValueString(),
		data.FolderLocation,
		extra,
		custom,
		data.Notes,
	)
}

// UpdateHasMutations reports whether plan differs from state on updatable software license fields.
func UpdateHasMutations(plan, state SoftwareLicenseModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.SoftwareLicenseKey.Equal(state.SoftwareLicenseKey) ||
		!plan.ExpirationDate.Equal(state.ExpirationDate) ||
		!plan.DateActive.Equal(state.DateActive) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed software license fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state SoftwareLicenseModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagLicenseNumber, plan.SoftwareLicenseKey, state.SoftwareLicenseKey)
	commonrecordsutils.AppendChangedStringField(&extra, FlagExpirationDate, plan.ExpirationDate, state.ExpirationDate)
	commonrecordsutils.AppendChangedStringField(&extra, FlagDateActive, plan.DateActive, state.DateActive)

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

// MapVaultRecordGetResponseToSoftwareLicenseModel fills state from a `get <uid> --format json` payload.
// ExpirationDate and DateActive use DateOrEpochMillisField* so NSF YYYY-MM-DD strings and
// classic epoch-ms values both map to YYYY-MM-DD without changing EpochMillisField itself.
func MapVaultRecordGetResponseToSoftwareLicenseModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *SoftwareLicenseModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.SoftwareLicenseKey = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeLicenseNumber)
	m.ExpirationDate = commonrecordsutils.DateOrEpochMillisFieldUnlabeled(rec.Fields, commonrecordsutils.FieldTypeExpirationDate)
	m.DateActive = commonrecordsutils.DateOrEpochMillisField(rec.Fields, commonrecordsutils.FieldTypeDate, DateActiveLabel)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
