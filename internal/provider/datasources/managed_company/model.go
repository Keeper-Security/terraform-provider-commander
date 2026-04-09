// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManagedCompanyDataSourceModel struct {
	ManagedCompany types.String `tfsdk:"managed_company"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Node           types.String `tfsdk:"node"`
	NodeName       types.String `tfsdk:"node_name"`
	Plan           types.String `tfsdk:"plan"`
	FilePlan       types.String `tfsdk:"file_plan"`
	Seats          types.Int64  `tfsdk:"seats"`
	AddOns         types.Set    `tfsdk:"add_ons"`
}
