// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"fmt"
	"strings"
)

// buildCreateCommand builds the scim create command: scim create --node NODE --prefix PREFIX --unique-groups on|off.
func buildCreateCommand(data *EnterpriseScimResourceModel) string {
	uniqueGroups := ValueOff
	if data.UniqueGroups.ValueBool() {
		uniqueGroups = ValueOn
	}
	return strings.Join([]string{
		CmdScimCreate,
		fmt.Sprintf("%s '%s'", FlagNode, data.Node.ValueString()),
		fmt.Sprintf("%s '%s'", FlagPrefix, data.Prefix.ValueString()),
		fmt.Sprintf("%s '%s'", FlagUniqueGroups, uniqueGroups),
	}, " ")
}

// buildUpdateCommand builds the scim edit command: scim edit ID --prefix PREFIX --unique-groups on|off.
func buildUpdateCommand(scimId string, data *EnterpriseScimResourceModel) string {
	uniqueGroups := ValueOff
	if data.UniqueGroups.ValueBool() {
		uniqueGroups = ValueOn
	}
	return strings.Join([]string{
		CmdScimEdit,
		scimId,
		fmt.Sprintf("%s '%s'", FlagPrefix, data.Prefix.ValueString()),
		fmt.Sprintf("%s '%s'", FlagUniqueGroups, uniqueGroups),
	}, " ")
}
