// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseScimResourceModel struct {
	Id                types.String `tfsdk:"id"`
	ScimURL           types.String `tfsdk:"scim_url"`
	Node              types.String `tfsdk:"node"`
	Status            types.String `tfsdk:"status"`
	Prefix            types.String `tfsdk:"prefix"`
	UniqueGroups      types.Bool   `tfsdk:"unique_groups"`
	ProvisioningToken types.String `tfsdk:"provisioning_token"`
	ManagedCompany    types.String `tfsdk:"managed_company"`
}
