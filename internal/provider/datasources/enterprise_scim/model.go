// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseScimDataSourceModel struct {
	Scim           types.String `tfsdk:"scim"`
	ScimId         types.String `tfsdk:"scim_id"`
	ScimUrl        types.String `tfsdk:"scim_url"`
	NodeId         types.String `tfsdk:"node_id"`
	NodeName       types.String `tfsdk:"node_name"`
	Status         types.String `tfsdk:"status"`
	Prefix         types.String `tfsdk:"prefix"`
	UniqueGroups   types.Bool   `tfsdk:"unique_groups"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
