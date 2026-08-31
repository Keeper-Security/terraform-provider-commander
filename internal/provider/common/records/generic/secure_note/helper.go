// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for an encryptedNotes record.
func BuildAddCommand(cmd string, data SecureNoteModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagSecuredNote, data.SecuredNote)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagDate, data.Date)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(
		cmd,
		commonrecordsutils.RecordTypeEncryptedNotes,
		data.Title.ValueString(),
		data.FolderLocation,
		extra,
		custom,
		data.Notes,
	)
}

// UpdateHasMutations reports whether plan differs from state on updatable secure note fields.
func UpdateHasMutations(plan, state SecureNoteModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.SecuredNote.Equal(state.SecuredNote) ||
		!plan.Date.Equal(state.Date) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed secure note fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state SecureNoteModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagSecuredNote, plan.SecuredNote, state.SecuredNote)
	commonrecordsutils.AppendChangedStringField(&extra, FlagDate, plan.Date, state.Date)

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

// MapVaultRecordGetResponseToSecureNoteModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToSecureNoteModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *SecureNoteModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.SecuredNote = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeNote)
	m.Date = commonrecordsutils.EpochMillisFieldUnlabeled(rec.Fields, commonrecordsutils.FieldTypeDate)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
