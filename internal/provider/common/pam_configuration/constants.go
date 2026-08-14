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
	CmdPamConfig = "pam config"
	CmdPamNew    = "new"
	CmdPamEdit   = "edit"
	CmdPamRemove = "remove"
	CmdPamList   = "list"
)

const (
	FlagPamListConfig = "--config"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryReadFailed = "Read PAM Configuration Failed"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpFetchPamConfig = "Unable to fetch PAM configuration"
)
