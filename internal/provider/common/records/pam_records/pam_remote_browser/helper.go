// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pamremotebrowser holds shared Terraform schema fragments for PAM remote browser resources.
package pamremotebrowser

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// quoteShellSingle wraps s for use as a single-quoted shell argument
// (bash-style escaping of ').
func quoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

// BuildAddCommand builds a record-add style command for a PAM Remote Browser
// record. The CLI command (`record-add` or `nsf-record-add`) is provided by
// the caller. Settings are applied separately via BuildPamRbiEditCommand.
func BuildAddCommand(cmd string, data PamRemoteBrowserResourceModel) string {
	parts := []string{cmd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamRemoteBrowser))
	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))
	parts = append(parts, fmt.Sprintf("'%s=%s'", utils.FlagRbiUrl, data.Url.ValueString()))

	if !data.FolderLocation.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagFolder, data.FolderLocation.ValueString()))
	}

	if !data.Notes.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}

// BuildUpdateCommand builds a record-update style command for a PAM Remote
// Browser record, including only flags for fields that changed. Settings are
// applied separately via BuildPamRbiEditCommand.
func BuildUpdateCommand(cmd, recordUID string, plan, state PamRemoteBrowserResourceModel) string {
	parts := []string{
		cmd,
		fmt.Sprintf("%s %s", utils.FlagRecord, quoteShellSingle(recordUID)),
	}
	if !plan.Title.Equal(state.Title) {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagTitle, quoteShellSingle(plan.Title.ValueString())))
	}
	if !plan.Url.Equal(state.Url) {
		parts = append(parts, fmt.Sprintf("'%s=%s'", utils.FlagRbiUrl, plan.Url.ValueString()))
	}
	if !plan.Notes.Equal(state.Notes) && !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		parts = append(parts, fmt.Sprintf("%s %s", utils.FlagNotes, quoteShellSingle(plan.Notes.ValueString())))
	}

	return strings.Join(parts, " ")
}

// RecordUpdateHasMutations returns true when record-update will include at
// least one field flag besides --record (omits notes/folder when plan clears
// them—no flag is sent).
func RecordUpdateHasMutations(plan, state PamRemoteBrowserResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!plan.Url.Equal(state.Url) ||
		(!plan.Notes.Equal(state.Notes) && !plan.Notes.IsNull() && !plan.Notes.IsUnknown()) ||
		(!plan.FolderLocation.Equal(state.FolderLocation) && !plan.FolderLocation.IsNull() && !plan.FolderLocation.IsUnknown())
}

// BuildPamRbiEditCommand builds `pam rbi edit --record <uid> ...` shared by
// Create (phase 2) and Update.
func BuildPamRbiEditCommand(recordUID string, settings *PamRemoteBrowserSettingsModel) string {
	parts := []string{
		CmdPamRbiEdit,
		fmt.Sprintf("%s %s", FlagRecord, quoteShellSingle(recordUID)),
	}
	AppendPamRbiEditSettingsFlags(&parts, settings)
	return strings.Join(parts, " ")
}

