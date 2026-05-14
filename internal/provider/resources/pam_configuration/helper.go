// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"fmt"
	"sort"
	"strings"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
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
	*parts = append(*parts, fmt.Sprintf("%s %s", flag, commonpamconfiguration.QuoteShellSingle(val)))
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
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagEnvironment, commonpamconfiguration.QuoteShellSingle(data.Environment.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTitle, commonpamconfiguration.QuoteShellSingle(data.Title.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagGateway, commonpamconfiguration.QuoteShellSingle(data.Gateway.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagSharedFolder, commonpamconfiguration.QuoteShellSingle(data.ApplicationFolder.ValueString())))

	appendFlagString(parts, FlagSchedule, data.Schedule)

	for _, pm := range sortedSetStrings(data.PortMapping) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagPortMapping, commonpamconfiguration.QuoteShellSingle(pm)))
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
	parts := []string{commonpamconfiguration.CmdPamConfig, commonpamconfiguration.CmdPamEdit, commonpamconfiguration.QuoteShellSingle(strings.TrimSpace(uid))}
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
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagRegionName, commonpamconfiguration.QuoteShellSingle(r)))
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
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagResourceGroup, commonpamconfiguration.QuoteShellSingle(rg)))
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
	return strings.Join([]string{commonpamconfiguration.CmdPamConfig, commonpamconfiguration.CmdPamRemove, commonpamconfiguration.QuoteShellSingle(strings.TrimSpace(uid))}, " ")
}
