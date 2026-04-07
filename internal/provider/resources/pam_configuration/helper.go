// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// quoteShellSingle wraps s for use as a single-quoted shell argument (bash-style escaping of ').
func quoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

func appendFlagString(parts *[]string, flag string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	val := strings.TrimSpace(v.ValueString())
	if val == "" {
		return
	}
	*parts = append(*parts, fmt.Sprintf("%s %s", flag, quoteShellSingle(val)))
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
// isUpdate: true for `pam config edit` (omits AI permission flags; Commander does not support them on edit yet).
// false for `pam config new`.
func appendPamConfigBodyFlags(parts *[]string, data *PamConfigurationResourceModel, isUpdate bool) {
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagEnvironment, quoteShellSingle(data.Environment.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTitle, quoteShellSingle(data.Title.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagGateway, quoteShellSingle(data.Gateway.ValueString())))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagSharedFolder, quoteShellSingle(data.ApplicationFolder.ValueString())))

	appendFlagString(parts, FlagSchedule, data.Schedule)

	for _, pm := range sortedSetStrings(data.PortMapping) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagPortMapping, quoteShellSingle(pm)))
	}

	*parts = append(*parts, fmt.Sprintf("%s %s", FlagConnections, boolToOnOff(data.Connections)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTunneling, boolToOnOff(data.Tunneling)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagRotation, boolToOnOff(data.Rotation)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagRemoteBrowserIsolation, boolToOnOff(data.RemoteBrowserIsolation)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagConnectionsRecording, boolToOnOff(data.ConnectionsRecording)))
	*parts = append(*parts, fmt.Sprintf("%s %s", FlagTypescriptRecording, boolToOnOff(data.TypescriptRecording)))
	if !isUpdate {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagAIThreatDetection, boolToOnOff(data.AIThreatDetection)))
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagAITerminateSessionOnDetection, boolToOnOff(data.AITerminateSessionOnDetection)))
	}

	env := data.Environment.ValueString()
	switch env {
	case EnvLocal:
		appendLocalNetworkFlags(parts, data.LocalNetwork)
	case EnvAWS:
		appendAwsFlags(parts, data.Aws)
	case EnvAzure:
		appendAzureFlags(parts, data.Azure)
	case EnvDomain:
		appendDomainFlags(parts, data.Domain)
	case EnvGCP:
		appendGcpFlags(parts, data.Gcp)
	}
}

// buildPamConfigNewCommand builds `pam config new ...` from the resource model.
// Required attributes (environment, title, gateway, application_folder) are enforced by the schema.
func buildPamConfigNewCommand(data *PamConfigurationResourceModel) string {
	parts := []string{CmdPamConfig, CmdPamNew}
	appendPamConfigBodyFlags(&parts, data, false)
	return strings.Join(parts, " ")
}

// buildPamConfigEditCommand builds `pam config edit '<uid>' ...` from the resource model.
func buildPamConfigEditCommand(uid string, data *PamConfigurationResourceModel) string {
	parts := []string{CmdPamConfig, CmdPamEdit, quoteShellSingle(strings.TrimSpace(uid))}
	appendPamConfigBodyFlags(&parts, data, true)
	return strings.Join(parts, " ")
}

func appendLocalNetworkFlags(parts *[]string, m *PamLocalNetworkModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagNetworkId, m.NetworkId)
	appendFlagString(parts, FlagNetworkCidr, m.NetworkCidr)
}

func appendAwsFlags(parts *[]string, m *PamAwsModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagAwsId, m.AwsId)
	appendFlagString(parts, FlagAccessKeyId, m.AccessKeyId)
	appendFlagString(parts, FlagAccessSecretKey, m.AccessSecretKey)
	for _, r := range sortedSetStrings(m.RegionNames) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagRegionName, quoteShellSingle(r)))
	}
}

func appendAzureFlags(parts *[]string, m *PamAzureModel) {
	if m == nil {
		return
	}
	appendFlagString(parts, FlagAzureId, m.AzureId)
	appendFlagString(parts, FlagClientId, m.ClientId)
	appendFlagString(parts, FlagClientSecret, m.ClientSecret)
	appendFlagString(parts, FlagSubscriptionId, m.SubscriptionId)
	appendFlagString(parts, FlagTenantId, m.TenantId)
	for _, rg := range sortedSetStrings(m.ResourceGroups) {
		*parts = append(*parts, fmt.Sprintf("%s %s", FlagResourceGroup, quoteShellSingle(rg)))
	}
}

func appendDomainFlags(parts *[]string, m *PamDomainModel) {
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

func appendGcpFlags(parts *[]string, m *PamGcpModel) {
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
	return strings.Join([]string{CmdPamConfig, CmdPamRemove, quoteShellSingle(strings.TrimSpace(uid))}, " ")
}

// buildPamConfigListCommand builds `pam config list --config '<uid>' --format json`.
func buildPamConfigListCommand(uid string) string {
	return strings.Join([]string{
		CmdPamConfig, CmdPamList,
		FlagPamListConfig, quoteShellSingle(strings.TrimSpace(uid)),
		FlagFormat, FormatJSON,
	}, " ")
}

// pamConfigTypeToEnvironment maps API config_type values to Terraform environment.
var pamConfigTypeToEnvironment = map[string]string{
	"pamAwsConfiguration":     EnvAWS,
	"pamAzureConfiguration":   EnvAzure,
	"pamDomainConfiguration":  EnvDomain,
	"pamGcpConfiguration":     EnvGCP,
	"pamNetworkConfiguration": EnvLocal,
}

func configTypeToEnvironment(ct string) (string, error) {
	ct = strings.TrimSpace(ct)
	if e, ok := pamConfigTypeToEnvironment[ct]; ok {
		return e, nil
	}
	return "", fmt.Errorf("unknown config_type %q", ct)
}

func pickApplicationFolderFromListAPI(sf utils.PamConfigListResponse, prior types.String) types.String {
	uid := strings.TrimSpace(sf.SharedFolder.UID)
	name := strings.TrimSpace(sf.SharedFolder.Name)
	if !prior.IsNull() && !prior.IsUnknown() {
		p := strings.TrimSpace(prior.ValueString())
		if p != "" {
			if p == uid {
				return types.StringValue(uid)
			}
			if p == name {
				return types.StringValue(name)
			}
		}
	}
	if uid != "" {
		return types.StringValue(uid)
	}
	return types.StringValue(name)
}

// MapPamConfigAPIResponseToModel refreshes fields returned by the list API. Other schema fields are unchanged.
func mapPamConfigAPIResponseToModel(state *PamConfigurationResourceModel, api *utils.PamConfigListResponse) error {
	env, err := configTypeToEnvironment(api.ConfigType)
	if err != nil {
		return err
	}
	state.Title = types.StringValue(api.Name)
	state.Environment = types.StringValue(env)
	state.Gateway = types.StringValue(strings.TrimSpace(api.GatewayUID))
	state.ApplicationFolder = pickApplicationFolderFromListAPI(*api, state.ApplicationFolder)
	return nil
}
