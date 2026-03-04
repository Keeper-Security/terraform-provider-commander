// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// mapScimReadResponseToModel sets the resource model from the API read response.
func mapScimReadResponseToModel(scim *utils.EnterpriseScimResponse, state *EnterpriseScimResourceModel) {
	state.Id = types.StringValue(strconv.Itoa(scim.ScimID))
	state.ScimURL = types.StringValue(scim.ScimURL)
	// Use node_name for display; API also has node_id
	state.Node = types.StringValue(state.Node.ValueString())
	state.Status = types.StringValue(scim.Status)
	state.Prefix = types.StringValue(scim.Prefix)
	state.UniqueGroups = types.BoolValue(scim.UniqueGroups)
}
