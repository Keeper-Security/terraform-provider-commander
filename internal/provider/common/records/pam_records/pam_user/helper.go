// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// quoteShellSingle wraps s for use as a single-quoted shell argument.
func quoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

// BuildAddCommand builds a record-add style command for a PAM User record.
// The CLI command (`record-add` or `nsf-record-add`) is provided by the caller.
//
// Example output:
//
//	record-add --folder 'nptw7AealHr8LC-hk8aAHg' --title 'AD service account' --record-type pamUser
//	  f.login=svc_myapp f.password='123' f.text.distinguishedName='CN=...' f.secret.privatePEMKey='abc'
//	  f.text.connectDatabase='test' f.checkbox.managed=true
func BuildAddCommand(cmd string, data PamUserSharedModel) string {
	parts := []string{cmd}

	if !data.FolderLocation.IsNull() && !data.FolderLocation.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagFolder, quoteShellSingle(data.FolderLocation.ValueString())))
	}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, quoteShellSingle(data.Title.ValueString())))

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamUser))

	if !data.Login.IsNull() && !data.Login.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldLogin, quoteShellSingle(data.Login.ValueString())))
	}

	commonpamrecords.AppendOptionalTextField(&parts, FieldPassword, data.Password)
	commonpamrecords.AppendOptionalTextField(&parts, FieldDistinguishedName, data.DistinguishedName)
	commonpamrecords.AppendOptionalTextField(&parts, FieldPrivatePEMKey, data.PrivatePEMKey)
	commonpamrecords.AppendOptionalTextField(&parts, FieldPublicKey, data.PublicKey)
	commonpamrecords.AppendOptionalTextField(&parts, FieldPrivateKeyPassphrase, data.PrivateKeyPassphrase)
	commonpamrecords.AppendOptionalTextField(&parts, FieldConnectDatabase, data.ConnectDatabase)
	commonpamrecords.AppendOptionalCheckboxField(&parts, FieldManaged, data.Managed)

	if !data.Notes.IsNull() && !data.Notes.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, quoteShellSingle(data.Notes.ValueString())))
	}

	return strings.Join(parts, " ")
}

// UpdateHasMutations returns true when at least one record-level field differs
// between plan and state (rotation_settings is checked separately).
func UpdateHasMutations(plan, state PamUserSharedModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Password.Equal(state.Password) ||
		!plan.DistinguishedName.Equal(state.DistinguishedName) ||
		!plan.PrivatePEMKey.Equal(state.PrivatePEMKey) ||
		!plan.PublicKey.Equal(state.PublicKey) ||
		!plan.PrivateKeyPassphrase.Equal(state.PrivateKeyPassphrase) ||
		!plan.ConnectDatabase.Equal(state.ConnectDatabase) ||
		!plan.Managed.Equal(state.Managed) ||
		(!plan.Notes.Equal(state.Notes) && !plan.Notes.IsNull() && !plan.Notes.IsUnknown())
}

// BuildUpdateCommand builds a record-update style command for a PAM User
// record, including only the fields that changed between plan and state.
//
// Example output:
//
//	record-update --record 'nb-dp0Xm7KZWRTKt5UK83w' --title 'AD service account' f.login='svc_myapp'
func BuildUpdateCommand(cmd, recordUID string, plan, state PamUserSharedModel) string {
	parts := []string{
		cmd,
		fmt.Sprintf("%s %s", utils.FlagRecord, quoteShellSingle(recordUID)),
	}

	if !plan.Title.Equal(state.Title) {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, quoteShellSingle(plan.Title.ValueString())))
	}

	if !plan.Login.Equal(state.Login) {
		parts = append(parts, fmt.Sprintf("%s=%s", FieldLogin, quoteShellSingle(plan.Login.ValueString())))
	}

	if !plan.Password.Equal(state.Password) {
		// if the password is null or unknown, we need to add the flag with an empty value
		if plan.Password.IsNull() || plan.Password.IsUnknown() {
			parts = append(parts, fmt.Sprintf("'%s='", FieldPassword))
		} else {
			parts = append(parts, fmt.Sprintf("'%s=%s'", FieldPassword, quoteShellSingle(plan.Password.ValueString())))
		}
	}

	commonpamrecords.AppendChangedTextField(&parts, FieldDistinguishedName, plan.DistinguishedName, state.DistinguishedName)
	commonpamrecords.AppendChangedTextField(&parts, FieldPrivatePEMKey, plan.PrivatePEMKey, state.PrivatePEMKey)
	commonpamrecords.AppendChangedTextField(&parts, FieldPublicKey, plan.PublicKey, state.PublicKey)
	commonpamrecords.AppendChangedTextField(&parts, FieldPrivateKeyPassphrase, plan.PrivateKeyPassphrase, state.PrivateKeyPassphrase)
	commonpamrecords.AppendChangedTextField(&parts, FieldConnectDatabase, plan.ConnectDatabase, state.ConnectDatabase)
	commonpamrecords.AppendChangedCheckboxField(&parts, FieldManaged, plan.Managed, state.Managed)

	if !plan.Notes.Equal(state.Notes) && !plan.Notes.IsUnknown() {
		if plan.Notes.IsNull() {
			parts = append(parts, fmt.Sprintf("%s ''", utils.FlagNotes))
		} else {
			parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, plan.Notes.ValueString()))
		}
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

