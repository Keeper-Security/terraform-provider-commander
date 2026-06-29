// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID format:
//   - policy_id (e.g. "12345")
//   - managed_company,policy_id (e.g. "My MC,12345")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EpmPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	state := EpmPolicyResourceModel{
		Id:                           types.StringValue(resourceID),
		ManagedCompany:               types.StringNull(),
		PolicyName:                   types.StringNull(),
		PolicyType:                   types.StringNull(),
		Status:                       types.StringNull(),
		Message:                      types.StringNull(),
		RequirePolicyAcknowledgement: types.BoolNull(),
		Control:                      types.SetNull(types.StringType),
		UserGroups:                   types.SetNull(types.StringType),
		MachineCollections:           types.SetNull(types.StringType),
		Applications:                 types.SetNull(types.StringType),
		DayFilter:                    types.SetValueMust(types.StringType, []attr.Value{}),
		TimeFilter:                   types.SetValueMust(types.StringType, []attr.Value{}),
		DateFilter:                   types.SetValueMust(types.StringType, []attr.Value{}),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
