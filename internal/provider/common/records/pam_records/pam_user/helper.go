// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"encoding/json"
	"fmt"
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

	if !rs.RotationProfile.IsNull() && !rs.RotationProfile.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagRotationProfile, rs.RotationProfile.ValueString()))
	}

	if !rs.Configuration.IsNull() && !rs.Configuration.IsUnknown() {
		flagToUse := FlagConfig

		// If the rotation profile is IAM User, we need to use the IAM/AAD config flag.
		if rs.RotationProfile.ValueString() == RotProfileIAMUser {
			flagToUse = FlagIamAadConfig
		}
		parts = append(parts, fmt.Sprintf("%s %s", flagToUse, quoteShellSingle(rs.Configuration.ValueString())))
	}

	if !rs.Resource.IsNull() && !rs.Resource.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagResource, quoteShellSingle(rs.Resource.ValueString())))
	}

	if !rs.SaaSConfig.IsNull() && !rs.SaaSConfig.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", FlagSaaSConfig, quoteShellSingle(rs.SaaSConfig.ValueString())))
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
		!plan.Resource.Equal(state.Resource) ||
		!plan.SaaSConfig.Equal(state.SaaSConfig) ||
		!plan.Enabled.Equal(state.Enabled) ||
		!plan.ScheduleCron.Equal(state.ScheduleCron) ||
		!plan.ScheduleJSON.Equal(state.ScheduleJSON) ||
		!plan.OnDemand.Equal(state.OnDemand) ||
		!plan.ScheduleConfig.Equal(state.ScheduleConfig) ||
		!plan.Complexity.Equal(state.Complexity)
}

// ParseRotationInfoMessage parses the message lines from
// `pam rotation info -r <uid>` and populates the rotation settings on the
// state model. Fields not present in the response are preserved from the
// existing state (config-as-source-of-truth).
func MapRotationSettingsToState(rotInfo *PamRotationInfoResponse, rec *utils.VaultRecordGetResponse, existing *PamUserRotationSettings, state *PamUserSharedModel) {
	var rs *PamUserRotationSettings
	if existing != nil {
		rs = &PamUserRotationSettings{
			RotationProfile: existing.RotationProfile,
			Configuration:   existing.Configuration,
			Resource:        existing.Resource,
			Enabled:         existing.Enabled,
			ScheduleCron:    existing.ScheduleCron,
			ScheduleJSON:    existing.ScheduleJSON,
			OnDemand:        existing.OnDemand,
			ScheduleConfig:  existing.ScheduleConfig,
			Complexity:      existing.Complexity,
			SaaSConfig:      existing.SaaSConfig,
		}
	} else {
		rs = &PamUserRotationSettings{}
	}

	if rotInfo == nil {
		state.RotationSettings = rs
		return
	}

	// set configuration
	if uid := configurationFromRotInfoOrRecord(rotInfo, rec); uid != "" {
		rs.Configuration = types.StringValue(uid)
	}

	// extract dagDebug from vault record
	dagDebug := DagDebugResponseFromVaultRecord(rec)

	// set rotation profile
	if rotationProfile := setRotationProfile(dagDebug, rec.RotationProfile); rotationProfile != "" {
		rs.RotationProfile = types.StringValue(rotationProfile)
	}

	if rotInfo.Disable {
		rs.Enabled = types.BoolValue(false)
		rs.ScheduleCron = types.StringNull()
		rs.ScheduleJSON = types.StringNull()
		rs.OnDemand = types.BoolNull()
	} else {
		rs.Enabled = types.BoolValue(true)

		if rotInfo.ScheduleType == RotProfileScheduleTypeManual {
			rs.OnDemand = types.BoolValue(true)
		} else {
			rs.OnDemand = types.BoolValue(false)
		}
	}

	if !mapSaaSRotationFromDagDebug(dagDebug, rs) &&
		rec.RotationProfile != nil &&
		rec.RotationProfile.Type == RotProfileGeneral {
		rs.Resource = types.StringValue(rec.RotationProfile.ResourceUID)
	}

	parseScheduleValue(rotInfo.ScheduleData, rs)

	rs.Complexity = MapPasswordComplexityResponseToState(rotInfo.PasswordComplexityDetails)
	state.RotationSettings = rs
}

func setRotationProfile(dagDebug *utils.DagDebugResponse, rotationProfile *utils.RotationProfileResponse) string {
	if dagDebug == nil {
		return ""
	}

	// SaaS is inferred from parent_acl_edge; rotationProfile.type is never "saas".
	if rotSettings := parentAclEdgeRotationSettings(dagDebug); rotSettings != nil {
		if rotSettings.Noop && len(rotSettings.SaaSRecordUIDList) > 0 {
			return RotProfileSaaS
		}
	}

	if rotationProfile == nil {
		return ""
	}
	return rotationProfile.Type
}

// mapSaaSRotationFromDagDebug infers SaaS rotation from parent_acl_edge noop + saas_record_uid_list.
// dagDebug.rotation_profile.type never returns "saas".
func mapSaaSRotationFromDagDebug(dagDebug *utils.DagDebugResponse, rs *PamUserRotationSettings) bool {
	rotSettings := parentAclEdgeRotationSettings(dagDebug)
	if rotSettings == nil {
		return false
	}

	saasUIDs := rotSettings.SaaSRecordUIDList
	if !rotSettings.Noop || len(saasUIDs) == 0 {
		return false
	}

	if uid := strings.TrimSpace(saasUIDs[0]); uid != "" {
		rs.SaaSConfig = types.StringValue(uid)
	}
	rs.Resource = types.StringNull()
	return true
}

func parentAclEdgeRotationSettings(dagDebug *utils.DagDebugResponse) *utils.DagDebugParentAclEdgeContentRotationSettingsResponse {
	if dagDebug == nil {
		return nil
	}
	parentEdge := dagDebug.ParentAclEdge
	if parentEdge == nil || parentEdge.Content == nil {
		return nil
	}
	return parentEdge.Content.RotationSettings
}

func configurationFromRotInfoOrRecord(rotInfo *PamRotationInfoResponse, rec *utils.VaultRecordGetResponse) string {
	if rotInfo != nil {
		if uid := strings.TrimSpace(rotInfo.PamConfigUID); uid != "" {
			return uid
		}
	}
	if rec != nil {
		return strings.TrimSpace(rec.PamConfigurationUID)
	}
	return ""
}

// DagDebugResponseFromVaultRecord returns dagDebug from a vault record, if present.
func DagDebugResponseFromVaultRecord(rec *utils.VaultRecordGetResponse) *utils.DagDebugResponse {
	if rec == nil || rec.DagDebug == nil {
		return nil
	}
	return rec.DagDebug
}

// MapPasswordComplexityResponseToState maps the password complexity response from the API to the state model.
func MapPasswordComplexityResponseToState(response *PamRotationInfoPasswordComplexityDetailsDataResponse) types.String {
	if response == nil {
		return types.StringNull()
	}
	return types.StringValue(fmt.Sprintf("%d,%d,%d,%d,%d", response.Length, response.Capital, response.Lowercase, response.Digits, response.Special))
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
	case "DAILY", "WEEKLY", "MONTHLY_BY_WEEKDAY", "YEARLY":
		rs.ScheduleJSON = types.StringValue(raw)
		rs.OnDemand = types.BoolNull()
		rs.ScheduleCron = types.StringNull()
	}
}
