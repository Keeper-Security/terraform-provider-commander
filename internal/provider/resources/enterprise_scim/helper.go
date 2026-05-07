// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"fmt"
	"strings"
)

// buildCreateCommand builds the scim create command: scim create --node NODE --prefix PREFIX --unique-groups on|off.
func buildCreateCommand(data *EnterpriseScimResourceModel) string {

	var parts []string
	parts = append(parts, CmdScimCreate)
	parts = append(parts, fmt.Sprintf("%s '%s'", FlagNode, data.Node.ValueString()))

	if !data.Prefix.IsUnknown() && !data.Prefix.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", FlagPrefix, data.Prefix.ValueString()))
	}

	uniqueGroups := ValueOff
	if data.UniqueGroups.ValueBool() {
		uniqueGroups = ValueOn
	}
	parts = append(parts, fmt.Sprintf("%s '%s'", FlagUniqueGroups, uniqueGroups))

	return strings.Join(parts, " ")
}

// buildUpdateCommand builds the scim edit command: scim edit ID --prefix PREFIX --unique-groups on|off.
func buildUpdateCommand(scimId string, data *EnterpriseScimResourceModel) string {
	var parts []string
	parts = append(parts, CmdScimEdit)
	parts = append(parts, scimId)

	if !data.Prefix.IsUnknown() && !data.Prefix.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", FlagPrefix, data.Prefix.ValueString()))
	} else {
		parts = append(parts, fmt.Sprintf("%s '%s'", FlagPrefix, ""))
	}

	uniqueGroups := ValueOff
	if data.UniqueGroups.ValueBool() {
		uniqueGroups = ValueOn
	}
	parts = append(parts, fmt.Sprintf("%s '%s'", FlagUniqueGroups, uniqueGroups))

	return strings.Join(parts, " ")
}
