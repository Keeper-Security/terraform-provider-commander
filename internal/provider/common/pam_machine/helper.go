// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"encoding/json"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func setStringOrNull(val string) types.String {
	if strings.TrimSpace(val) == "" {
		return types.StringNull()
	}
	return types.StringValue(val)
}

// MapVaultRecordGetResponseToPamMachineModel fills state from `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPamMachineModel(rec *utils.VaultRecordGetResponse, state *PamMachineResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if strings.TrimSpace(rec.RecordUID) != "" {
		state.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}
	state.Title = setStringOrNull(rec.Title)
	state.Notes = setStringOrNull(rec.Notes)
	state.Folder = setStringOrNull(rec.Folder)

	// pamHostname field
	state.HostnameOrIP = ExtractPamHostnameFieldValue(rec.Fields)

	// Text fields extracted by label
	state.OperatingSystem = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "operatingSystem"))
	state.InstanceName = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "instanceName"))
	state.InstanceId = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "instanceId"))
	state.ProviderGroup = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "providerGroup"))
	state.ProviderRegion = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "providerRegion"))

	// TODO: Currently we are not getting --folder data from the API.

	// TODO: map pamSettings when PamSettingsModel fields are defined.
	state.PamSettings = nil

	return diags
}

// ExtractPamHostnameFieldValue extracts the pamHostname field value from the fields array.
func ExtractPamHostnameFieldValue(fields []utils.VaultRecordFieldResponse) *HostnameOrIPModel {
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
			return &HostnameOrIPModel{
				HostName: setStringOrNull(vals[0].HostName),
				Port:     setStringOrNull(vals[0].Port),
			}
		}
	}
	return nil
}
