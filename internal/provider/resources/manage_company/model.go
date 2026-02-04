// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManageCompanyResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Node     types.String `tfsdk:"node"`
	Plan     types.String `tfsdk:"plan"`
	Seats    types.Int64  `tfsdk:"seats"`
	FilePlan types.String `tfsdk:"file_plan"`
	AddOns   types.Set    `tfsdk:"add_ons"`
}
