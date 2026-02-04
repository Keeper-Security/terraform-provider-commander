package enterprisenode

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnterpriseNodeResourceModel

	// Get state from Terraform
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

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {
		// Build delete command
		command := fmt.Sprintf("enterprise-node --delete '%s'", state.Id.ValueString())

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to delete enterprise node")
		if err != nil {
			return fmt.Errorf("Delete Enterprise Node Failed: %w", err)
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Enterprise Node Failed",
			err.Error(),
		)
		return
	}
}
