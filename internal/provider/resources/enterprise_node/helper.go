// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"fmt"
	"strings"
)

func buildEnterpriseNodeAddCommand(data EnterpriseNodeResourceModel) string {
	var parts []string

	parts = append(parts, "enterprise-node")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.Parent.IsNull() {
		value := data.Parent.ValueString()

		// We need to make the parent as "root" as if the parent is the same as the managed company bec. this like functionality implemente in commander cli.
		if !data.ManagedCompany.IsNull() && data.Parent.ValueString() == data.ManagedCompany.ValueString() {
			value = "root"
		}

		parts = append(parts, fmt.Sprintf("--parent '%s'", value))
	}

	// TODO: Currently its not working in
	// if !data.WipeOut.IsNull() {
	// 	parts = append(parts, "--wipe-out")
	// }

	// TODO: NEED TO CHECK IF THIS FLAG IS REQUIRED / usecase CHECK?
	if !data.ToggleIsolated.IsNull() {
		parts = append(parts, "--toggle-isolated")
	}

	// TODO: NEED TO CHECK HOW WE CAN ADD LOGO FILE - NOT WORKING IN COMMANDER CLI
	// if !data.LogoFile.IsNull() {
	// 	parts = append(parts, fmt.Sprintf("--logo-file '%s'", data.LogoFile.ValueString()))
	// }

	return strings.Join(parts, " ")
}

func buildEnterpriseNodeUpdateCommand(plan *EnterpriseNodeResourceModel, state *EnterpriseNodeResourceModel) string {
	var parts []string

	parts = append(parts, "enterprise-node")

	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	if !state.Parent.Equal(plan.Parent) {
		value := plan.Parent.ValueString()

		// We need to make the parent as "root" as we if the parent is the same as the managed company
		if !plan.ManagedCompany.IsNull() && plan.Parent.ValueString() == plan.ManagedCompany.ValueString() {
			value = "root"
		}
		parts = append(parts, fmt.Sprintf("--parent '%s'", value))
	}

	// if !state.WipeOut.Equal(plan.WipeOut) {
	// 	parts = append(parts, "--wipe-out")
	// }

	// if !state.ToggleIsolated.Equal(plan.ToggleIsolated) {
	// 	parts = append(parts, "--toggle-isolated")
	// }

	// if !state.LogoFile.Equal(plan.LogoFile) {
	// 	parts = append(parts, fmt.Sprintf("--logo-file '%s'", plan.LogoFile.ValueString()))
	// }

	// Node to update
	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	return strings.Join(parts, " ")
}
