// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"fmt"
	"strings"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// listToStrings returns non-empty string values from a types.List (e.g. user_groups, machine_collections).
// func listToStrings(l types.List) []string {
// 	if l.IsNull() || l.IsUnknown() {
// 		return nil
// 	}
// 	var out []string
// 	for _, el := range l.Elements() {
// 		if s, ok := el.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
// 			v := strings.TrimSpace(s.ValueString())
// 			if v != "" {
// 				out = append(out, v)
// 			}
// 		}
// 	}
// 	return out
// }

// setToStrings returns non-empty string values from a types.Set (e.g. control).
func setToStrings(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	var out []string
	for _, el := range s.Elements() {
		if str, ok := el.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
			v := strings.TrimSpace(str.ValueString())
			if v != "" {
				out = append(out, v)
			}
		}
	}
	return out
}

// appendListFlags appends multiple flag arguments for each value (e.g. --user-filter "id1" --user-filter "id2").
func appendListFlags(parts []string, flag string, values []string) []string {
	for _, v := range values {
		parts = append(parts, fmt.Sprintf("%s '%s'", flag, v))
	}
	return parts
}

// appendEpmPolicyAttributeFlags appends policy name, type, status, and optional list flags (shared by add and update).
func appendEpmPolicyAttributeFlags(parts []string, data *EpmPolicyResourceModel) []string {
	parts = append(parts, fmt.Sprintf("%s '%s'", commonepm.FlagPolicyName, data.PolicyName.ValueString()))
	parts = append(parts, fmt.Sprintf("%s '%s'", commonepm.FlagPolicyType, data.PolicyType.ValueString()))

	// If status is off, set enable to off otherwise set status
	if data.Status.ValueString() == commonepm.StatusOff {
		parts = append(parts, fmt.Sprintf("%s '%s'", commonepm.FlagEnable, commonepm.StatusOff))
	} else {
		parts = append(parts, fmt.Sprintf("%s '%s'", commonepm.FlagStatus, data.Status.ValueString()))
	}

	if data.Status.ValueString() == commonepm.StatusMonitorAndNotify {
		if !data.Message.IsNull() && !data.Message.IsUnknown() {
			msg := strings.TrimSpace(data.Message.ValueString())
			if msg != "" {
				parts = append(parts, fmt.Sprintf("%s '%s'", commonepm.FlagMessage, msg))
			}
		}
		if !data.RequirePolicyAcknowledgement.IsNull() && !data.RequirePolicyAcknowledgement.IsUnknown() {
			ack := utils.ValueOff
			if data.RequirePolicyAcknowledgement.ValueBool() {
				ack = utils.ValueOn
			}
			parts = append(parts, fmt.Sprintf("%s '%s'", commonepm.FlagRequireAcknowledgement, ack))
		}
	}

	parts = appendListFlags(parts, commonepm.FlagControl, setToStrings(data.Control))
	parts = appendListFlags(parts, commonepm.FlagDayFilter, setToStrings(data.DayFilter))
	parts = appendListFlags(parts, commonepm.FlagUserFilter, setToStrings(data.UserGroups))
	parts = appendListFlags(parts, commonepm.FlagMachineFilter, setToStrings(data.MachineCollections))
	parts = appendListFlags(parts, commonepm.FlagAppFilter, setToStrings(data.Applications))
	parts = appendListFlags(parts, commonepm.FlagTimeFilter, setToStrings(data.TimeFilter))
	parts = appendListFlags(parts, commonepm.FlagDateFilter, setToStrings(data.DateFilter))
	return parts
}

// buildCreateCommand builds the "epm policy add" command with flags from the resource model.
func buildCreateCommand(data *EpmPolicyResourceModel) string {
	var parts []string
	parts = append(parts, commonepm.CmdEpmPolicyAdd)
	return strings.Join(appendEpmPolicyAttributeFlags(parts, data), " ")
}

// buildUpdateCommand builds the "epm policy update <id>" command with the same flags as add.
func buildUpdateCommand(policyID string, data *EpmPolicyResourceModel) string {
	var parts []string
	parts = append(parts, commonepm.CmdEpmPolicyEdit, policyID)
	return strings.Join(appendEpmPolicyAttributeFlags(parts, data), " ")
}

// buildViewCommand builds `epm policy view <id> --format json`.
func buildViewCommand(policyID string) string {
	return commonepm.BuildPolicyViewCommand(policyID)
}
