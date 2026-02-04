package enterprisenode

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseNodeResourceModel
	var state EnterpriseNodeResourceModel

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

	// NOTE: We should not allow user to update managed company, bec. once node is created in managed company, if we allow user to update managed company then switching to that MC we will not able to find that node, so command will fail.
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			"Managed Company Cannot Be Updated",
			"Cannot update the managed_company field. Once an enterprise node is created in a managed company, the managed company cannot be changed.",
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

		if err := updateEnterpriseNode(ctx, r.apiManager, &plan, &state); err != nil {
			return err
		}

		// Keep the same ID
		plan.Id = state.Id

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise Node Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func updateEnterpriseNode(ctx context.Context, apiManager *api.ApiManager, plan *EnterpriseNodeResourceModel, state *EnterpriseNodeResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-node")

	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	if !state.Parent.Equal(plan.Parent) {
		value := plan.Parent.ValueString()

		// We need to make the parent as "root" as we if the parent is the same as the managed company
		if !plan.ManagedCompany.IsNull() && plan.Parent.ValueString() == plan.ManagedCompany.ValueString() {
			value = "root"
		}
		parts = append(parts, fmt.Sprintf("--parent '%s'", value))
	}

	// if !state.WipeOut.Equal(plan.WipeOut) {
	// 	parts = append(parts, "--wipe-out")
	// }

	// if !state.ToggleIsolated.Equal(plan.ToggleIsolated) {
	// 	parts = append(parts, "--toggle-isolated")
	// }

	// if !state.LogoFile.Equal(plan.LogoFile) {
	// 	parts = append(parts, fmt.Sprintf("--logo-file '%s'", plan.LogoFile.ValueString()))
	// }

	// Node to update
	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise node")
	if err != nil {
		return err
	}

	return nil
}
