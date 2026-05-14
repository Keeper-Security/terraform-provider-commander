// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildRecordAddCommand(data WifiResourceModel) string {
	var extra []string

	records.AppendOptionalTextField(&extra, FlagSSID, data.SSID)
	records.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)
	records.AppendOptionalJSONStringAdd(&extra, FlagEncryption, data.Encryption)
	records.AppendOptionalJSONBoolAdd(&extra, FlagIsSSIDHidden, data.IsSSIDHidden)

	custom := records.NormalizeCustomFromPlan(data.Custom)
	return records.BuildRecordAdd(data.Folder, data.Title.ValueString(), records.RecordTypeWifiCredentials, extra, custom, data.Notes)
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
	return !records.CustomFieldsEqual(plan.Custom, state.Custom)
}

func buildRecordUpdateCommand(recordUID string, plan, state WifiResourceModel) string {
	var extra []string

	records.AppendChangedStringField(&extra, FlagSSID, plan.SSID, state.SSID)
	records.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)
	records.AppendChangedJSONStringField(&extra, FlagEncryption, plan.Encryption, state.Encryption)
	records.AppendChangedJSONBoolField(&extra, FlagIsSSIDHidden, plan.IsSSIDHidden, state.IsSSIDHidden)

	customPlan := records.NormalizeCustomFromPlan(plan.Custom)
	customState := records.NormalizeCustomFromPlan(state.Custom)
	return records.BuildRecordUpdate(recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// mapVaultRecordToModel hydrates the Terraform state from a `get <uid> --format json` payload.
func mapVaultRecordToModel(ctx context.Context, rec *utils.VaultRecordGetResponse, stateFolder types.String, m *WifiResourceModel) diag.Diagnostics {
	records.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.SSID = records.FirstStringField(rec.Fields, records.FieldTypeText, "SSID")
	m.Password = records.FirstStringFieldAnyLabel(rec.Fields, records.FieldTypePassword)
	m.Encryption = records.FirstStringFieldAnyLabel(rec.Fields, records.FieldTypeWifiEncryption)
	m.IsSSIDHidden = records.FirstBoolFieldAnyLabel(rec.Fields, records.FieldTypeIsSSIDHidden)

	// Parse share record permissions from the API response.
	shareMap, diags := records.ParseSharePermissionsFromResponse(ctx, rec.UserPermissions)
	if diags.HasError() {
		return diags
	}
	m.Share = shareMap

	return diags
}
