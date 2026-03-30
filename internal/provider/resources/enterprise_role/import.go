// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - Role only: role name or role ID (e.g. "Admin" or 1234567890)
//   - With managed company: "managed_company_name_or_id,role_name_or_id" (comma-separated, e.g. "Test Company,Admin" or "1169425105420462,1234567890")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	resourceID, managedCompany, diags := utils.ParseManagedCompanyImportID(req.ID, "role")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := EnterpriseRoleResourceModel{
		Id:                  types.StringValue(resourceID),
		Name:                types.StringNull(),
		Node:                types.StringNull(),
		Users:               types.SetNull(types.StringType),
		Teams:               types.SetNull(types.StringType),
		ManagingNodes:       types.MapNull(utils.ManagingNodesMapElemType),
		EnforcementPolicies: types.MapNull(types.StringType),
		ManagedCompany:      types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
