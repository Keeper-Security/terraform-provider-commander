// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add style command for a PAM Directory record.
// The CLI command (`record-add` or `nsf-record-add`) is provided by the caller.
func BuildAddCommand(cmd string, data PamDirectoryResourceModel) string {
	parts := []string{cmd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamDirectory))
	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))

	commonpamrecords.AppendHostnameOrIPField(&parts, data.HostnameOrIP)
	commonpamrecords.AppendOptionalCheckboxField(&parts, FlagUseSSL, data.UseSSL)
	commonpamrecords.AppendOptionalTextField(&parts, FlagDomainName, data.DomainName)
	appendAlternativeIPsField(&parts, data.AlternativeIPs)
	commonpamrecords.AppendOptionalTextField(&parts, FlagDirectoryId, data.DirectoryId)
	appendOptionalDirectoryTypeField(&parts, data.DirectoryType)
	commonpamrecords.AppendOptionalTextField(&parts, FlagUserMatch, data.UserMatch)
	commonpamrecords.AppendOptionalTextField(&parts, FlagProviderGroup, data.ProviderGroup)
	commonpamrecords.AppendOptionalTextField(&parts, FlagProviderRegion, data.ProviderRegion)

	if !data.Folder.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagFolder, data.Folder.ValueString()))
	}

	if !data.Notes.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}

// BuildUpdateCommand builds a record-update style command for a PAM Directory
// record, including only the fields that changed between plan and state.
func BuildUpdateCommand(cmd, recordUID string, plan, state PamDirectoryResourceModel) string {
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

	commonpamrecords.AppendChangedCheckboxField(&parts, FlagUseSSL, plan.UseSSL, state.UseSSL)
	commonpamrecords.AppendChangedTextField(&parts, FlagDomainName, plan.DomainName, state.DomainName)

	if !plan.AlternativeIPs.Equal(state.AlternativeIPs) {
		appendAlternativeIPsField(&parts, plan.AlternativeIPs)
	}

	commonpamrecords.AppendChangedTextField(&parts, FlagDirectoryId, plan.DirectoryId, state.DirectoryId)

	if !plan.DirectoryType.Equal(state.DirectoryType) {
		appendOptionalDirectoryTypeField(&parts, plan.DirectoryType)
	}

	commonpamrecords.AppendChangedTextField(&parts, FlagUserMatch, plan.UserMatch, state.UserMatch)
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
func RecordUpdateHasMutations(plan, state PamDirectoryResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!commonpamrecords.HostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) ||
		!plan.UseSSL.Equal(state.UseSSL) ||
		!plan.DomainName.Equal(state.DomainName) ||
		!plan.AlternativeIPs.Equal(state.AlternativeIPs) ||
		!plan.DirectoryId.Equal(state.DirectoryId) ||
		!plan.DirectoryType.Equal(state.DirectoryType) ||
		!plan.UserMatch.Equal(state.UserMatch) ||
		!plan.ProviderGroup.Equal(state.ProviderGroup) ||
		!plan.ProviderRegion.Equal(state.ProviderRegion) ||
		!plan.Notes.Equal(state.Notes)
}

func appendOptionalDirectoryTypeField(parts *[]string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	*parts = append(*parts, fmt.Sprintf("'%s=%s'", FlagDirectoryType, v.ValueString()))
}

// appendAlternativeIPsField joins the set elements with newlines and writes
// them as a single multiline field value.
func appendAlternativeIPsField(parts *[]string, s types.Set) {
	if s.IsNull() || s.IsUnknown() || len(s.Elements()) == 0 {
		return
	}
	var ips []string
	diags := s.ElementsAs(context.Background(), &ips, false)
	if diags.HasError() {
		return
	}
	joined := strings.Join(ips, "\\n")
	*parts = append(*parts, fmt.Sprintf("'%s=%s'", FlagAlternativeIPs, joined))
}

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

	state.Folder = utils.ExtractFolderValue(rec.FolderLocation, state.Folder)

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
