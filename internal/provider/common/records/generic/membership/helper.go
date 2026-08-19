// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a membership record.
func BuildAddCommand(cmd string, data MembershipModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagAccountNumber, data.AccountNumber)
	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeMembership, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable membership fields.
func UpdateHasMutations(plan, state MembershipModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.AccountNumber.Equal(state.AccountNumber) ||
		!plan.Password.Equal(state.Password) {
		return true
	}
	if !commonrecordsutils.NameEqual(plan.Name, state.Name) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed membership fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state MembershipModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagAccountNumber, plan.AccountNumber, state.AccountNumber)

	namePlanJSON, namePlanErr := plan.Name.ToJSON()
	nameStateJSON, nameStateErr := state.Name.ToJSON()
	nameChanged := namePlanJSON != nameStateJSON || namePlanErr != nameStateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagName, namePlanJSON, nameStateJSON, nameChanged)

	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToMembershipModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToMembershipModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *MembershipModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.AccountNumber = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeAccountNumber)
	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
