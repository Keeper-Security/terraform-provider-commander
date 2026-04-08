// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import "github.com/hashicorp/terraform-plugin-framework/types"

type EpmPolicyResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	ManagedCompany     types.String `tfsdk:"managed_company"`
	PolicyName         types.String `tfsdk:"policy_name"`
	PolicyType         types.String `tfsdk:"policy_type"`
	Status             types.String `tfsdk:"status"`
	Control            types.Set    `tfsdk:"control"`
	UserGroups         types.Set    `tfsdk:"user_groups"`
	MachineCollections types.Set    `tfsdk:"machine_collections"`
	Applications       types.Set    `tfsdk:"applications"`
	DayFilter          types.Set    `tfsdk:"day_filter"`
	TimeFilter         types.Set    `tfsdk:"time_filter"`
	DateFilter         types.Set    `tfsdk:"date_filter"`
}