// MapVaultRecordToState maps a `get <uid> --format json` response onto the
// shared PAM User model.
func MapVaultRecordToState(rec *utils.VaultRecordGetResponse, state *PamUserSharedModel) {
	if strings.TrimSpace(rec.RecordUID) != "" {
		state.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}

	state.Title = stringOrNull(rec.Title)
	state.Notes = stringOrNull(rec.Notes)
	state.FolderLocation = utils.ExtractFolderValue(rec.FolderLocation, state.FolderLocation)

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
			if f.Label == "publicKey" {
				state.PublicKey = firstStringValue(f.Value)
			}
			if f.Label == "privateKeyPassphrase" {
				state.PrivateKeyPassphrase = firstStringValue(f.Value)
			}
		case vaultFieldTypeCheckbox:
			if f.Label == "managed" {
				state.Managed = firstBoolValue(f.Value)
			}
		}
	}
}

func stringOrNull(s string) types.String {
	if strings.TrimSpace(s) == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

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

// BuildPamRotationEditCommand builds `pam rotation edit -r <uid> --config ... --force`.
func BuildPamRotationEditCommand(recordUID string, rs *PamUserRotationSettings) string {
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

// RotationSettingsNeedApply returns true when the plan has a rotation block
// that differs from state.
func RotationSettingsNeedApply(plan, state *PamUserRotationSettings) bool {
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

// HasRotationData returns true when the message lines contain a PAM Config
// UID, indicating the record has rotation configured in the vault.
func HasRotationData(messages []string) bool {
	for _, line := range messages {
		if strings.HasPrefix(strings.TrimSpace(line), "PAM Config UID:") {
			return true
		}
	}
	return false
}

// ParseRotationInfoMessage parses the message lines from
// `pam rotation info -r <uid>` and populates the rotation settings on the
// state model. Fields not present in the response are preserved from the
// existing state (config-as-source-of-truth).
func ParseRotationInfoMessage(messages []string, existing *PamUserRotationSettings, state *PamUserSharedModel) {
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
				if rs.RotationProfile.ValueString() != "general" {
					rs.Configuration = stringOrNull(v)
				} else {
					rs.Configuration = types.StringNull()
				}
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
				rs.Complexity = ParsePasswordComplexityData(v)
			}
		}
	}

	state.RotationSettings = rs
}

// ParsePasswordComplexityData converts Commander rotation info complexity text
// (e.g. "Length: 32; Lowercase: 1; Uppercase: 5; Digits: 1; Symbols: 2; Symbols Chars: ...")
// into the Terraform schema format: length,upper,lower,digits,symbols.
func ParsePasswordComplexityData(raw string) types.String {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return types.StringNull()
	}

	length, upper, lower, digits, symbols, ok := parsePasswordComplexityFields(raw)
	if !ok {
		return types.StringNull()
	}

	return types.StringValue(fmt.Sprintf("%d,%d,%d,%d,%d", length, upper, lower, digits, symbols))
}

func complexityDataWithoutSymbolsChars(raw string) string {
	lower := strings.ToLower(raw)
	marker := "; symbols chars:"
	if idx := strings.Index(lower, marker); idx >= 0 {
		return strings.TrimSpace(raw[:idx])
	}
	if strings.HasPrefix(lower, "symbols chars:") {
		return ""
	}
	return raw
}

func parsePasswordComplexityFields(raw string) (length, upper, lower, digits, symbols int, ok bool) {
	const missing = -1
	length, upper, lower, digits, symbols = missing, missing, missing, missing, missing

	for _, segment := range strings.Split(complexityDataWithoutSymbolsChars(raw), ";") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		key, value, found := strings.Cut(segment, ":")
		if !found {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			continue
		}

		switch strings.ToLower(strings.TrimSpace(key)) {
		case "length":
			length = n
		case "uppercase":
			upper = n
		case "lowercase":
			lower = n
		case "digits":
			digits = n
		case "symbols":
			symbols = n
		}
	}

	ok = length != missing && upper != missing && lower != missing && digits != missing && symbols != missing
	return length, upper, lower, digits, symbols, ok
}

// parseScheduleValue parses the JSON schedule from `pam rotation info`
// responses such as:
//   - {"type":"CRON","cron":"0 56 17 * * ?","tz":"Etc/UTC"}
//   - {"type":"DAILY","utcTime":"17:56","intervalCount":1}
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

// ParseFlexibleMessageToLines converts the FlexibleMessage string (may be a
// JSON array or plain text) into individual lines for parsing.
func ParseFlexibleMessageToLines(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err == nil {
		return arr
	}

	return strings.Split(raw, "\n")
}
