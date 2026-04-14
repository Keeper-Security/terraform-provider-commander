// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"fmt"
	"sort"
	"strings"

	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// quoteShellSingle wraps s for use as a single-quoted shell argument (bash-style escaping of ').
func quoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

// appendPamRbiStringFlag appends `flag '<value>'`. Null (attribute removed from config) sends an empty argument; unknown is omitted.
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

// appendPamRbiBoolOnOff appends `flag on|off`. Null (attribute removed from config) sends `off`; unknown is omitted.
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

// appendPamRbiRepeatedStringFlags appends one or more `flag '<value>'` arguments.
// Null set (attribute omitted) produces a single `flag ”` for CLI clear semantics.
// Non-null sets must be non-empty per schema; unknown set is omitted.
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

// AppendPamRbiEditSettingsFlags appends all `pam rbi edit` flags except the leading command and `--record`.
// Shared by Create (phase 2) and Update. Bool attributes unset in config are sent as `off`; string attributes unset as `”`.
// Set attributes omitted (null) send one `flag ”` per list field; non-null sets must be non-empty in schema.
func AppendPamRbiEditSettingsFlags(parts *[]string, settings *commonpamremotebrowser.PamRemoteBrowserSettingsModel) {
	if settings == nil {
		return
	}

	appendPamRbiStringFlag(parts, FlagConfiguration, settings.Configuration)

	appendPamRbiBoolOnOff(parts, FlagRemoteBrowserIsolation, settings.RemoteBrowserIsolation)
	appendPamRbiBoolOnOff(parts, FlagConnectionsRecording, settings.ConnectionsRecording)
	appendPamRbiBoolOnOff(parts, FlagKeyEvents, settings.KeyEvents)
	appendPamRbiBoolOnOff(parts, FlagAllowURLNavigation, settings.AllowUrlNavigation)
	appendPamRbiBoolOnOff(parts, FlagIgnoreServerCert, settings.IgnoreServerCert)

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

// BuildPamRbiEditCommand builds `pam rbi edit --record <uid> ...` for Create phase 2 and Update.
func BuildPamRbiEditCommand(recordUID string, settings *commonpamremotebrowser.PamRemoteBrowserSettingsModel) string {
	parts := []string{
		CmdPamRbiEdit,
		fmt.Sprintf("%s %s", FlagRecord, quoteShellSingle(recordUID)),
	}
	AppendPamRbiEditSettingsFlags(&parts, settings)
	return strings.Join(parts, " ")
}

// recordUpdateHasMutations is true when record-update will include at least one field flag besides --record
// (omits notes/folder when plan clears them—no flag is sent).
func recordUpdateHasMutations(plan, state commonpamremotebrowser.PamRemoteBrowserResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!plan.Url.Equal(state.Url) ||
		(!plan.Notes.Equal(state.Notes) && !plan.Notes.IsNull() && !plan.Notes.IsUnknown()) ||
		(!plan.Folder.Equal(state.Folder) && !plan.Folder.IsNull() && !plan.Folder.IsUnknown())
}

// buildUpdatePamRemoteBrowserRecordCommand builds `record-update --record <uid> ...` with only flags for fields that changed.
func buildUpdatePamRemoteBrowserRecordCommand(recordUID string, plan, state commonpamremotebrowser.PamRemoteBrowserResourceModel) string {
	parts := []string{
		utils.CmdRecordUpdate,
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
	// if !plan.Folder.Equal(state.Folder) && !plan.Folder.IsNull() && !plan.Folder.IsUnknown() {
	// 	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagFolder, quoteShellSingle(plan.Folder.ValueString())))
	// }

	return strings.Join(parts, " ")
}

func pamRemoteBrowserSettingsEqual(plan, state *commonpamremotebrowser.PamRemoteBrowserSettingsModel) bool {
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
		plan.AudioSampleRate.Equal(state.AudioSampleRate)
}

// pamRemoteBrowserSettingsNeedApply is true when plan has a settings block that differs from state (including first-time apply).
func pamRemoteBrowserSettingsNeedApply(plan, state *commonpamremotebrowser.PamRemoteBrowserSettingsModel) bool {
	if plan == nil {
		return false
	}
	if state == nil {
		return true
	}
	return !pamRemoteBrowserSettingsEqual(plan, state)
}
