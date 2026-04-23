// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// quoteShellSingle wraps s for use as a single-quoted shell argument.
func quoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

// buildRecordAddPamUserCommand builds a `record-add` command string for a pamUser record.
//
// Example output:
//
//	record-add --folder 'nptw7AealHr8LC-hk8aAHg' --title 'AD service account' --record-type pamUser
//	  f.login=svc_myapp f.password='123' f.text.distinguishedName='CN=...' f.secret.privatePEMKey='abc'
//	  f.text.connectDatabase='test' f.checkbox.managed=true
func buildRecordAddPamUserCommand(data PamUserResourceModel) string {
	parts := []string{utils.CmdRecordAdd}

	// Folder (before title, matches CLI convention)
	if !data.Folder.IsNull() && !data.Folder.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagFolder, quoteShellSingle(data.Folder.ValueString())))
	}

	// Title
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, quoteShellSingle(data.Title.ValueString())))

	// Record type
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamUser))

	// Login
	if !data.Login.IsNull() && !data.Login.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldLogin, quoteShellSingle(data.Login.ValueString())))
	}

	// Password
	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldPassword, quoteShellSingle(data.Password.ValueString())))
	}

	// Distinguished name
	if !data.DistinguishedName.IsNull() && !data.DistinguishedName.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldDistinguishedName, quoteShellSingle(data.DistinguishedName.ValueString())))
	}

	// Private PEM key
	if !data.PrivatePEMKey.IsNull() && !data.PrivatePEMKey.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldPrivatePEMKey, quoteShellSingle(data.PrivatePEMKey.ValueString())))
	}

	// Connect database
	if !data.ConnectDatabase.IsNull() && !data.ConnectDatabase.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldConnectDatabase, quoteShellSingle(data.ConnectDatabase.ValueString())))
	}

	// Managed checkbox
	if !data.Managed.IsNull() && !data.Managed.IsUnknown() {
		val := "false"
		if data.Managed.ValueBool() {
			val = "true"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", FieldManaged, val))
	}

	// Notes
	if !data.Notes.IsNull() && !data.Notes.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, quoteShellSingle(data.Notes.ValueString())))
	}

	return strings.Join(parts, " ")
}

// updateHasMutations returns true when at least one field differs between plan and state.
func updateHasMutations(plan, state PamUserResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Password.Equal(state.Password) ||
		!plan.DistinguishedName.Equal(state.DistinguishedName) ||
		!plan.PrivatePEMKey.Equal(state.PrivatePEMKey) ||
		!plan.ConnectDatabase.Equal(state.ConnectDatabase) ||
		!plan.Managed.Equal(state.Managed) ||
		(!plan.Notes.Equal(state.Notes) && !plan.Notes.IsNull() && !plan.Notes.IsUnknown())
}

// buildRecordUpdatePamUserCommand builds a `record-update -r '<uid>'` command with only changed fields.
//
// Example output:
//
//	record-update --record 'nb-dp0Xm7KZWRTKt5UK83w' --title 'AD service account' f.login='svc_myapp'
func buildRecordUpdatePamUserCommand(recordUID string, plan, state PamUserResourceModel) string {
	parts := []string{
		utils.CmdRecordUpdate,
		fmt.Sprintf("%s %s", utils.FlagRecord, quoteShellSingle(recordUID)),
	}

	if !plan.Title.Equal(state.Title) {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, quoteShellSingle(plan.Title.ValueString())))
	}

	if !plan.Login.Equal(state.Login) {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldLogin, quoteShellSingle(plan.Login.ValueString())))
	}

	if !plan.Password.Equal(state.Password) {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldPassword, quoteShellSingle(plan.Password.ValueString())))
	}

	if !plan.DistinguishedName.Equal(state.DistinguishedName) {
		if plan.DistinguishedName.IsNull() || plan.DistinguishedName.IsUnknown() {
			parts = append(parts, fmt.Sprintf("%s=", FieldDistinguishedName))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", FieldDistinguishedName, quoteShellSingle(plan.DistinguishedName.ValueString())))
		}
	}

	if !plan.PrivatePEMKey.Equal(state.PrivatePEMKey) {
		if plan.PrivatePEMKey.IsNull() || plan.PrivatePEMKey.IsUnknown() {
			parts = append(parts, fmt.Sprintf("%s=", FieldPrivatePEMKey))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", FieldPrivatePEMKey, quoteShellSingle(plan.PrivatePEMKey.ValueString())))
		}
	}

	if !plan.ConnectDatabase.Equal(state.ConnectDatabase) {
		if plan.ConnectDatabase.IsNull() || plan.ConnectDatabase.IsUnknown() {
			parts = append(parts, fmt.Sprintf("%s=", FieldConnectDatabase))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", FieldConnectDatabase, quoteShellSingle(plan.ConnectDatabase.ValueString())))
		}
	}

	if !plan.Managed.Equal(state.Managed) {
		val := "false"
		if !plan.Managed.IsNull() && !plan.Managed.IsUnknown() && plan.Managed.ValueBool() {
			val = "true"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", FieldManaged, val))
	}

	if !plan.Notes.Equal(state.Notes) && !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, quoteShellSingle(plan.Notes.ValueString())))
	}

	return strings.Join(parts, " ")
}

