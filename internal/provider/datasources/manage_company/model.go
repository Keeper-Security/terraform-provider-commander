// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManageCompanyDataSourceModel struct {
	ManagedCompany types.String `tfsdk:"managed_company"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Node           types.String `tfsdk:"node"`
	Plan           types.String `tfsdk:"plan"`
	FilePlan       types.String `tfsdk:"file_plan"`
}
