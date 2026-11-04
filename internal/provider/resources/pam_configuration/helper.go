// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"fmt"
	"sort"
	"strings"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func appendFlagString(parts *[]string, flag string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	val := strings.TrimSpace(v.ValueString())
	if val == "" {
		return
	}
	*parts = append(*parts, fmt.Sprintf("%s %s", flag, utils.QuoteShellSingle(val)))
}

func boolToOnOff(b types.Bool) string {
	if !b.IsNull() && !b.IsUnknown() && b.ValueBool() {
		return ValueOn
	}
	return ValueOff
}

func boolToTrueFalse(b types.Bool) string {
	if !b.IsNull() && !b.IsUnknown() && b.ValueBool() {
		return ValueTrue
	}
	return ValueFalse
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

// appendPamConfigBodyFlags appends shared `pam config new|edit` flags (after the subcommand / UID).
func appendPamConfigBodyFlags(parts *[]string, data *commonpamconfiguration.PamConfigurationResourceModel) {
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagEnvironment, utils.QuoteShellSingle(data.Environment.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTitle, utils.QuoteShellSingle(data.Title.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagGateway, utils.QuoteShellSingle(data.Gateway.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagSharedFolder, utils.QuoteShellSingle(data.ApplicationFolder.ValueString())))

	appendFlagString(parts, FlagSchedule, data.Schedule)

	for _, pm := range sortedSetStrings(data.PortMapping) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagPortMapping, utils.QuoteShellSingle(pm)))
	}

	*parts = append(*parts, fmt.Sprintf("%s %s", FlagConnections, boolToOnOff(data.Connections)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTunneling, boolToOnOff(data.Tunneling)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagRotation, boolToOnOff(data.Rotation)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagRemoteBrowserIsolation, boolToOnOff(data.RemoteBrowserIsolation)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagConnectionsRecording, boolToOnOff(data.ConnectionsRecording)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTypescriptRecording, boolToOnOff(data.TypescriptRecording)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagAIThreatDetection, boolToOnOff(data.AIThreatDetection)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagAITerminateSessionOnDetection, boolToOnOff(data.AITerminateSessionOnDetection)))

	env := data.Environment.ValueString()
	switch env {
	case commonpamconfiguration.EnvLocal:
		appendLocalNetworkFlags(parts, data.LocalNetwork)
	case commonpamconfiguration.EnvAWS:
		appendAwsFlags(parts, data.Aws)
	case commonpamconfiguration.EnvAzure:
		appendAzureFlags(parts, data.Azure)
	case commonpamconfiguration.EnvDomain:
		appendDomainFlags(parts, data.Domain)
	case commonpamconfiguration.EnvGCP:
		appendGcpFlags(parts, data.Gcp)
	}
}

// buildPamConfigNewCommand builds `pam config new ...` from the resource model.
// Required attributes (environment, title, gateway, application_folder) are enforced by the schema.
func buildPamConfigNewCommand(data *commonpamconfiguration.PamConfigurationResourceModel) string {
	parts := []string{commonpamconfiguration.CmdPamConfig, commonpamconfiguration.CmdPamNew}
	appendPamConfigBodyFlags(&parts, data)
	return strings.Join(parts, " ")
}

// buildPamConfigEditCommand builds `pam config edit '<uid>' ...` from the resource model.
func buildPamConfigEditCommand(uid string, data *commonpamconfiguration.PamConfigurationResourceModel) string {
	parts := []string{commonpamconfiguration.CmdPamConfig, commonpamconfiguration.CmdPamEdit, utils.QuoteShellSingle(strings.TrimSpace(uid))}
	appendPamConfigBodyFlags(&parts, data)
	return strings.Join(parts, " ")
}

func appendLocalNetworkFlags(parts *[]string, m *commonpamconfiguration.PamLocalNetworkModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagNetworkId, m.NetworkId)
	appendFlagString(parts, FlagNetworkCidr, m.NetworkCidr)
}

