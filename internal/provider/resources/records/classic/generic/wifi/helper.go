// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildRecordAddCommand(data WifiResourceModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalTextField(&extra, FlagSSID, data.SSID)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)
	commonrecordsutils.AppendOptionalJSONStringAdd(&extra, FlagEncryption, data.Encryption)
	commonrecordsutils.AppendOptionalJSONBoolAdd(&extra, FlagIsSSIDHidden, data.IsSSIDHidden)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(data.FolderLocation, data.Title.ValueString(), commonrecordsutils.RecordTypeWifiCredentials, extra, custom, data.Notes)
}

func updateHasMutations(plan, state WifiResourceModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.SSID.Equal(state.SSID) ||
		!plan.Password.Equal(state.Password) ||
		!plan.Encryption.Equal(state.Encryption) ||
		!plan.IsSSIDHidden.Equal(state.IsSSIDHidden) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

func buildRecordUpdateCommand(recordUID string, plan, state WifiResourceModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagSSID, plan.SSID, state.SSID)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)
	commonrecordsutils.AppendChangedJSONStringField(&extra, FlagEncryption, plan.Encryption, state.Encryption)
	commonrecordsutils.AppendChangedJSONBoolField(&extra, FlagIsSSIDHidden, plan.IsSSIDHidden, state.IsSSIDHidden)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// mapVaultRecordToModel hydrates the Terraform state from a `get <uid> --format json` payload.
func mapVaultRecordToModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *WifiResourceModel) {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.SSID = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, "SSID")
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	m.Encryption = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeWifiEncryption)
	m.IsSSIDHidden = commonrecordsutils.FirstBoolFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeIsSSIDHidden)

	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
}
