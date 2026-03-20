// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - User only: user email or user ID (e.g. user@example.com or 1326075447607317)
//   - With managed company: "managed_company_name_or_id,user_email_or_id" (comma-separated, e.g. "Test Company,1326075447607317" or "1169425105420462,user@example.com")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	resourceID, managedCompany, diags := utils.ParseManagedCompanyImportID(req.ID, "user")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := EnterpriseUserResourceModel{
		Id:             types.StringValue(resourceID),
		Email:          types.StringNull(),
		Name:           types.StringNull(),
		JobTitle:       types.StringNull(),
		Roles:          types.SetNull(types.StringType),
		Teams:          types.SetNull(types.StringType),
		Node:           types.StringNull(),
		ManagedCompany: types.StringNull(),
		Status:         types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
