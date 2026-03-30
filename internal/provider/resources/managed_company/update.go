// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ManagedCompanyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ManagedCompanyResourceModel
	var state ManagedCompanyResourceModel

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

	if err := utils.RunWithMspContext(ctx, r.ApiManager, func() error {
		return updateManagedCompany(ctx, r.ApiManager, &plan, &state)
	}, "Update Managed Company Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Keep the same ID
	plan.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// buildManagedCompanyUpdateCommand builds the Commander CLI command for updating a managed company.
func updateManagedCompany(ctx context.Context, apiManager *api.ApiManager, plan *ManagedCompanyResourceModel, state *ManagedCompanyResourceModel) error {
	var parts []string

	parts = append(parts, "msp-update")
	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	// Required fields - no null check needed (validation ensures they exist)
	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	if !state.Plan.Equal(plan.Plan) {
		parts = append(parts, fmt.Sprintf("--plan '%s'", plan.Plan.ValueString()))
	}

	// Optional fields
	if !state.Seats.Equal(plan.Seats) {
		// check null (user might have removed them)
		if plan.Seats.IsNull() {
			parts = append(parts, fmt.Sprintf("--seats %d", 0))

		} else {
			parts = append(parts, fmt.Sprintf("--seats %d", plan.Seats.ValueInt64()))
		}
	}

	if !state.FilePlan.Equal(plan.FilePlan) && !plan.FilePlan.IsNull() {
		if plan.FilePlan.IsNull() {
			parts = append(parts, fmt.Sprintf("--file-plan '%s'", FilePlan100GB))
		} else {
			parts = append(parts, fmt.Sprintf("--file-plan '%s'", plan.FilePlan.ValueString()))
		}
	}

	// Add-ons

	// if there is add-on in state but not in plan, means we need to remove it by adding --remove-addon flag and appending add-ons
	// if there is add-on in plan but not in state, means we need to add it
	// if there is add-on with same base name but different number (e.g., connection_manager:2 -> connection_manager:1), we need to update it

	// Normalize add-ons from state and plan
	stateAddOns := normalizeAddOns(state.AddOns)
	planAddOns := normalizeAddOns(plan.AddOns)

	// Convert to maps for easier comparison
	stateAddOnsMap := make(map[string]bool)
	for _, addOn := range stateAddOns {
		stateAddOnsMap[addOn] = true
	}

	planAddOnsMap := make(map[string]bool)
	for _, addOn := range planAddOns {
		planAddOnsMap[addOn] = true
	}

	// Helper function to extract base add-on name (without number suffix)
	getBaseAddOnName := func(addOn string) string {
		if matches := addOnWithNumberRegex.FindStringSubmatch(addOn); matches != nil {
			return matches[1]
		}
		return addOn
	}

	// Find add-ons to remove (in state but not in plan)
	var addOnsToRemove []string
	// Track base names that are being updated (same base, different number)
	updatedBaseNames := make(map[string]bool)

	for _, stateAddOn := range stateAddOns {
		if !planAddOnsMap[stateAddOn] {
			// Check if this add-on has a number suffix and if there's a matching base name in plan
			baseName := getBaseAddOnName(stateAddOn)
			hasUpdate := false

			// Check if there's an add-on in plan with the same base name but different number
			for _, planAddOn := range planAddOns {
				if getBaseAddOnName(planAddOn) == baseName && planAddOn != stateAddOn {
					// This is an update case (same base, different number)
					updatedBaseNames[baseName] = true
					hasUpdate = true
					break
				}
			}

			// Only add to remove list if it's not an update case
			if !hasUpdate {
				addOnsToRemove = append(addOnsToRemove, stateAddOn)
			}
		}
	}

	// Find add-ons to add (in plan but not in state, or updates with different numbers)
	var addOnsToAdd []string
	for _, planAddOn := range planAddOns {
		if !stateAddOnsMap[planAddOn] {
			// Check if this is an update case (same base name exists in state with different number)
			baseName := getBaseAddOnName(planAddOn)
			if updatedBaseNames[baseName] {
				// This is an update - include it in add list
				addOnsToAdd = append(addOnsToAdd, planAddOn)
			} else {
				// This is a new add-on
				addOnsToAdd = append(addOnsToAdd, planAddOn)
			}
		}
	}

	// Add remove-addon flag for each add-on to remove
	for _, addOn := range addOnsToRemove {
		parts = append(parts, fmt.Sprintf("--remove-addon '%s'", addOn))
	}

	// Add add-addon flag for each add-on to add
	for _, addOn := range addOnsToAdd {
		parts = append(parts, fmt.Sprintf("--add-addon '%s'", addOn))
	}

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update managed company")
	if err != nil {
		return err
	}

	return nil
}