// Vault field type constants returned by `get <uid> --format json`.
const (
	vaultFieldTypeLogin    = "login"
	vaultFieldTypePassword = "password"
	vaultFieldTypeText     = "text"
	vaultFieldTypeSecret   = "secret"
	vaultFieldTypeCheckbox = "checkbox"
)

// mapVaultRecordToState maps a `get <uid> --format json` response onto the Terraform state model.
func mapVaultRecordToState(rec *utils.VaultRecordGetResponse, state *PamUserResourceModel) {
	if strings.TrimSpace(rec.RecordUID) != "" {
		state.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}

	state.Title = stringOrNull(rec.Title)
	state.Notes = stringOrNull(rec.Notes)
	state.Folder = stringOrNull(rec.Folder)

	for i := range rec.Fields {
		f := &rec.Fields[i]
		switch f.Type {
		case vaultFieldTypeLogin:
			state.Login = firstStringValue(f.Value)
		case vaultFieldTypePassword:
			state.Password = firstStringValue(f.Value)
		case vaultFieldTypeText:
			switch f.Label {
			case "distinguishedName":
				state.DistinguishedName = firstStringValue(f.Value)
			case "connectDatabase":
				state.ConnectDatabase = firstStringValue(f.Value)
			}
		case vaultFieldTypeSecret:
			if f.Label == "privatePEMKey" {
				state.PrivatePEMKey = firstStringValue(f.Value)
			}
		case vaultFieldTypeCheckbox:
			if f.Label == "managed" {
				state.Managed = firstBoolValue(f.Value)
			}
		}
	}
}

// stringOrNull returns a types.StringValue when s is non-empty, types.StringNull otherwise.
func stringOrNull(s string) types.String {
	if strings.TrimSpace(s) == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// firstStringValue unmarshals a JSON array and returns the first non-empty string, or null.
func firstStringValue(raw json.RawMessage) types.String {
	var vals []string
	if err := json.Unmarshal(raw, &vals); err != nil {
		return types.StringNull()
	}
	if len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
		return types.StringValue(vals[0])
	}
	return types.StringNull()
}

// firstBoolValue unmarshals a JSON array and returns the first bool, or null.
func firstBoolValue(raw json.RawMessage) types.Bool {
	var vals []bool
	if err := json.Unmarshal(raw, &vals); err != nil {
		return types.BoolNull()
	}
	if len(vals) > 0 {
		return types.BoolValue(vals[0])
	}
	return types.BoolNull()
}

// buildPamRotationEditCommand builds `pam rotation edit -r <uid> --config ... --force`.
func buildPamRotationEditCommand(recordUID string, rs *PamUserRotationSettings) string {
	parts := []string{
		CmdPamRotationEdit,
		fmt.Sprintf("%s %s", FlagRecordShort, quoteShellSingle(recordUID)),
	}

	if !rs.Configuration.IsNull() && !rs.Configuration.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagConfig, quoteShellSingle(rs.Configuration.ValueString())))
	}

	if !rs.RotationProfile.IsNull() && !rs.RotationProfile.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagRotationProfile, rs.RotationProfile.ValueString()))
	}

	if !rs.IamAadConfig.IsNull() && !rs.IamAadConfig.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagIamAadConfig, quoteShellSingle(rs.IamAadConfig.ValueString())))
	}

	if !rs.Resource.IsNull() && !rs.Resource.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagResource, quoteShellSingle(rs.Resource.ValueString())))
	}

	if !rs.AdminUser.IsNull() && !rs.AdminUser.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagAdminUser, quoteShellSingle(rs.AdminUser.ValueString())))
	}

	if !rs.Enabled.IsNull() && !rs.Enabled.IsUnknown() {
		if rs.Enabled.ValueBool() {
			parts = append(parts, FlagEnable)
		} else {
			parts = append(parts, FlagDisable)
		}
	}

	if !rs.OnDemand.IsNull() && !rs.OnDemand.IsUnknown() && rs.OnDemand.ValueBool() {
		parts = append(parts, FlagOnDemand)
	} else if !rs.ScheduleConfig.IsNull() && !rs.ScheduleConfig.IsUnknown() && rs.ScheduleConfig.ValueBool() {
		parts = append(parts, FlagScheduleConfig)
	} else if !rs.ScheduleCron.IsNull() && !rs.ScheduleCron.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagScheduleCron, quoteShellSingle(rs.ScheduleCron.ValueString())))
	} else if !rs.ScheduleJSON.IsNull() && !rs.ScheduleJSON.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagScheduleJSON, quoteShellSingle(rs.ScheduleJSON.ValueString())))
	}

	if !rs.Complexity.IsNull() && !rs.Complexity.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagComplexity, quoteShellSingle(rs.Complexity.ValueString())))
	}

	parts = append(parts, FlagForce)

	return strings.Join(parts, " ")
}

