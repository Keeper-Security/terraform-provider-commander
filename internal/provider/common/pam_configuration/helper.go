// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// buildPamConfigListCommand builds `pam config list --config '<uid>' --format json -v`.
func FetchPamConfigByUIDCommand(uid string) string {
	return strings.Join([]string{
		CmdPamConfig, CmdPamList,
		FlagPamListConfig, utils.QuoteShellSingle(strings.TrimSpace(uid)),
		utils.FlagFormatJSON,
		utils.FlagVerbose,
	}, " ")
}

// FetchPamUssConfigByNetworkUIDCommand builds
// `pam universal-sync-config list --network '<uid>' --format json`.
func FetchPamUssConfigByNetworkUIDCommand(uid string) string {
	return strings.Join([]string{
		CmdPamUniversalSyncConfigList,
		FlagPamUssNetwork, utils.QuoteShellSingle(strings.TrimSpace(uid)),
		utils.FlagFormatJSON,
	}, " ")
}

// BuildPamRemoveUssConfigCommand builds `pam universal-sync-config remove --network '<uid>'`.
func BuildPamRemoveUssConfigCommand(uid string) string {
	return strings.Join([]string{
		CmdPamUniversalSyncConfigRemove,
		FlagPamUssNetwork, utils.QuoteShellSingle(strings.TrimSpace(uid)),
	}, " ")
}

// ussConfigNotFoundPhrase is the substring (case-insensitive) Commander's
// `pam universal-sync-config list` returns in its error when no Universal Secrets Sync
// configuration has ever been created for the given network/PAM configuration. This is a
// normal, expected condition - not an out-of-band deletion of the PAM configuration itself -
// so callers must treat it as "leave UniversalSecretsSync nil", not as api.ErrResourceNotFound
// (whose usual handler removes the whole resource from state) or a hard Read failure. It
// doesn't contain "not found"/"cannot find", so the generic not-found classifier in the api
// package doesn't catch it.
const ussConfigNotFoundPhrase = "no uss configuration found"

// IsUssConfigNotFoundError reports whether err is Commander's "no Universal Secrets Sync
// configuration exists yet for this network" response.
func IsUssConfigNotFoundError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), ussConfigNotFoundPhrase)
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

func pickGatewayFromListAPI(gatewayUID, gatewayName string, prior types.String) types.String {
	uid := strings.TrimSpace(gatewayUID)
	name := strings.TrimSpace(gatewayName)
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

// firstFieldValue returns the first element of a string slice, trimmed, or "" if empty.
func firstFieldValue(vals []string) string {
	if len(vals) == 0 {
		return ""
	}
	return strings.TrimSpace(vals[0])
}

// setStringFromField sets a types.String from a field slice; null if empty.
func setStringFromField(vals []string) types.String {
	v := firstFieldValue(vals)
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// newlineToStringSet splits a newline-delimited field value into a Terraform string set.
func newlineToStringSet(vals []string) types.Set {
	raw := firstFieldValue(vals)
	if raw == "" {
		return types.SetNull(types.StringType)
	}
	var elems []attr.Value
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			elems = append(elems, types.StringValue(line))
		}
	}
	if len(elems) == 0 {
		return types.SetNull(types.StringType)
	}
	sv, _ := types.SetValue(types.StringType, elems)
	return sv
}

// parseDomainBool converts "True"/"False" (case-insensitive) field values to types.Bool.
func parseDomainBool(vals []string) types.Bool {
	v := strings.ToLower(firstFieldValue(vals))
	switch v {
	case "true":
		return types.BoolValue(true)
	case "false":
		return types.BoolValue(false)
	default:
		return types.BoolNull()
	}
}

// parseDomainHostnamePort splits "hostname:port" into separate hostname and port strings.
func parseDomainHostnamePort(vals []string) (types.String, types.String) {
	raw := firstFieldValue(vals)
	if raw == "" {
		return types.StringNull(), types.StringNull()
	}
	idx := strings.LastIndex(raw, ":")
	if idx < 0 || idx == len(raw)-1 {
		return types.StringValue(raw), types.StringNull()
	}
	return types.StringValue(strings.TrimSpace(raw[:idx])), types.StringValue(strings.TrimSpace(raw[idx+1:]))
}

func mapFieldsToLocalNetwork(f *utils.PamConfigFieldsResponse) *PamLocalNetworkModel {
	return &PamLocalNetworkModel{
		NetworkId:   setStringFromField(f.NetworkId),
		NetworkCidr: setStringFromField(f.NetworkCIDR),
	}
}

