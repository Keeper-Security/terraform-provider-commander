// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseNodeResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Parent         types.String `tfsdk:"parent"`
	ToggleIsolated types.Bool   `tfsdk:"toggle_isolated"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
