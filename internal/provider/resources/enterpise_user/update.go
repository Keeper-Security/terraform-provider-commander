package enterpiseuser

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan EnterpriseUserResourceModel
	var state EnterpriseUserResourceModel

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

	// NOTE: We should not allow user to update managed company, bec. once role is created in managed company, if we allow user to update managed company then switching to that MC we will not able to find that role, so command will fail.
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			"Managed Company Cannot Be Updated",
			"Cannot update the managed_company field. Once an enterprise role is created in a managed company, the managed company cannot be changed.",
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
		// Teams cannot be updated when user status is Invited
		if !state.Status.IsNull() && !state.Status.IsUnknown() &&
			state.Status.ValueString() == UserInvitedStatus && !state.Teams.Equal(plan.Teams) {
			return fmt.Errorf("User with 'Invited' status cannot be added to teams.")
		}

		if !state.Email.Equal(plan.Email) {
			return fmt.Errorf("Email can not be changed.")
		}

		if err := updateUserAttributes(ctx, r.apiManager, &plan, &state); err != nil {
			return err
		}

		// Keep the same ID
		plan.Id = state.Id

		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise User Failed",
			err.Error(),
		)
		return
	}

	plan.Id = state.Id

	plan.Status = state.Status

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func updateUserAttributes(ctx context.Context, apiManager *api.ApiManager, plan *EnterpriseUserResourceModel, state *EnterpriseUserResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-user")

	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	if !state.JobTitle.Equal(plan.JobTitle) {
		parts = append(parts, fmt.Sprintf("--job-title '%s'", plan.JobTitle.ValueString()))
	}

	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	if !state.Roles.Equal(plan.Roles) {
		roles, err := utils.FetchAndProcessRoles(ctx, apiManager, state.Roles, plan.Roles, "--add-role", "--remove-role")
		if err != nil {
			return err
		}
		if roles != "" {
			parts = append(parts, fmt.Sprintf("%s", roles))
		}
	}

	if !state.Teams.Equal(plan.Teams) {
		teams, err := utils.FetchAndProcessTeams(ctx, apiManager, state.Teams, plan.Teams, "--add-team", "--remove-team")
		if err != nil {
			return err
		}
		if teams != "" {
			parts = append(parts, fmt.Sprintf("%s", teams))
		}
	}

	parts = append(parts, fmt.Sprintf("'%s' -f", state.Id.ValueString()))

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise user")
	if err != nil {
		return err
	}

	return nil
}
