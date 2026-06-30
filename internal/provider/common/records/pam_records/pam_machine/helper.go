// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add style command for a PAM Machine record.
// The CLI command (`record-add` or `nsf-record-add`) is provided by the caller.
func BuildAddCommand(cmd string, data PamMachineResourceModel) string {
	parts := []string{cmd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamMachine))
	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))

	commonpamrecords.AppendHostnameOrIPField(&parts, data.HostnameOrIP)

	commonpamrecords.AppendOptionalTextField(&parts, FlagOperatingSystem, data.OperatingSystem)
	commonpamrecords.AppendOptionalTextField(&parts, FlagInstanceName, data.InstanceName)
	commonpamrecords.AppendOptionalTextField(&parts, FlagInstanceId, data.InstanceId)
	commonpamrecords.AppendOptionalTextField(&parts, FlagProviderGroup, data.ProviderGroup)
	commonpamrecords.AppendOptionalTextField(&parts, FlagProviderRegion, data.ProviderRegion)

	if !data.FolderLocation.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagFolder, data.FolderLocation.ValueString()))
	}

	if !data.Notes.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}

// BuildUpdateCommand builds a record-update style command for a PAM Machine
// record, including only the fields that changed between plan and state.
func BuildUpdateCommand(cmd, recordUID string, plan, state PamMachineResourceModel) string {
	parts := []string{
		cmd,
		fmt.Sprintf("%s '%s'", utils.FlagRecord, recordUID),
	}

	if !plan.Title.Equal(state.Title) {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, plan.Title.ValueString()))
	}

	if !commonpamrecords.HostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) {
		commonpamrecords.AppendHostnameOrIPField(&parts, plan.HostnameOrIP)
	}

	commonpamrecords.AppendChangedTextField(&parts, FlagOperatingSystem, plan.OperatingSystem, state.OperatingSystem)
	commonpamrecords.AppendChangedTextField(&parts, FlagInstanceName, plan.InstanceName, state.InstanceName)
	commonpamrecords.AppendChangedTextField(&parts, FlagInstanceId, plan.InstanceId, state.InstanceId)
	commonpamrecords.AppendChangedTextField(&parts, FlagProviderGroup, plan.ProviderGroup, state.ProviderGroup)
	commonpamrecords.AppendChangedTextField(&parts, FlagProviderRegion, plan.ProviderRegion, state.ProviderRegion)

	if !plan.Notes.Equal(state.Notes) && !plan.Notes.IsUnknown() {
		if plan.Notes.IsNull() {
			parts = append(parts, fmt.Sprintf("%s ''", utils.FlagNotes))
		} else {
			parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, plan.Notes.ValueString()))
		}
	}

	return strings.Join(parts, " ")
}

// RecordUpdateHasMutations returns true when at least one record-level field
// differs between plan and state (pam_settings is checked separately).
func RecordUpdateHasMutations(plan, state PamMachineResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!commonpamrecords.HostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) ||
		!plan.OperatingSystem.Equal(state.OperatingSystem) ||
		!plan.InstanceName.Equal(state.InstanceName) ||
		!plan.InstanceId.Equal(state.InstanceId) ||
		!plan.ProviderGroup.Equal(state.ProviderGroup) ||
		!plan.ProviderRegion.Equal(state.ProviderRegion) ||
		!plan.Notes.Equal(state.Notes)
}

// MapVaultRecordGetResponseToPamMachineModel fills state from `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPamMachineModel(rec *utils.VaultRecordGetResponse, state *PamMachineResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	commonrecordsutils.MapBaseVaultRecord(rec, state.FolderLocation, &state.BaseVaultRecordModel)

	// pamHostname field
	state.HostnameOrIP = ExtractPamHostnameFieldValue(rec.Fields)

	// Text fields extracted by label
	state.OperatingSystem = utils.StringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "operatingSystem"))
	state.InstanceName = utils.StringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "instanceName"))
	state.InstanceId = utils.StringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "instanceId"))
	state.ProviderGroup = utils.StringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "providerGroup"))
	state.ProviderRegion = utils.StringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "providerRegion"))

	state.PamSettings = commonpamrecords.ExtractMachineDirectoryPamSettingsFromResponse(rec, state.PamSettings)

	return diags
}

// ExtractPamHostnameFieldValue extracts the pamHostname field value from the fields array.
func ExtractPamHostnameFieldValue(fields []utils.VaultRecordFieldResponse) *commonpamrecords.HostnameOrIPModel {
	for i := range fields {
		f := &fields[i]
		if f.Type != "pamHostname" {
			continue
		}
		var vals []utils.PamRemoteBrowserHostnameFieldResponse
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			return nil
		}
		if len(vals) > 0 {
			model := &commonpamrecords.HostnameOrIPModel{
				HostName: utils.StringOrNull(vals[0].HostName),
			}
			portStr := strings.TrimSpace(vals[0].AdministrativePort)
			if portStr != "" {
				if parsed, err := strconv.ParseInt(portStr, 10, 32); err == nil {
					model.AdministrativePort = types.Int32Value(int32(parsed))
				} else {
					model.AdministrativePort = types.Int32Null()
				}
			} else {
				model.AdministrativePort = types.Int32Null()
			}
			return model
		}
	}
	return nil
}