func mapFieldsToAws(f *utils.PamConfigFieldsResponse) *PamAwsModel {
	return &PamAwsModel{
		AwsId:           setStringFromField(f.AwsId),
		AccessKeyId:     setStringFromField(f.AccessKeyId),
		AccessSecretKey: setStringFromField(f.AccessSecretKey),
		RegionNames:     newlineToStringSet(f.RegionNames),
	}
}

func mapFieldsToAzure(f *utils.PamConfigFieldsResponse) *PamAzureModel {
	return &PamAzureModel{
		AzureId:        setStringFromField(f.AzureId),
		ClientId:       setStringFromField(f.ClientId),
		ClientSecret:   setStringFromField(f.ClientSecret),
		SubscriptionId: setStringFromField(f.SubscriptionId),
		TenantId:       setStringFromField(f.TenantId),
		ResourceGroups: newlineToStringSet(f.ResourceGroups),
	}
}

func mapFieldsToDomain(f *utils.PamConfigFieldsResponse, domainAdministrativeCredential string) *PamDomainModel {
	hostname, port := parseDomainHostnamePort(f.PamHostname)
	return &PamDomainModel{
		DomainId:          setStringFromField(f.PamDomainId),
		DomainHostname:    hostname,
		DomainPort:        port,
		DomainUseSsl:      parseDomainBool(f.UseSSL),
		DomainScanDcCidr:  parseDomainBool(f.ScanDCCIDR),
		DomainNetworkCidr: setStringFromField(f.NetworkCIDR),
		DomainAdmin:       types.StringValue(domainAdministrativeCredential),
		UserMatch:         setStringFromField(f.UserMatch),
	}
}

func mapFieldsToGcp(f *utils.PamConfigFieldsResponse) *PamGcpModel {
	return &PamGcpModel{
		GcpId:             setStringFromField(f.PamGcpId),
		ServiceAccountKey: setStringFromField(f.PamServiceAccountKey),
		GoogleAdminEmail:  setStringFromField(f.PamGoogleAdminEmail),
		GcpRegion:         setStringFromField(f.PamGcpRegionName),
	}
}

// MapPamConfigAPIResponseToModel refreshes fields returned by the list API. Other schema fields are unchanged.
func MapPamConfigAPIResponseToModel(state *PamConfigurationResourceModel, api *utils.PamConfigListResponse) error {
	env, err := configTypeToEnvironment(api.ConfigType)
	if err != nil {
		return err
	}
	state.Title = types.StringValue(api.Name)
	state.Environment = types.StringValue(env)
	state.Gateway = pickGatewayFromListAPI(api.GatewayUID, api.GatewayName, state.Gateway)
	state.ApplicationFolder = pickApplicationFolderFromListAPI(*api, state.ApplicationFolder)

	if api.AllowedSettings != nil {
		state.Connections = types.BoolValue(api.AllowedSettings.Connections)
		state.Tunneling = types.BoolValue(api.AllowedSettings.Tunneling)
		state.Rotation = types.BoolValue(api.AllowedSettings.Rotation)
		state.RemoteBrowserIsolation = types.BoolValue(api.AllowedSettings.RemoteBrowserIsolation)
		state.ConnectionsRecording = types.BoolValue(api.AllowedSettings.ConnectionsRecording)
		state.TypescriptRecording = types.BoolValue(api.AllowedSettings.TypescriptRecording)
		state.AIThreatDetection = types.BoolValue(api.AllowedSettings.AIThreatDetection)
		state.AITerminateSessionOnDetection = types.BoolValue(api.AllowedSettings.AITerminateSessionOnDetection)
	}

	if f := api.Fields; f != nil {
		state.Schedule = setStringFromField(f.DefaultSchedule)
		state.PortMapping = newlineToStringSet(f.PortMapping)

		state.LocalNetwork = nil
		state.Aws = nil
		state.Azure = nil
		state.Domain = nil
		state.Gcp = nil

		switch env {
		case EnvLocal:
			state.LocalNetwork = mapFieldsToLocalNetwork(f)
		case EnvAWS:
			state.Aws = mapFieldsToAws(f)
		case EnvAzure:
			state.Azure = mapFieldsToAzure(f)
		case EnvDomain:
			state.Domain = mapFieldsToDomain(f, api.DomainAdministrativeCredential)
		case EnvGCP:
			state.Gcp = mapFieldsToGcp(f)
		}
	}

	return nil
}

