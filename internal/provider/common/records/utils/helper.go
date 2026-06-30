// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FormatFieldAssignment returns one token: key='value' or key=” for clear.
// useJSON wraps value as $JSON:... inside quotes when true.
func FormatFieldAssignment(key, value string, useJSON bool) string {
	if useJSON && strings.TrimSpace(value) != "" {
		payload := "$JSON:" + value
		return fmt.Sprintf("%s=%s", key, utils.QuoteShellSingle(payload))
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Sprintf("%s=%s", key, utils.QuoteShellSingle(""))
	}
	return fmt.Sprintf("%s=%s", key, utils.QuoteShellSingle(value))
}

// BuildRecordAdd builds a record-add command (standard fields as extraParts, then custom, then notes).
func BuildRecordAdd(folder types.String, title, recordType string, extraParts []string, custom []CustomFieldModel, notes types.String) string {
	parts := []string{utils.CmdRecordAdd}
	if !folder.IsNull() && !folder.IsUnknown() && strings.TrimSpace(folder.ValueString()) != "" {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagFolder, utils.QuoteShellSingle(strings.TrimSpace(folder.ValueString()))))
	}
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, utils.QuoteShellSingle(strings.TrimSpace(title))))
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, recordType))
	parts = append(parts, extraParts...)
	AppendCustomFieldsAdd(&parts, custom)
	if !notes.IsNull() && !notes.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, utils.QuoteShellSingle(notes.ValueString())))
	}
	return strings.Join(parts, " ")
}

// BuildRecordUpdate builds record-update with changed title, extra field parts, custom, notes.
func BuildRecordUpdate(recordUID string, titlePlan, titleState types.String, extraParts []string, customPlan, customState []CustomFieldModel, notesPlan, notesState types.String) string {
	parts := []string{
		utils.CmdRecordUpdate,
		fmt.Sprintf("%s %s", utils.FlagRecord, utils.QuoteShellSingle(strings.TrimSpace(recordUID))),
	}
	if !titlePlan.Equal(titleState) && !titlePlan.IsNull() && !titlePlan.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, utils.QuoteShellSingle(titlePlan.ValueString())))
	}
	parts = append(parts, extraParts...)
	AppendCustomFieldsUpdate(&parts, customPlan, customState)
	if !notesPlan.Equal(notesState) && !notesPlan.IsNull() && !notesPlan.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, utils.QuoteShellSingle(notesPlan.ValueString())))
	}
	return strings.Join(parts, " ")
}

// AppendChangedStringField appends f.text.label or typed unlabeled field when plan != state.
func AppendChangedStringField(parts *[]string, key string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, FormatFieldAssignment(key, "", false))
		return
	}
	*parts = append(*parts, FormatFieldAssignment(key, plan.ValueString(), false))
}

// AppendChangedJSONField appends key=$JSON:... when JSON payload changed.
func AppendChangedJSONField(parts *[]string, key string, planJSON, stateJSON string, changed bool) {
	if !changed {
		return
	}
	if strings.TrimSpace(planJSON) == "" {
		*parts = append(*parts, FormatFieldAssignment(key, "", false))
		return
	}
	*parts = append(*parts, FormatFieldAssignment(key, planJSON, true))
}

// AppendChangedRefUpdate compares plan vs state UIDs and emits addressRef/cardRef JSON array when changed.
func AppendChangedRefUpdate(parts *[]string, key string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, FormatFieldAssignment(key, "", false))
		return
	}
	u := strings.TrimSpace(plan.ValueString())
	arr, _ := json.Marshal([]string{u})
	*parts = append(*parts, FormatFieldAssignment(key, string(arr), true))
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
// pass the complete Commander key e.g. "f.text.company"). Values are shell-quoted via FormatFieldAssignment.
func AppendOptionalTextField(parts *[]string, fullFlag string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fullFlag, strings.TrimSpace(v.ValueString()), false))
}

// AppendOptionalStringAdd adds f.text.label=value for record-add when set.
func AppendOptionalStringAdd(parts *[]string, label string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(FText(label), strings.TrimSpace(v.ValueString()), false))
}

// AppendOptionalScalarAdd adds key=value for simple scalar field types (login, password, url, email).
func AppendOptionalScalarAdd(parts *[]string, fieldKey string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, strings.TrimSpace(v.ValueString()), false))
}

// AppendOptionalJSONAdd adds name=$JSON:... style.
func AppendOptionalJSONAdd(parts *[]string, fieldKey, jsonPayload string) {
	if strings.TrimSpace(jsonPayload) == "" {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, jsonPayload, true))
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
	*parts = append(*parts, FormatFieldAssignment(fieldKey, strconv.FormatInt(ms, 10), false))
}

// AppendChangedEpochDateField emits fieldKey epoch when plan != state.
func AppendChangedEpochDateField(parts *[]string, fieldKey string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, FormatFieldAssignment(fieldKey, "", false))
		return
	}
	ms, err := DateStringToEpochMillisOrZero(plan.ValueString())
	if err != nil || ms == 0 {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, strconv.FormatInt(ms, 10), false))
}

