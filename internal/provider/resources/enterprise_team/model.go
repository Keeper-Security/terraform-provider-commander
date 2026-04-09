// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseTeamResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	RestrictEdit   types.Bool   `tfsdk:"restrict_record_edit"`
	RestrictShare  types.Bool   `tfsdk:"restrict_record_re_share"`
	RestrictView   types.Bool   `tfsdk:"enable_privacy_screen"`
	Users          types.Set    `tfsdk:"users"`
	Roles          types.Set    `tfsdk:"roles"`
	Node           types.String `tfsdk:"node"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