// AppendPamRbiEditSettingsFlags appends all `pam rbi edit` flags except the
// leading command and `--record`.
func AppendPamRbiEditSettingsFlags(parts *[]string, settings *PamRemoteBrowserSettingsModel) {
	if settings == nil {
		return
	}

	appendPamRbiStringFlag(parts, FlagConfiguration, settings.Configuration)

	appendPamRbiBoolOnOff(parts, FlagRemoteBrowserIsolation, settings.RemoteBrowserIsolation)
	appendPamRbiBoolOnOff(parts, FlagConnectionsRecording, settings.ConnectionsRecording)
	appendPamRbiBoolOnOff(parts, FlagKeyEvents, settings.KeyEvents)
	appendPamRbiBoolOnOff(parts, FlagAllowURLNavigation, settings.AllowUrlNavigation)
	appendPamRbiBoolOnOff(parts, FlagIgnoreServerCert, settings.IgnoreServerCert)

	appendPamRbiStringFlag(parts, FlagSessionPersistence, settings.SessionPersistence)
	appendPamRbiRepeatedStringFlags(parts, FlagAllowedURLs, settings.AllowedUrls)
	appendPamRbiRepeatedStringFlags(parts, FlagAllowedResourceURLs, settings.AllowedResourceUrls)

	appendPamRbiStringFlag(parts, FlagAutofillCredentials, settings.AutoFillCredentials)
	appendPamRbiRepeatedStringFlags(parts, FlagAutofillTargets, settings.AutoFillTargets)

	appendPamRbiBoolOnOff(parts, FlagAllowCopy, settings.AllowCopy)
	appendPamRbiBoolOnOff(parts, FlagAllowPaste, settings.AllowPaste)
	appendPamRbiBoolOnOff(parts, FlagDisableAudio, settings.DisableAudio)

	appendOptionalInt32(parts, FlagAudioChannels, settings.AudioChannels)
	appendOptionalInt64(parts, FlagAudioBitDepth, settings.AudioBitDepth)
	appendOptionalInt64(parts, FlagAudioSampleRate, settings.AudioSampleRate)
}

// appendPamRbiStringFlag appends `flag '<value>'`. Null (attribute removed
// from config) sends an empty argument; unknown is omitted.
func appendPamRbiStringFlag(parts *[]string, flag string, v types.String) {
	if v.IsUnknown() {
		return
	}
	if v.IsNull() {
		*parts = append(*parts, fmt.Sprintf("%s %s", flag, quoteShellSingle("")))
		return
	}
	*parts = append(*parts, fmt.Sprintf("%s %s", flag, quoteShellSingle(v.ValueString())))
}

// appendPamRbiBoolOnOff appends `flag on|off`. Null (attribute removed from
// config) sends `off`; unknown is omitted.
func appendPamRbiBoolOnOff(parts *[]string, flag string, v types.Bool) {
	if v.IsUnknown() {
		return
	}
	if v.IsNull() {
		*parts = append(*parts, fmt.Sprintf("%s %s", flag, utils.ValueOff))
		return
	}
	val := utils.ValueOff
	if v.ValueBool() {
		val = utils.ValueOn
	}
	*parts = append(*parts, fmt.Sprintf("%s %s", flag, val))
}

func appendOptionalInt32(parts *[]string, flag string, v types.Int32) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	*parts = append(*parts, fmt.Sprintf("%s %d", flag, v.ValueInt32()))
}

func appendOptionalInt64(parts *[]string, flag string, v types.Int64) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	*parts = append(*parts, fmt.Sprintf("%s %d", flag, v.ValueInt64()))
}

// sortedSetStrings returns sorted non-null string elements from a Terraform set.
func sortedSetStrings(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var out []string
	for _, elem := range set.Elements() {
		s, ok := elem.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			continue
		}
		out = append(out, s.ValueString())
	}
	sort.Strings(out)
	return out
}

// appendPamRbiRepeatedStringFlags appends one or more `flag '<value>'`
// arguments. Null set (attribute omitted) produces a single `flag ”` for CLI
// clear semantics. Non-null sets must be non-empty per schema; unknown set is
// omitted.
func appendPamRbiRepeatedStringFlags(parts *[]string, flag string, set types.Set) {
	if set.IsUnknown() {
		return
	}
	if set.IsNull() {
		*parts = append(*parts, fmt.Sprintf("%s %s", flag, quoteShellSingle("")))
		return
	}
	for _, s := range sortedSetStrings(set) {
		t := strings.TrimSpace(s)
		if t == "" {
			continue
		}
		*parts = append(*parts, fmt.Sprintf("%s %s", flag, quoteShellSingle(t)))
	}
}