func appendAwsFlags(parts *[]string, m *commonpamconfiguration.PamAwsModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagAwsId, m.AwsId)
	appendFlagString(parts, FlagAccessKeyId, m.AccessKeyId)
	appendFlagString(parts, FlagAccessSecretKey, m.AccessSecretKey)
	for _, r := range sortedSetStrings(m.RegionNames) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagRegionName, utils.QuoteShellSingle(r)))
	}
}

func appendAzureFlags(parts *[]string, m *commonpamconfiguration.PamAzureModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagAzureId, m.AzureId)
	appendFlagString(parts, FlagClientId, m.ClientId)
	appendFlagString(parts, FlagClientSecret, m.ClientSecret)
	appendFlagString(parts, FlagSubscriptionId, m.SubscriptionId)
	appendFlagString(parts, FlagTenantId, m.TenantId)
	for _, rg := range sortedSetStrings(m.ResourceGroups) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagResourceGroup, utils.QuoteShellSingle(rg)))
	}
}

func appendDomainFlags(parts *[]string, m *commonpamconfiguration.PamDomainModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagDomainId, m.DomainId)
	appendFlagString(parts, FlagDomainHostname, m.DomainHostname)
	appendFlagString(parts, FlagDomainPort, m.DomainPort)
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagDomainUseSsl, boolToTrueFalse(m.DomainUseSsl)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagDomainScanDcCidr, boolToTrueFalse(m.DomainScanDcCidr)))
	appendFlagString(parts, FlagDomainNetworkCidr, m.DomainNetworkCidr)
	appendFlagString(parts, FlagDomainAdmin, m.DomainAdmin)
	appendFlagString(parts, FlagDomainUserMatch, m.UserMatch)
}

func appendGcpFlags(parts *[]string, m *commonpamconfiguration.PamGcpModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagGcpId, m.GcpId)
	appendFlagString(parts, FlagServiceAccountKey, m.ServiceAccountKey)
	appendFlagString(parts, FlagGoogleAdminEmail, m.GoogleAdminEmail)
	appendFlagString(parts, FlagGcpRegion, m.GcpRegion)
}

// buildPamConfigRemoveCommand builds `pam config remove '<uid>'`.
func buildPamConfigRemoveCommand(uid string) string {
	return strings.Join([]string{commonpamconfiguration.CmdPamConfig, commonpamconfiguration.CmdPamRemove, utils.QuoteShellSingle(strings.TrimSpace(uid))}, " ")
}

// ussConfigured reports whether universal_secrets_sync is set for data's active environment.
func ussConfigured(data *commonpamconfiguration.PamConfigurationResourceModel) bool {
	return ussCommonModel(data) != nil
}

// ussCommonModel returns the active environment's universal_secrets_sync common fields, or nil if not set.
func ussCommonModel(data *commonpamconfiguration.PamConfigurationResourceModel) *commonpamconfiguration.PamCommonUniversalSecretsSyncModel {
	switch data.Environment.ValueString() {
	case commonpamconfiguration.EnvAWS:
		if data.Aws != nil && data.Aws.UniversalSecretsSync != nil {
			return &data.Aws.UniversalSecretsSync.PamCommonUniversalSecretsSyncModel
		}
	case commonpamconfiguration.EnvAzure:
		if data.Azure != nil && data.Azure.UniversalSecretsSync != nil {
			return &data.Azure.UniversalSecretsSync.PamCommonUniversalSecretsSyncModel
		}
	case commonpamconfiguration.EnvGCP:
		if data.Gcp != nil && data.Gcp.UniversalSecretsSync != nil {
			return &data.Gcp.UniversalSecretsSync.PamCommonUniversalSecretsSyncModel
		}
	}
	return nil
}

// ussVaultName returns azure.universal_secrets_sync.vault_name, or null for any other environment.
func ussVaultName(data *commonpamconfiguration.PamConfigurationResourceModel) types.String {
	if data.Environment.ValueString() == commonpamconfiguration.EnvAzure && data.Azure != nil && data.Azure.UniversalSecretsSync != nil {
		return data.Azure.UniversalSecretsSync.VaultName
	}
	return types.StringNull()
}

