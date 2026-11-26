// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add style command for a PAM Cloud record.
// The CLI command (`record-add` or `nsf-record-add`) is provided by the caller.
func BuildAddCommand(cmd string, data PamCloudResourceModel) string {
	parts := []string{cmd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamCloud))
	parts = append(parts, utils.FlagArg(utils.FlagTitle, data.Title.ValueString()))

	commonpamrecords.AppendOptionalTextField(&parts, FlagURL, data.URL)
	commonpamrecords.AppendOptionalTextField(&parts, FlagAccountIdentifier, data.AccountIdentifier)

	if !data.FolderLocation.IsNull() && !data.FolderLocation.IsUnknown() && strings.TrimSpace(data.FolderLocation.ValueString()) != "" {
		parts = append(parts, utils.FlagArg(utils.FlagFolder, strings.TrimSpace(data.FolderLocation.ValueString())))
	}

	if !data.Notes.IsNull() {
		parts = append(parts, utils.FlagArg(utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}

// BuildUpdateCommand builds a record-update style command for a PAM Cloud
// record, including only the fields that changed between plan and state.
func BuildUpdateCommand(cmd, recordUID string, plan, state PamCloudResourceModel) string {
	parts := []string{
		cmd,
		utils.FlagArg(utils.FlagRecord, recordUID),
	}

	if !plan.Title.Equal(state.Title) {
		parts = append(parts, utils.FlagArg(utils.FlagTitle, plan.Title.ValueString()))
	}

	commonpamrecords.AppendChangedTextField(&parts, FlagURL, plan.URL, state.URL)
	commonpamrecords.AppendChangedTextField(&parts, FlagAccountIdentifier, plan.AccountIdentifier, state.AccountIdentifier)

	if !plan.Notes.Equal(state.Notes) && !plan.Notes.IsUnknown() {
		if plan.Notes.IsNull() {
			parts = append(parts, fmt.Sprintf("%s ''", utils.FlagNotes))
		} else {
			parts = append(parts, utils.FlagArg(utils.FlagNotes, plan.Notes.ValueString()))
		}
	}

	return strings.Join(parts, " ")
}

// RecordUpdateHasMutations returns true when at least one record-level field
// differs between plan and state.
func RecordUpdateHasMutations(plan, state PamCloudResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!plan.URL.Equal(state.URL) ||
		!plan.AccountIdentifier.Equal(state.AccountIdentifier) ||
		!plan.Notes.Equal(state.Notes)
}

// MapVaultRecordGetResponseToPamCloudModel fills state from `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPamCloudModel(rec *utils.VaultRecordGetResponse, state *PamCloudResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	commonrecordsutils.MapBaseVaultRecord(rec, state.FolderLocation, &state.BaseVaultRecordModel)

	state.URL = extractFirstFieldValueByType(rec.Fields, FieldTypeURL)
	state.AccountIdentifier = extractFirstFieldValueByType(rec.Fields, FieldTypeAccountIdentifier)

	return diags
}

// extractFirstFieldValueByType extracts the first string value of the field
// with the given type, or null when the field is absent or empty.
func extractFirstFieldValueByType(fields []utils.VaultRecordFieldResponse, fieldType string) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			return types.StringNull()
		}
		if len(vals) > 0 {
			return utils.StringOrNull(vals[0])
		}
	}
	return types.StringNull()
}
