// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

const (
	EnvLocal  = "local"
	EnvAWS    = "aws"
	EnvAzure  = "azure"
	EnvGCP    = "gcp"
	EnvDomain = "domain"
)

// Commander CLI: `pam config new` / `pam config edit` and flags.
const (
	CmdPamConfig                    = "pam config"
	CmdPamNew                       = "new"
	CmdPamEdit                      = "edit"
	CmdPamRemove                    = "remove"
	CmdPamList                      = "list"
	CmdPamUniversalSyncConfigList   = "pam universal-sync-config list"
	CmdPamUniversalSyncConfigNew    = "pam universal-sync-config add"
	CmdPamUniversalSyncConfigEdit   = "pam universal-sync-config edit"
	CmdPamUniversalSyncConfigRemove = "pam universal-sync-config remove"
)

const (
	FlagPamListConfig = "--config"

	// Flags for `pam universal-sync-config`.
	FlagPamUssNetwork      = "--network"
	FlagPamUssEnabled      = "--enabled"
	FlagPamUssFolder       = "--folder"
	FlagPamUssSyncIdentity = "--sync-identity"
	FlagPamUssVaultName    = "--vault-name"
	FlagPamUssDryRun       = "--dry-run"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryReadFailed = "Read PAM Configuration Failed"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpFetchPamConfig    = "Unable to fetch PAM configuration"
	ErrOpFetchPamUssConfig = "Unable to get PAM universal secrets sync configuration"
)
