// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a birthCertificate record.
func BuildAddCommand(cmd string, data BirthCertificateModel) string {
	var extra []string

	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagBirthDate, data.BirthDate)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeBirthCertificate, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable birthCertificate fields.
func UpdateHasMutations(plan, state BirthCertificateModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.BirthDate.Equal(state.BirthDate) {
		return true
	}
	if !commonrecordsutils.NameEqual(plan.Name, state.Name) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed birthCertificate fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state BirthCertificateModel) string {
	var extra []string

	namePlanJSON, namePlanErr := plan.Name.ToJSON()
	nameStateJSON, nameStateErr := state.Name.ToJSON()
	nameChanged := namePlanJSON != nameStateJSON || namePlanErr != nameStateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagName, namePlanJSON, nameStateJSON, nameChanged)

	commonrecordsutils.AppendChangedStringField(&extra, FlagBirthDate, plan.BirthDate, state.BirthDate)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToBirthCertificateModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToBirthCertificateModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *BirthCertificateModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.BirthDate = commonrecordsutils.EpochMillisFieldUnlabeled(rec.Fields, commonrecordsutils.FieldTypeBirthDate)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
