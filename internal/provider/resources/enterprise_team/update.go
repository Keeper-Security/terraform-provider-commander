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
	if err := r.EnsureApiManager(); err != nil {
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

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, managedCompany, func() error {
		if err := updateEnterpriseTeam(ctx, r.ApiManager, &plan, &state); err != nil {
			return err
		}
		plan.Id = state.Id
		return nil
	}, "Update Enterprise Team Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
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
	// Restrict flags: send off when user had it on and then removed or turned it off
	stateRestrictEditOn := !state.RestrictEdit.IsNull() && state.RestrictEdit.ValueBool()
	planRestrictEditOn := !plan.RestrictEdit.IsNull() && plan.RestrictEdit.ValueBool()
	if planRestrictEditOn && !stateRestrictEditOn {
		parts = append(parts, "--restrict-edit on")
	} else if stateRestrictEditOn && !planRestrictEditOn {
		parts = append(parts, "--restrict-edit off")
	}

	stateRestrictShareOn := !state.RestrictShare.IsNull() && state.RestrictShare.ValueBool()
	planRestrictShareOn := !plan.RestrictShare.IsNull() && plan.RestrictShare.ValueBool()
	if planRestrictShareOn && !stateRestrictShareOn {
		parts = append(parts, "--restrict-share on")
	} else if stateRestrictShareOn && !planRestrictShareOn {
		parts = append(parts, "--restrict-share off")
	}

	// Send --restrict-view when value changed; explicitly send off when user had it on and then removed or turned it off
	stateRestrictViewOn := !state.RestrictView.IsNull() && state.RestrictView.ValueBool()
	planRestrictViewOn := !plan.RestrictView.IsNull() && plan.RestrictView.ValueBool()
	if planRestrictViewOn && !stateRestrictViewOn {
		parts = append(parts, "--restrict-view on")
	} else if stateRestrictViewOn && !planRestrictViewOn {
		parts = append(parts, "--restrict-view off")
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
