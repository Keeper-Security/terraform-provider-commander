// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Update allows only prefix and unique_groups to be changed.
func (r *EnterpriseScimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state EnterpriseScimResourceModel

	// Get planned data (new values)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	utils.RestrictAttributeUpdate(&resp.Diagnostics, []utils.ImmutableAttributeCheck{
		{
			PlanValue:  plan.ManagedCompany,
			StateValue: state.ManagedCompany,
			Summary:    ErrSummaryManagedCompanyCannotBeUpdated,
			Detail:     ErrDetailManagedCompany,
		},
		{
			PlanValue:  plan.Node,
			StateValue: state.Node,
			Summary:    ErrSummaryNodeCannotBeUpdated,
			Detail:     ErrDetailNode,
		},
	})
	if resp.Diagnostics.HasError() {
		return
	}

	// Use state's managed_company for API context (plan equals state here; state is source of truth for identity).
	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, plan.ManagedCompany, func() error {
		command := buildUpdateCommand(state.Id.ValueString(), &plan)
		updatedScimResponse, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpUpdateScim)
		if err != nil {
			return err
		}

		var updatedScimInfo utils.EnterpriseScimResponse
		if err := utils.UnmarshalApiResponse(updatedScimResponse.Data, &updatedScimInfo); err != nil {
			return err
		}

		// Update returns a new provisioning token in res.message; store it so state has the current token.
		plan.ProvisioningToken = types.StringValue(updatedScimInfo.ProvisioningToken)
		return nil
	}, ErrSummaryUpdateFailed, &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve id, node, scim_url from state (provisioning_token already set from API response above).
	plan.Id = state.Id
	plan.Node = state.Node
	plan.ScimURL = state.ScimURL

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
