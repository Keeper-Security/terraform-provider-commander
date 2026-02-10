// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - Node only: node name or node ID (e.g. "Root" or 1169425105420462)
//   - With managed company: "managed_company_name_or_id,node_name_or_id" (comma-separated, e.g. "Test Company,Root" or "1169425105420462,1169425105420462")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	resourceID, managedCompany, diags := utils.ParseManagedCompanyImportID(req.ID, "node")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := EnterpriseNodeResourceModel{
		Id:             types.StringValue(resourceID),
		Name:           types.StringNull(),
		Parent:         types.StringNull(),
		ToggleIsolated: types.BoolNull(),
		ManagedCompany: types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
