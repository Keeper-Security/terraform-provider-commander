// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for an address record.
func BuildAddCommand(cmd string, data AddressModel) string {
	var extra []string

	if data.Address != nil && !data.Address.IsNull() {
		if j, err := data.Address.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagAddress, j)
		}
	}

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeAddress, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable address fields.
func UpdateHasMutations(plan, state AddressModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) {
		return true
	}
	if !commonrecordsutils.AddressEqual(plan.Address, state.Address) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed address fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state AddressModel) string {
	var extra []string

	planJSON, planErr := plan.Address.ToJSON()
	stateJSON, stateErr := state.Address.ToJSON()
	changed := planJSON != stateJSON || planErr != stateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagAddress, planJSON, stateJSON, changed)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToAddressModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToAddressModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *AddressModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Address = commonrecordsutils.AddressFromFields(rec.Fields, "")
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