// AppendOptionalBoolAdd adds fieldKey=true|false when the value is set.
// Use for typed boolean fields such as `f.isSSIDHidden` where the API expects a literal bool string.
func AppendOptionalBoolAdd(parts *[]string, fieldKey string, v types.Bool) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	val := "false"
	if v.ValueBool() {
		val = "true"
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, val, false))
}

// AppendChangedBoolField emits fieldKey=true|false when plan != state.
// Sends an empty value when the plan is null so the field is cleared on the server.
func AppendChangedBoolField(parts *[]string, fieldKey string, plan, state types.Bool) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() {
		*parts = append(*parts, FormatFieldAssignment(fieldKey, "", false))
		return
	}
	val := "false"
	if plan.ValueBool() {
		val = "true"
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, val, false))
}

// AppendOptionalJSONStringAdd marshals a string value as JSON and emits
// fieldKey='$JSON:"value"' when set. Used for Keeper fields that require
// JSON-wrapped string values (e.g. `wifiEncryption='$JSON:"wpa"'`).
func AppendOptionalJSONStringAdd(parts *[]string, fieldKey string, v types.String) {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return
	}
	b, err := json.Marshal(strings.TrimSpace(v.ValueString()))
	if err != nil {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, string(b), true))
}

// AppendChangedJSONStringField emits fieldKey='$JSON:"value"' when plan != state.
// Clears the field with an empty value when the plan is null.
func AppendChangedJSONStringField(parts *[]string, fieldKey string, plan, state types.String) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() || strings.TrimSpace(plan.ValueString()) == "" {
		*parts = append(*parts, FormatFieldAssignment(fieldKey, "", false))
		return
	}
	b, err := json.Marshal(strings.TrimSpace(plan.ValueString()))
	if err != nil {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, string(b), true))
}

// AppendOptionalJSONBoolAdd marshals a bool value as JSON and emits
// fieldKey='$JSON:true' / fieldKey='$JSON:false' when set. Used for Keeper
// fields that require JSON-wrapped boolean values (e.g. `isSSIDHidden='$JSON:false'`).
func AppendOptionalJSONBoolAdd(parts *[]string, fieldKey string, v types.Bool) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	b, err := json.Marshal(v.ValueBool())
	if err != nil {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, string(b), true))
}

// AppendChangedJSONBoolField emits fieldKey='$JSON:true|false' when plan != state.
// Clears the field with an empty value when the plan is null.
func AppendChangedJSONBoolField(parts *[]string, fieldKey string, plan, state types.Bool) {
	if plan.Equal(state) || plan.IsUnknown() {
		return
	}
	if plan.IsNull() {
		*parts = append(*parts, FormatFieldAssignment(fieldKey, "", false))
		return
	}
	b, err := json.Marshal(plan.ValueBool())
	if err != nil {
		return
	}
	*parts = append(*parts, FormatFieldAssignment(fieldKey, string(b), true))
}

// AppendOptionalRefAdd adds addressRef / cardRef first UID.
func AppendOptionalRefAdd(parts *[]string, fieldKey string, uid types.String) {
	if uid.IsNull() || uid.IsUnknown() || strings.TrimSpace(uid.ValueString()) == "" {
		return
	}
	u := strings.TrimSpace(uid.ValueString())
	arr, _ := json.Marshal([]string{u})
	*parts = append(*parts, FormatFieldAssignment(fieldKey, string(arr), true))
}

// FetchVaultRecord runs `get <recordUID> --format json --include-dag`.
func FetchVaultRecord(ctx context.Context, apiManager *api.ApiManager, recordUID string) (*api.RequestResultResponse, error) {
	command := fmt.Sprintf("%s '%s' %s %s", utils.CmdGet, recordUID, utils.FlagFormatJSON, utils.FlagIncludeDag)
	return apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryFetchVaultRecordFailed)
}

// MoveRecordFromSourceToDestination moves a record when plan and state folder paths differ.
func MoveRecordFromSourceToDestination(ctx context.Context, apiManager *api.ApiManager, recordUID string, planFolderData string, stateFolderData string) error {
	if planFolderData == stateFolderData {
		return nil
	}

	dest := planFolderData
	if dest == "" {
		dest = "/"
	}

	command := fmt.Sprintf("%s '%s' '%s' %s", utils.CmdMv, recordUID, dest, utils.FlagForce)
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryMoveRecordFailed)
	return err
}

// MapBaseVaultRecord maps record_uid, title, notes, and folder from API onto base.
func MapBaseVaultRecord(rec *utils.VaultRecordGetResponse, stateFolderLocation types.String, base *BaseVaultRecordModel) {
	if rec == nil || base == nil {
		return
	}
	if strings.TrimSpace(rec.RecordUID) != "" {
		base.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}
	base.Title = utils.StringOrNull(rec.Title)
	base.Notes = utils.StringOrNull(rec.Notes)
	base.FolderLocation = utils.ExtractFolderValue(rec.FolderLocation, stateFolderLocation)
}
