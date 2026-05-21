// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
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

// MapVaultRecordGetResponseToPamDirectoryModel fills state from `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPamDirectoryModel(rec *utils.VaultRecordGetResponse, state *PamDirectoryResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if strings.TrimSpace(rec.RecordUID) != "" {
		state.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}
	state.Title = setStringOrNull(rec.Title)
	state.Notes = setStringOrNull(rec.Notes)

	state.Folder = commonpamrecords.ExtractFolderValue(rec.Folder, state.Folder)

	state.HostnameOrIP = ExtractPamHostnameFieldValue(rec.Fields)
	state.UseSSL = extractCheckboxFieldValue(rec.Fields, "useSSL")
	state.DomainName = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "domainName"))
	state.AlternativeIPs = extractMultilineAsSet(rec.Fields, "alternativeIPs")
	state.DirectoryId = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "directoryId"))
	state.DirectoryType = extractDirectoryTypeFieldValue(rec.Fields)
	state.UserMatch = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "userMatch"))
	state.ProviderGroup = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "providerGroup"))
	state.ProviderRegion = setStringOrNull(commonpamrecords.ExtractFirstTextFieldValue(rec.Fields, "providerRegion"))

	state.PamSettings = commonpamrecords.ExtractPamSettingsFromResponse(rec, state.PamSettings)

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
				HostName: setStringOrNull(vals[0].HostName),
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

// extractCheckboxFieldValue extracts a boolean from a checkbox-type field.
func extractCheckboxFieldValue(fields []utils.VaultRecordFieldResponse, label string) types.Bool {
	for i := range fields {
		f := &fields[i]
		if f.Type != "checkbox" || f.Label != label {
			continue
		}
		var vals []bool
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			return types.BoolNull()
		}
		if len(vals) > 0 {
			return types.BoolValue(vals[0])
		}
	}
	return types.BoolNull()
}

// extractDirectoryTypeFieldValue extracts the directoryType field value.
func extractDirectoryTypeFieldValue(fields []utils.VaultRecordFieldResponse) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != "directoryType" {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			return types.StringNull()
		}
		if len(vals) > 0 {
			return setStringOrNull(vals[0])
		}
	}
	return types.StringNull()
}

// extractMultilineAsSet reads a multiline field (newline-separated values) and
// returns them as a Terraform Set of strings. Returns null if not found.
func extractMultilineAsSet(fields []utils.VaultRecordFieldResponse, label string) types.Set {
	for i := range fields {
		f := &fields[i]
		if f.Type != "multiline" || f.Label != label {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			return types.SetNull(types.StringType)
		}
		if len(vals) == 0 || strings.TrimSpace(vals[0]) == "" {
			return types.SetNull(types.StringType)
		}
		lines := strings.Split(vals[0], "\n")
		cleaned := filterNonEmptyLines(lines)
		if len(cleaned) == 0 {
			return types.SetNull(types.StringType)
		}
		sv, diags := types.SetValueFrom(context.Background(), types.StringType, cleaned)
		if diags.HasError() {
			return types.SetNull(types.StringType)
		}
		return sv
	}
	return types.SetNull(types.StringType)
}

func filterNonEmptyLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
