// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseTeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseTeamResourceModel
	var state EnterpriseTeamResourceModel

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

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// NOTE: We should not allow user to update managed company, bec. once team is created in managed company, if we allow user to update managed company then switching to that MC we will not able to find that team, so command will fail.
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			"Managed Company Cannot Be Updated",
			"Cannot update the managed_company field. Once an enterprise team is created in a managed company, the managed company cannot be changed.",
		)
		return
	}

	// Use managed company from plan (or state if plan doesn't have it)
	managedCompany := plan.ManagedCompany
	if managedCompany.IsNull() || managedCompany.IsUnknown() {
		managedCompany = state.ManagedCompany
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, managedCompany, func() error {
		if err := updateEnterpriseTeam(ctx, r.apiManager, &plan, &state); err != nil {
			return err
		}

		plan.Id = state.Id

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise Team Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func updateEnterpriseTeam(ctx context.Context, apiManager *api.ApiManager, plan *EnterpriseTeamResourceModel, state *EnterpriseTeamResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-team")

	// Required parameters
	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	// Optional parameters
	if !state.RestrictEdit.Equal(plan.RestrictEdit) {
		if !plan.RestrictEdit.IsNull() && plan.RestrictEdit.ValueBool() {
			parts = append(parts, "--restrict-edit on")
		} else {
			parts = append(parts, "--restrict-edit off")
		}
	}

	if !state.RestrictShare.Equal(plan.RestrictShare) {
		if !plan.RestrictShare.IsNull() && plan.RestrictShare.ValueBool() {
			parts = append(parts, "--restrict-share on")
		} else {
			parts = append(parts, "--restrict-share off")
		}
	}

	if !state.RestrictView.Equal(plan.RestrictView) {
		if !plan.RestrictView.IsNull() && plan.RestrictView.ValueBool() {
			parts = append(parts, "--restrict-view on")
		} else {
			parts = append(parts, "--restrict-view off")
		}
	}

	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	// Process users and roles changes
	if !state.Users.Equal(plan.Users) {
		users, err := utils.FetchAndProcessUsers(ctx, apiManager, state.Users, plan.Users)
		if err != nil {
			return err
		}
		if users != "" {
			parts = append(parts, users)
		}
	}

	if !state.Roles.Equal(plan.Roles) {
		roles, err := utils.FetchAndProcessRoles(ctx, apiManager, state.Roles, plan.Roles)
		if err != nil {
			return err
		}
		if roles != "" {
			parts = append(parts, roles)
		}
	}

	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise team")
	if err != nil {
		return err
	}

	return nil
}
