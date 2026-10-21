// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"errors"
	"fmt"
	"testing"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFetchPamUssConfigByNetworkUIDCommand(t *testing.T) {
	got := commonpamconfiguration.FetchPamUssConfigByNetworkUIDCommand("net-uid-1")
	want := "pam universal-sync-config list --network 'net-uid-1' --format json"
	if got != want {
		t.Errorf("command = %q, want %q", got, want)
	}
}

func TestIsUssConfigNotFoundError(t *testing.T) {
	// Real Commander wording for a network that never had a USS config created.
	wrapped := fmt.Errorf("Unable to get PAM universal secrets sync configuration: internal server error (500): No USS configuration found for this network")
	if !commonpamconfiguration.IsUssConfigNotFoundError(wrapped) {
		t.Errorf("expected IsUssConfigNotFoundError to match Commander's real error text, got false for: %v", wrapped)
	}

	// Case-insensitivity.
	if !commonpamconfiguration.IsUssConfigNotFoundError(errors.New("NO USS CONFIGURATION FOUND for network xyz")) {
		t.Error("expected case-insensitive match")
	}

	if commonpamconfiguration.IsUssConfigNotFoundError(errors.New("some unrelated failure")) {
		t.Error("expected no match for an unrelated error")
	}

	if commonpamconfiguration.IsUssConfigNotFoundError(nil) {
		t.Error("expected no match for a nil error")
	}
}

// gcpUssAPIResponse mirrors the exact example response for
// `pam universal-sync-config list --network <uid> --format json` against a GCP config.
func gcpUssAPIResponse() *utils.PamUniversalSyncConfigResponse {
	return &utils.PamUniversalSyncConfigResponse{
		RecordUID:     "08OV7gNVRky9BtStSsBGEw",
		RecordType:    "pamGcpConfiguration",
		Enabled:       true,
		DryRunEnabled: true,
		VaultName:     "N/A",
		SyncIdentity:  "abc",
		Folders: []utils.PamUniversalSyncConfigFolderResponse{
			{UID: "JOk6cztcISx6FduQTLo-hQ"},
			{UID: "zz4QAzzkPExqMbDkR4HjPg"},
		},
	}
}

func TestMapUssConfigAPIResponseToModel_Gcp(t *testing.T) {
	state := &commonpamconfiguration.PamConfigurationResourceModel{
		Environment: types.StringValue(commonpamconfiguration.EnvGCP),
		Gcp:         &commonpamconfiguration.PamGcpModel{},
	}

	commonpamconfiguration.MapUssConfigAPIResponseToModel(state, gcpUssAPIResponse(), false)

	uss := state.Gcp.UniversalSecretsSync
	if uss == nil {
		t.Fatal("expected Gcp.UniversalSecretsSync to be set")
	}
	if !uss.Enabled.ValueBool() {
		t.Error("expected enabled to be true")
	}
	if !uss.DryRun.ValueBool() {
		t.Error("expected dry_run to be true")
	}
	if uss.SyncIdentity.ValueString() != "abc" {
		t.Errorf("sync_identity = %q, want abc", uss.SyncIdentity.ValueString())
	}
	if uss.FolderLocation.IsNull() {
		t.Fatal("expected folder_location to be set")
	}
	elems := uss.FolderLocation.Elements()
	if len(elems) != 2 {
		t.Fatalf("expected 2 folder UIDs, got %d: %v", len(elems), elems)
	}

	// aws/azure must stay untouched.
	if state.Aws != nil {
		t.Error("expected Aws to remain nil")
	}
}

func TestMapUssConfigAPIResponseToModel_AzureSetsVaultName(t *testing.T) {
	state := &commonpamconfiguration.PamConfigurationResourceModel{
		Environment: types.StringValue(commonpamconfiguration.EnvAzure),
		Azure:       &commonpamconfiguration.PamAzureModel{},
	}
	api := &utils.PamUniversalSyncConfigResponse{
		RecordType:    "pamAzureConfiguration",
		Enabled:       true,
		DryRunEnabled: false,
		VaultName:     "my-vault",
		SyncIdentity:  "N/A",
	}

	commonpamconfiguration.MapUssConfigAPIResponseToModel(state, api, false)

	uss := state.Azure.UniversalSecretsSync
	if uss == nil {
		t.Fatal("expected Azure.UniversalSecretsSync to be set")
	}
	if uss.VaultName.ValueString() != "my-vault" {
		t.Errorf("vault_name = %q, want my-vault", uss.VaultName.ValueString())
	}
	if !uss.SyncIdentity.IsNull() {
		t.Errorf("expected sync_identity %q (\"N/A\") to map to null, got %v", "N/A", uss.SyncIdentity)
	}
}

