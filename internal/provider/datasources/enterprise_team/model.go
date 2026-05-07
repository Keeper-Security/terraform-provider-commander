// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseTeamDataSourceModel struct {
	Team           types.String `tfsdk:"team"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Users          types.Set    `tfsdk:"users"`
	Roles          types.Set    `tfsdk:"roles"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