// PamRemoteBrowserSettingsNeedApply is true when plan has a settings block
// that differs from state (including first-time apply).
func PamRemoteBrowserSettingsNeedApply(plan, state *PamRemoteBrowserSettingsModel) bool {
	if plan == nil {
		return false
	}
	if state == nil {
		return true
	}
	return !pamRemoteBrowserSettingsEqual(plan, state)
}

func pamRemoteBrowserSettingsEqual(plan, state *PamRemoteBrowserSettingsModel) bool {
	if plan == nil && state == nil {
		return true
	}
	if plan == nil || state == nil {
		return false
	}
	return plan.Configuration.Equal(state.Configuration) &&
		plan.RemoteBrowserIsolation.Equal(state.RemoteBrowserIsolation) &&
		plan.ConnectionsRecording.Equal(state.ConnectionsRecording) &&
		plan.KeyEvents.Equal(state.KeyEvents) &&
		plan.AllowUrlNavigation.Equal(state.AllowUrlNavigation) &&
		plan.IgnoreServerCert.Equal(state.IgnoreServerCert) &&
		plan.AllowedUrls.Equal(state.AllowedUrls) &&
		plan.AllowedResourceUrls.Equal(state.AllowedResourceUrls) &&
		plan.AutoFillTargets.Equal(state.AutoFillTargets) &&
		plan.AutoFillCredentials.Equal(state.AutoFillCredentials) &&
		plan.AllowCopy.Equal(state.AllowCopy) &&
		plan.AllowPaste.Equal(state.AllowPaste) &&
		plan.DisableAudio.Equal(state.DisableAudio) &&
		plan.AudioChannels.Equal(state.AudioChannels) &&
		plan.AudioBitDepth.Equal(state.AudioBitDepth) &&
		plan.AudioSampleRate.Equal(state.AudioSampleRate) &&
		plan.SessionPersistence.Equal(state.SessionPersistence)
}

const (
	vaultFieldTypeRbiURL                   = "rbiUrl"
	vaultFieldTypePamRemoteBrowserSettings = "pamRemoteBrowserSettings"
)

// MapVaultRecordGetResponseToPamRemoteBrowserModel fills state from `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPamRemoteBrowserModel(ctx context.Context, rec *utils.VaultRecordGetResponse, state *PamRemoteBrowserResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if strings.TrimSpace(rec.RecordUID) != "" {
		state.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}
	if strings.TrimSpace(rec.Title) == "" {
		state.Title = types.StringNull()
	} else {
		state.Title = types.StringValue(rec.Title)
	}

	if strings.TrimSpace(rec.Notes) == "" {
		state.Notes = types.StringNull()
	} else {
		state.Notes = types.StringValue(rec.Notes)
	}

	state.FolderLocation = utils.ExtractFolderValue(rec.FolderLocation, state.FolderLocation)

	var rbiURL string
	var settingsConn *utils.PamRemoteBrowserSettingsFieldConnectionResponse

	for i := range rec.Fields {
		f := &rec.Fields[i]
		switch f.Type {
		case vaultFieldTypeRbiURL:
			var vals []string
			if err := json.Unmarshal(f.Value, &vals); err != nil {
				diags.AddWarning("PAM remote browser read", "Could not parse rbiUrl field: "+err.Error())
				continue
			}
			if len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
				rbiURL = strings.TrimSpace(vals[0])
			}
		case vaultFieldTypePamRemoteBrowserSettings:
			var entries []utils.PamRemoteBrowserSettingsFieldResponse
			if err := json.Unmarshal(f.Value, &entries); err != nil {
				diags.AddWarning("PAM remote browser read", "Could not parse pamRemoteBrowserSettings field: "+err.Error())
				continue
			}
			if len(entries) > 0 {
				settingsConn = &entries[0].Connection
			}
		}
	}

	if rbiURL != "" {
		state.Url = types.StringValue(rbiURL)
	} else {
		state.Url = types.StringNull()
	}

	if settingsConn != nil {
		m, d := mapConnectionToPamRemoteBrowserSettingsModel(ctx, settingsConn, rec)
		diags.Append(d...)
		state.PamRemoteBrowserSettings = m
	} else {
		state.PamRemoteBrowserSettings = nil
	}

	return diags
}

