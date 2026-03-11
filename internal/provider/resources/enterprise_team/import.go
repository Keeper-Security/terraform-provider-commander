// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - Team only: team name or team UID (e.g. "Engineering" or 1234567890123456)
//   - With managed company: "managed_company_name_or_id,team_name_or_uid" (comma-separated, e.g. "Test Company,Engineering" or "1169425105420462,1234567890123456")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseTeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	resourceID, managedCompany, diags := utils.ParseManagedCompanyImportID(req.ID, "team")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := EnterpriseTeamResourceModel{
		Id:             types.StringValue(resourceID),
		Name:           types.StringNull(),
		RestrictEdit:   types.BoolNull(),
		RestrictShare:  types.BoolNull(),
		RestrictView:   types.BoolNull(),
		Users:          types.SetNull(types.StringType),
		Roles:          types.SetNull(types.StringType),
		Node:           types.StringNull(),
		ManagedCompany: types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