// rotationSettingsNeedApply returns true when the plan has a rotation block that differs from state.
func rotationSettingsNeedApply(plan, state *PamUserRotationSettings) bool {
	if plan == nil {
		return false
	}
	if state == nil {
		return true
	}
	return !plan.RotationProfile.Equal(state.RotationProfile) ||
		!plan.Configuration.Equal(state.Configuration) ||
		!plan.IamAadConfig.Equal(state.IamAadConfig) ||
		!plan.Resource.Equal(state.Resource) ||
		!plan.AdminUser.Equal(state.AdminUser) ||
		!plan.Enabled.Equal(state.Enabled) ||
		!plan.ScheduleCron.Equal(state.ScheduleCron) ||
		!plan.ScheduleJSON.Equal(state.ScheduleJSON) ||
		!plan.OnDemand.Equal(state.OnDemand) ||
		!plan.ScheduleConfig.Equal(state.ScheduleConfig) ||
		!plan.Complexity.Equal(state.Complexity)
}

// parseRotationInfoMessage parses the message lines from `pam rotation info -r <uid>` response
// and populates the rotation settings on the state model. Fields not present in the response
// are preserved from the existing state (config-as-source-of-truth).
// hasRotationData returns true when the message lines contain a PAM Config UID,
// indicating the record has rotation configured in the vault.
func hasRotationData(messages []string) bool {
	for _, line := range messages {
		if strings.HasPrefix(strings.TrimSpace(line), "PAM Config UID:") {
			return true
		}
	}
	return false
}

func parseRotationInfoMessage(messages []string, existing *PamUserRotationSettings, state *PamUserResourceModel) {
	var rs *PamUserRotationSettings
	if existing != nil {
		rs = &PamUserRotationSettings{
			RotationProfile: existing.RotationProfile,
			Configuration:   existing.Configuration,
			IamAadConfig:    existing.IamAadConfig,
			Resource:        existing.Resource,
			AdminUser:       existing.AdminUser,
			Enabled:         existing.Enabled,
			ScheduleCron:    existing.ScheduleCron,
			ScheduleJSON:    existing.ScheduleJSON,
			OnDemand:        existing.OnDemand,
			ScheduleConfig:  existing.ScheduleConfig,
			Complexity:      existing.Complexity,
		}
	} else {
		rs = &PamUserRotationSettings{}
	}

	for _, line := range messages {
		line = strings.TrimSpace(line)
		if k, v, ok := strings.Cut(line, ": "); ok {
			switch k {
			case "PAM Config UID":
				rs.Configuration = stringOrNull(v)
			case "Admin Resource Uid":
				rs.AdminUser = stringOrNull(v)
			case "Is Rotation Disabled":
				rs.Enabled = types.BoolValue(strings.EqualFold(strings.TrimSpace(v), "False"))
			case "Schedule Type":
				if strings.Contains(strings.ToLower(v), "manual") {
					rs.OnDemand = types.BoolValue(true)
				}
			case "Schedule":
				parseScheduleValue(strings.TrimSpace(v), rs)
			case "Password Complexity Data":
				rs.Complexity = stringOrNull(v)
			}
		}
	}

	state.RotationSettings = rs
}

// parseScheduleValue parses the JSON schedule from `pam rotation info` response.
// Responses look like:
//   - {"type":"CRON","cron":"0 56 17 * * ?","tz":"Etc/UTC"}
//   - {"type":"DAILY","utcTime":"17:56","intervalCount":1}
//   - {"type":"DAILY","intervalCount":19,"time":"00:00:00","tz":"Asia/Calcutta"}
//   - {"type":"WEEKLY","utcTime":"15:44","weekday":"SUNDAY","intervalCount":1}
func parseScheduleValue(raw string, rs *PamUserRotationSettings) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}

	var sched map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &sched); err != nil {
		return
	}

	schedType, _ := sched["type"].(string)
	switch strings.ToUpper(schedType) {
	case "CRON":
		if cron, ok := sched["cron"].(string); ok {
			rs.ScheduleCron = types.StringValue(cron)
			rs.OnDemand = types.BoolNull()
			rs.ScheduleJSON = types.StringNull()
		}
	case "DAILY", "WEEKLY", "MONTHLY":
		rs.ScheduleJSON = types.StringValue(raw)
		rs.OnDemand = types.BoolNull()
		rs.ScheduleCron = types.StringNull()
	}
}