func mapConnectionToPamRemoteBrowserSettingsModel(ctx context.Context, c *utils.PamRemoteBrowserSettingsFieldConnectionResponse, rec *utils.VaultRecordGetResponse) (*PamRemoteBrowserSettingsModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	autofillTokens := splitAutofillTargetTokens(c.AutofillConfiguration)

	configUID := types.StringNull()
	if strings.TrimSpace(rec.PamConfigurationUID) != "" {
		configUID = types.StringValue(strings.TrimSpace(rec.PamConfigurationUID))
	}

	remoteBrowserIsolation := types.BoolNull()
	connectionsRecording := types.BoolNull()
	if rec.ConfigurationAllowedSettings != nil {
		remoteBrowserIsolation = types.BoolValue(rec.ConfigurationAllowedSettings.RemoteBrowserIsolation)
		connectionsRecording = types.BoolValue(rec.ConfigurationAllowedSettings.ConnectionsRecording)
	}

	m := &PamRemoteBrowserSettingsModel{
		Configuration:          configUID,
		RemoteBrowserIsolation: remoteBrowserIsolation,
		ConnectionsRecording:   connectionsRecording,
		KeyEvents:              types.BoolValue(c.RecordingIncludeKeys),
		AllowUrlNavigation:     types.BoolValue(c.AllowUrlManipulation),
		IgnoreServerCert:       types.BoolValue(c.IgnoreInitialSslCert),
		AllowCopy:              types.BoolValue(!c.DisableCopy),
		AllowPaste:             types.BoolValue(!c.DisablePaste),
		DisableAudio:           types.BoolValue(c.DisableAudio),
		SessionPersistence:     types.StringValue(c.SessionPersistence),
	}

	if strings.TrimSpace(c.HttpCredentialsUID) == "" {
		m.AutoFillCredentials = types.StringNull()
	} else {
		m.AutoFillCredentials = types.StringValue(strings.TrimSpace(c.HttpCredentialsUID))
	}

	urlSet, d := stringSliceToStringSet(ctx, nonEmptyLines(c.AllowedURLPatterns))
	diags.Append(d...)
	m.AllowedUrls = urlSet

	resSet, d := stringSliceToStringSet(ctx, nonEmptyLines(c.AllowedResourceURLPatterns))
	diags.Append(d...)
	m.AllowedResourceUrls = resSet

	tgtSet, d := stringSliceToStringSet(ctx, autofillTokens)
	diags.Append(d...)
	m.AutoFillTargets = tgtSet

	if c.AudioChannels > 0 {
		m.AudioChannels = types.Int32Value(int32(c.AudioChannels))
	} else {
		m.AudioChannels = types.Int32Null()
	}
	if c.AudioBps > 0 {
		m.AudioBitDepth = types.Int64Value(int64(c.AudioBps))
	} else {
		m.AudioBitDepth = types.Int64Null()
	}
	if c.AudioSampleRate > 0 {
		m.AudioSampleRate = types.Int64Value(int64(c.AudioSampleRate))
	} else {
		m.AudioSampleRate = types.Int64Null()
	}

	return m, diags
}

func nonEmptyLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var out []string
	for _, line := range strings.Split(s, "\n") {
		t := strings.TrimSpace(line)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func splitAutofillTargetTokens(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	seen := make(map[string]struct{})
	var out []string
	for _, chunk := range strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	}) {
		t := strings.TrimSpace(chunk)
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func stringSliceToStringSet(ctx context.Context, vals []string) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(vals) == 0 {
		return types.SetNull(types.StringType), diags
	}
	elems := make([]attr.Value, len(vals))
	for i, v := range vals {
		elems[i] = types.StringValue(v)
	}
	set, d := types.SetValueFrom(ctx, types.StringType, elems)
	diags.Append(d...)
	return set, diags
}
