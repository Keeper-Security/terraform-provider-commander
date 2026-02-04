// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManageCompanyDataSourceModel struct {
	// Input fields (optional)
	Id   types.Number `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Output fields
	Node     types.String `tfsdk:"node"`
	Plan     types.String `tfsdk:"plan"`
	FilePlan types.String `tfsdk:"file_plan"`
}