// ussChanged reports whether universal_secrets_sync differs between plan and state.
func ussChanged(plan, state *commonpamconfiguration.PamConfigurationResourceModel) bool {
	p, s := ussCommonModel(plan), ussCommonModel(state)
	if p == nil || s == nil {
		return p != s
	}
	return !p.Enabled.Equal(s.Enabled) ||
		!p.DryRun.Equal(s.DryRun) ||
		!p.SyncIdentity.Equal(s.SyncIdentity) ||
		!p.FolderLocation.Equal(s.FolderLocation) ||
		!ussVaultName(plan).Equal(ussVaultName(state))
}

// clearUssConfig removes universal_secrets_sync from every cloud environment block.
func clearUssConfig(data *commonpamconfiguration.PamConfigurationResourceModel) {
	if data.Aws != nil {
		data.Aws.UniversalSecretsSync = nil
	}
	if data.Azure != nil {
		data.Azure.UniversalSecretsSync = nil
	}
	if data.Gcp != nil {
		data.Gcp.UniversalSecretsSync = nil
	}
}

// appendPamUssConfigBodyFlags appends `pam universal-sync-config` body flags (enabled,
// folder, sync-identity, vault-name, dry-run) for the active environment's
// universal_secrets_sync block. Returns false (and appends nothing) if the environment
// doesn't support Universal Secrets Sync, or the universal_secrets_sync block isn't
// configured for it - callers should skip issuing the command entirely in that case.
func appendPamUssConfigBodyFlags(parts *[]string, data *commonpamconfiguration.PamConfigurationResourceModel) bool {
	common := ussCommonModel(data)
	if common == nil {
		return false
	}

	*parts = append(*parts, fmt.Sprintf("%s %s", commonpamconfiguration.FlagPamUssEnabled, boolToTrueFalse(common.Enabled)))
	for _, f := range sortedSetStrings(common.FolderLocation) {
		*parts = append(*parts, fmt.Sprintf("%s %s", commonpamconfiguration.FlagPamUssFolder, utils.QuoteShellSingle(f)))
	}
	appendFlagString(parts, commonpamconfiguration.FlagPamUssSyncIdentity, common.SyncIdentity)
	appendFlagString(parts, commonpamconfiguration.FlagPamUssVaultName, ussVaultName(data))
	*parts = append(*parts, fmt.Sprintf("%s %s", commonpamconfiguration.FlagPamUssDryRun, boolToTrueFalse(common.DryRun)))

	return true
}

// buildPamNewUssConfigCommand builds `pam universal-sync-config add --network '<uid>' ...`
// from the resource model. Returns "" if the environment doesn't support Universal Secrets
// Sync, or the universal_secrets_sync block isn't configured for it.
func buildPamNewUssConfigCommand(data *commonpamconfiguration.PamConfigurationResourceModel) string {
	parts := []string{commonpamconfiguration.CmdPamUniversalSyncConfigNew}
	appendFlagString(&parts, commonpamconfiguration.FlagPamUssNetwork, data.Id)

	if !appendPamUssConfigBodyFlags(&parts, data) {
		return ""
	}
	return strings.Join(parts, " ")
}

// buildPamEditUssConfigCommand builds `pam universal-sync-config edit --network '<uid>' ...`
// from the resource model. Returns "" if the environment doesn't support Universal Secrets
// Sync, or the universal_secrets_sync block isn't configured for it.
func buildPamEditUssConfigCommand(data *commonpamconfiguration.PamConfigurationResourceModel) string {
	parts := []string{commonpamconfiguration.CmdPamUniversalSyncConfigEdit}
	appendFlagString(&parts, commonpamconfiguration.FlagPamUssNetwork, data.Id)

	if !appendPamUssConfigBodyFlags(&parts, data) {
		return ""
	}
	return strings.Join(parts, " ")
}
