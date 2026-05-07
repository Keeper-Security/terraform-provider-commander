// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

// Allowed values for source attribute.
const (
	SourceGoogle = "google"
	SourceAD     = "ad"
	SourceRecord = "record"
)

// Commander CLI command for SCIM push.
const CmdScimPush = "scim push"

// Command flags.
const (
	FlagSource      = "--source"
	FlagRecord      = "--record"
	FlagAutoApprove = "--auto-approve"
)

// Auto-approve flag values for the CLI.
const (
	AutoApproveOn  = "on"
	AutoApproveOff = "off"
)
