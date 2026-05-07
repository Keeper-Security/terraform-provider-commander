// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID format:
//   - scim_id (e.g. "1169425105420640")
//   - managed_company,scim_id (e.g. "My MC,1169425105420640")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseScimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	resourceID, managedCompany, diags := utils.ParseManagedCompanyImportID(req.ID, ImportResourceType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := EnterpriseScimResourceModel{
		Id:             types.StringValue(resourceID),
		ScimURL:        types.StringNull(),
		Node:           types.StringNull(),
		Status:         types.StringNull(),
		Prefix:         types.StringNull(),
		UniqueGroups:   types.BoolNull(),
		ManagedCompany: types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
