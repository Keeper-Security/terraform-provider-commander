// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// QuoteShellSingle wraps s for use as a single-quoted shell argument.
func QuoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

// formatFieldAssignment returns one token: key='value' or key=” for clear.
// useJSON wraps value as $JSON:... inside quotes when true.
func formatFieldAssignment(key, value string, useJSON bool) string {
	if useJSON && strings.TrimSpace(value) != "" {
		payload := "$JSON:" + value
		return fmt.Sprintf("%s=%s", key, QuoteShellSingle(payload))
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Sprintf("%s=%s", key, QuoteShellSingle(""))
	}
	return fmt.Sprintf("%s=%s", key, QuoteShellSingle(value))
}

// MapBaseVaultRecord maps record_uid, title, notes, folder, custom from API onto base.
func MapBaseVaultRecord(rec *utils.VaultRecordGetResponse, stateFolder types.String, base *BaseVaultRecordModel) {
	if rec == nil || base == nil {
		return
	}
	if strings.TrimSpace(rec.RecordUID) != "" {
		base.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}
	base.Title = stringOrNull(rec.Title)
	base.Notes = stringOrNull(rec.Notes)
	base.Folder = ExtractFolderValue(rec.Folder, stateFolder)
	base.Custom = ParseCustomFields(rec.Custom)
}

// BuildRecordAdd builds a record-add command (standard fields as extraParts, then custom, then notes).
func BuildRecordAdd(folder types.String, title, recordType string, extraParts []string, custom []CustomFieldModel, notes types.String) string {
	parts := []string{utils.CmdRecordAdd}
	if !folder.IsNull() && !folder.IsUnknown() && strings.TrimSpace(folder.ValueString()) != "" {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagFolder, QuoteShellSingle(strings.TrimSpace(folder.ValueString()))))
	}
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, QuoteShellSingle(strings.TrimSpace(title))))
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, recordType))
	parts = append(parts, extraParts...)
	AppendCustomFieldsAdd(&parts, custom)
	if !notes.IsNull() && !notes.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, QuoteShellSingle(notes.ValueString())))
	}
	return strings.Join(parts, " ")
}

// BuildRecordUpdate builds record-update with changed title, extra field parts, custom, notes.
func BuildRecordUpdate(recordUID string, titlePlan, titleState types.String, extraParts []string, customPlan, customState []CustomFieldModel, notesPlan, notesState types.String) string {
	parts := []string{
		utils.CmdRecordUpdate,
		fmt.Sprintf("%s %s", utils.FlagRecord, QuoteShellSingle(strings.TrimSpace(recordUID))),
	}
	if !titlePlan.Equal(titleState) && !titlePlan.IsNull() && !titlePlan.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, QuoteShellSingle(titlePlan.ValueString())))
	}
	parts = append(parts, extraParts...)
	AppendCustomFieldsUpdate(&parts, customPlan, customState)
	if !notesPlan.Equal(notesState) && !notesPlan.IsNull() && !notesPlan.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, QuoteShellSingle(notesPlan.ValueString())))
	}
	return strings.Join(parts, " ")
}

// AppendChangedStringField appends f.text.label or typed unlabeled field when plan != state.
func AppendChangedStringField(parts *[]string, key string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, formatFieldAssignment(key, "", false))
		return
	}
	*parts = append(*parts, formatFieldAssignment(key, plan.ValueString(), false))
}

// AppendChangedJSONField appends key=$JSON:... when JSON payload changed.
func AppendChangedJSONField(parts *[]string, key string, planJSON, stateJSON string, changed bool) {
	if !changed {
		return
	}
	if strings.TrimSpace(planJSON) == "" {
		*parts = append(*parts, formatFieldAssignment(key, "", false))
		return
	}
	*parts = append(*parts, formatFieldAssignment(key, planJSON, true))
}

// AppendChangedRefUpdate compares plan vs state UIDs and emits addressRef/cardRef JSON array when changed.
func AppendChangedRefUpdate(parts *[]string, key string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, formatFieldAssignment(key, "", false))
		return
	}
	u := strings.TrimSpace(plan.ValueString())
	arr, _ := json.Marshal([]string{u})
	*parts = append(*parts, formatFieldAssignment(key, string(arr), true))
}

// FolderPathStrings returns trimmed folder path/UID from Terraform strings (empty if unset).
func FolderPathStrings(planFolder, stateFolder types.String) (plan, state string) {
	if !planFolder.IsNull() && !planFolder.IsUnknown() {
		plan = strings.TrimSpace(planFolder.ValueString())
	}
	if !stateFolder.IsNull() && !stateFolder.IsUnknown() {
		state = strings.TrimSpace(stateFolder.ValueString())
	}
	return plan, state
}

// FText returns Commander key f.text.<label>.
func FText(label string) string {
	return "f.text." + label
}

// AppendOptionalTextField adds fullFlag=value when set (same pattern as pam_records.AppendOptionalTextField:
// pass the complete Commander key e.g. "f.text.company"). Values are shell-quoted via formatFieldAssignment.
func AppendOptionalTextField(parts *[]string, fullFlag string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	*parts = append(*parts, formatFieldAssignment(fullFlag, strings.TrimSpace(v.ValueString()), false))
}

// AppendOptionalStringAdd adds f.text.label=value for record-add when set.
func AppendOptionalStringAdd(parts *[]string, label string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	*parts = append(*parts, formatFieldAssignment(FText(label), strings.TrimSpace(v.ValueString()), false))
}

// AppendOptionalScalarAdd adds key=value for simple scalar field types (login, password, url, email).
func AppendOptionalScalarAdd(parts *[]string, fieldKey string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	*parts = append(*parts, formatFieldAssignment(fieldKey, strings.TrimSpace(v.ValueString()), false))
}

// AppendOptionalJSONAdd adds name=$JSON:... style.
func AppendOptionalJSONAdd(parts *[]string, fieldKey, jsonPayload string) {
	if strings.TrimSpace(jsonPayload) == "" {
		return
	}
	*parts = append(*parts, formatFieldAssignment(fieldKey, jsonPayload, true))
}

// AppendOptionalEpochDateAdd adds date field as epoch number string.
func AppendOptionalEpochDateAdd(parts *[]string, fieldKey string, dateStr types.String) {
	if dateStr.IsNull() || dateStr.IsUnknown() || strings.TrimSpace(dateStr.ValueString()) == "" {
		return
	}
	ms, err := DateStringToEpochMillisOrZero(dateStr.ValueString())
	if err != nil {
		return
	}
	if ms == 0 {
		return
	}
	*parts = append(*parts, formatFieldAssignment(fieldKey, strconv.FormatInt(ms, 10), false))
}

// AppendChangedEpochDateField emits fieldKey epoch when plan != state.
func AppendChangedEpochDateField(parts *[]string, fieldKey string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, formatFieldAssignment(fieldKey, "", false))
		return
	}
	ms, err := DateStringToEpochMillisOrZero(plan.ValueString())
	if err != nil || ms == 0 {
		return
	}
	*parts = append(*parts, formatFieldAssignment(fieldKey, strconv.FormatInt(ms, 10), false))
}

// AppendOptionalRefAdd adds addressRef / cardRef first UID.
func AppendOptionalRefAdd(parts *[]string, fieldKey string, uid types.String) {
	if uid.IsNull() || uid.IsUnknown() || strings.TrimSpace(uid.ValueString()) == "" {
		return
	}
	u := strings.TrimSpace(uid.ValueString())
	arr, _ := json.Marshal([]string{u})
	*parts = append(*parts, formatFieldAssignment(fieldKey, string(arr), true))
}