// ussStringOrNull returns null for "" and the API's "N/A" placeholder (used when a field
// doesn't apply to this record's type), else the trimmed value.
func ussStringOrNull(s string) types.String {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "N/A") {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// ussFolderLocationSet converts the "folders" array from
// `pam universal-sync-config list` into a set of folder UIDs.
func ussFolderLocationSet(folders []utils.PamUniversalSyncConfigFolderResponse) types.Set {
	elems := make([]attr.Value, 0, len(folders))
	for _, f := range folders {
		if uid := strings.TrimSpace(f.UID); uid != "" {
			elems = append(elems, types.StringValue(uid))
		}
	}
	if len(elems) == 0 {
		return types.SetNull(types.StringType)
	}
	sv, _ := types.SetValue(types.StringType, elems)
	return sv
}

// ussHasMaterialConfig reports whether api represents an actually-configured Universal
// Secrets Sync setup, as opposed to Commander's placeholder response for a PAM
// configuration that has never had one created: `pam universal-sync-config list` returns
// "status":"success" with enabled/dry_run_enabled false, vault_name/sync_identity "N/A",
// and no folders even when nothing was ever configured, rather than erroring - so that
// shape must be treated the same as "not configured", not as real (if unusually blank)
// config.
func ussHasMaterialConfig(api *utils.PamUniversalSyncConfigResponse) bool {
	if api.Enabled || api.DryRunEnabled || len(api.Folders) > 0 {
		return true
	}
	return !ussStringOrNull(api.SyncIdentity).IsNull() || !ussStringOrNull(api.VaultName).IsNull()
}

// FetchAndMapUssConfig runs `pam universal-sync-config list` for networkUID and maps the result
// into state's active environment block. No-op for environments without USS, and when Commander
// reports no USS configuration exists. See MapUssConfigAPIResponseToModel for priorConfigured.
func FetchAndMapUssConfig(ctx context.Context, am *api.ApiManager, state *PamConfigurationResourceModel, networkUID string, priorConfigured bool) error {
	switch state.Environment.ValueString() {
	case EnvAWS, EnvAzure, EnvGCP:
	default:
		return nil
	}

	resp, err := am.ExecuteCommand(ctx, FetchPamUssConfigByNetworkUIDCommand(networkUID), ErrOpFetchPamUssConfig)
	if err != nil {
		if IsUssConfigNotFoundError(err) || errors.Is(err, api.ErrResourceNotFound) {
			return nil
		}
		return err
	}
	if resp == nil || resp.Data == nil {
		return nil
	}

	var data utils.PamUniversalSyncConfigResponse
	if err := utils.UnmarshalApiResponse(resp.Data, &data); err != nil {
		return err
	}
	MapUssConfigAPIResponseToModel(state, &data, priorConfigured)
	return nil
}

// MapUssConfigAPIResponseToModel sets the active environment's universal_secrets_sync block
// from `pam universal-sync-config list --network <uid> --format json`. No-op (leaving
// UniversalSecretsSync nil) if api is nil, state.Environment doesn't support USS (aws, azure,
// gcp), or api lacks material config (see ussHasMaterialConfig) and priorConfigured is false.
// priorConfigured must reflect prior state: an all-defaults response is indistinguishable
// from a user-configured `enabled = false, dry_run = false` block, so prior state decides.
func MapUssConfigAPIResponseToModel(state *PamConfigurationResourceModel, api *utils.PamUniversalSyncConfigResponse, priorConfigured bool) {
	if api == nil || (!priorConfigured && !ussHasMaterialConfig(api)) {
		return
	}

	common := PamCommonUniversalSecretsSyncModel{
		Enabled:        types.BoolValue(api.Enabled),
		DryRun:         types.BoolValue(api.DryRunEnabled),
		SyncIdentity:   ussStringOrNull(api.SyncIdentity),
		FolderLocation: ussFolderLocationSet(api.Folders),
	}

	switch state.Environment.ValueString() {
	case EnvAWS:
		if state.Aws == nil {
			return
		}
		state.Aws.UniversalSecretsSync = &PamAwsUniversalSecretsSyncModel{PamCommonUniversalSecretsSyncModel: common}
	case EnvAzure:
		if state.Azure == nil {
			return
		}
		state.Azure.UniversalSecretsSync = &PamAzureUniversalSecretsSyncModel{
			PamCommonUniversalSecretsSyncModel: common,
			VaultName:                          ussStringOrNull(api.VaultName),
		}
	case EnvGCP:
		if state.Gcp == nil {
			return
		}
		state.Gcp.UniversalSecretsSync = &PamGcpUniversalSecretsSyncModel{PamCommonUniversalSecretsSyncModel: common}
	}
}
