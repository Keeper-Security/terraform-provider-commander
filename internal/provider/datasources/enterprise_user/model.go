// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseUserDataSourceModel struct {
	User           types.String `tfsdk:"user"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	JobTitle       types.String `tfsdk:"job_title"`
	Roles          types.Set    `tfsdk:"roles"`
	Teams          types.Set    `tfsdk:"teams"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
