// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseScimPushResourceModel struct {
	Id              types.String `tfsdk:"id"`
	ScimId          types.String `tfsdk:"scim_id"`
	Source          types.String `tfsdk:"source"`
	Record          types.String `tfsdk:"record"`
	AutoApprove     types.Bool   `tfsdk:"auto_approve"`
	ManagedCompany  types.String `tfsdk:"managed_company"`
}
