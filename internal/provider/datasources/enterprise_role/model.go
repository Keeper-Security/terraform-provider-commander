// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseRoleDataSourceModel struct {
	Role           types.String `tfsdk:"role"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Users          types.Set    `tfsdk:"users"`
	Teams          types.Set    `tfsdk:"teams"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
