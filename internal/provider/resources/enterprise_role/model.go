// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

// THIS FILES STORE TERRAFORM STATE STRUCTS

package enterpriserole

import "github.com/hashicorp/terraform-plugin-framework/types"

// ManagingNodeModel represents a single managing node with its privileges and cascade option.
// Note: The node name/ID is the map key, so it's not stored in this struct.
type ManagingNodeModel struct {
	Privileges types.Set  `tfsdk:"privileges"`
	Cascade    types.Bool `tfsdk:"cascade"`
}

type EnterpriseRoleResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Node                types.String `tfsdk:"node"`
	Users               types.Set    `tfsdk:"users"`
	Teams               types.Set    `tfsdk:"teams"`
	ManagingNodes       types.Map    `tfsdk:"managing_nodes"`
	EnforcementPolicies types.Map    `tfsdk:"enforcement_policies"`
	ManagedCompany      types.String `tfsdk:"managed_company"`
}
