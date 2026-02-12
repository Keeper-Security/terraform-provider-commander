// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseNodeResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Parent types.String `tfsdk:"parent"`
	// WipeOut        types.Bool   `tfsdk:"wipe_out"` // NOTE: usecase dont fit - this falg removes all users/roles/teams/etc from the node, if there is all users,nodes,..etc then we need to remove manually from state
	ToggleIsolated types.Bool `tfsdk:"toggle_isolated"`
	// LogoFile       types.String `tfsdk:"logo_file"` // NOTE: In commander cli not working and in admin console there no feature like this
	ManagedCompany types.String `tfsdk:"managed_company"`
}
