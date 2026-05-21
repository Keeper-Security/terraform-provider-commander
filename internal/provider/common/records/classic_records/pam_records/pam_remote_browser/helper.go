// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pamremotebrowser holds shared Terraform schema fragments for PAM remote browser resources.
package pamremotebrowser

import (
	"context"
	"encoding/json"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

	state.Folder = commonpamrecords.ExtractFolderValue(rec.Folder, state.Folder)

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