// TestMapUssConfigAPIResponseToModel_PlaceholderResponseIsNoOp reproduces a real Commander
// response for a GCP PAM configuration that has never had Universal Secrets Sync created:
// `pam universal-sync-config list` still returns "status":"success" with
// enabled/dry_run_enabled false, vault_name/sync_identity "N/A", and no folders, instead of
// erroring. Without ussHasMaterialConfig, this was incorrectly treated as a real (disabled)
// config, producing a permanent `enabled: false -> null` / `dry_run: false -> null` diff on
// every plan even though the user never configured universal_secrets_sync at all.
func TestMapUssConfigAPIResponseToModel_PlaceholderResponseIsNoOp(t *testing.T) {
	state := &commonpamconfiguration.PamConfigurationResourceModel{
		Environment: types.StringValue(commonpamconfiguration.EnvGCP),
		Gcp:         &commonpamconfiguration.PamGcpModel{},
	}
	api := &utils.PamUniversalSyncConfigResponse{
		RecordType:    "pamGcpConfiguration",
		Enabled:       false,
		DryRunEnabled: false,
		VaultName:     "N/A",
		SyncIdentity:  "N/A",
	}

	commonpamconfiguration.MapUssConfigAPIResponseToModel(state, api, false)

	if state.Gcp.UniversalSecretsSync != nil {
		t.Errorf("expected UniversalSecretsSync to remain nil for Commander's never-configured placeholder response, got %+v", state.Gcp.UniversalSecretsSync)
	}
}

// An all-false response must be kept when the block was already in state, otherwise a
// user-configured `enabled = false, dry_run = false` block shows a permanent `+` diff.
func TestMapUssConfigAPIResponseToModel_AllFalseKeptWhenPriorConfigured(t *testing.T) {
	state := &commonpamconfiguration.PamConfigurationResourceModel{
		Environment: types.StringValue(commonpamconfiguration.EnvGCP),
		Gcp:         &commonpamconfiguration.PamGcpModel{},
	}
	api := &utils.PamUniversalSyncConfigResponse{
		RecordType:   "pamGcpConfiguration",
		VaultName:    "N/A",
		SyncIdentity: "N/A",
	}

	commonpamconfiguration.MapUssConfigAPIResponseToModel(state, api, true)

	uss := state.Gcp.UniversalSecretsSync
	if uss == nil {
		t.Fatal("expected UniversalSecretsSync to be set when prior state had it configured")
	}
	if uss.Enabled.ValueBool() || uss.DryRun.ValueBool() {
		t.Errorf("expected enabled/dry_run false, got enabled=%v dry_run=%v", uss.Enabled, uss.DryRun)
	}
	if !uss.SyncIdentity.IsNull() || !uss.FolderLocation.IsNull() {
		t.Errorf("expected sync_identity/folder_location null, got %v / %v", uss.SyncIdentity, uss.FolderLocation)
	}
}

func TestMapUssConfigAPIResponseToModel_NilApiIsNoOp(t *testing.T) {
	state := &commonpamconfiguration.PamConfigurationResourceModel{
		Environment: types.StringValue(commonpamconfiguration.EnvAWS),
		Aws:         &commonpamconfiguration.PamAwsModel{},
	}

	commonpamconfiguration.MapUssConfigAPIResponseToModel(state, nil, false)

	if state.Aws.UniversalSecretsSync != nil {
		t.Error("expected UniversalSecretsSync to remain nil when api is nil")
	}
}

func TestMapUssConfigAPIResponseToModel_UnsupportedEnvironmentIsNoOp(t *testing.T) {
	state := &commonpamconfiguration.PamConfigurationResourceModel{
		Environment: types.StringValue(commonpamconfiguration.EnvDomain),
		Domain:      &commonpamconfiguration.PamDomainModel{},
	}

	// Must not panic even though Domain has no UniversalSecretsSync field.
	commonpamconfiguration.MapUssConfigAPIResponseToModel(state, gcpUssAPIResponse(), false)

	if state.Aws != nil || state.Azure != nil || state.Gcp != nil {
		t.Error("expected no provider block to be populated for environment=domain")
	}
}
