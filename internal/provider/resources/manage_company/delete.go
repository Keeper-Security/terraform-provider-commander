package managecompany

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ManageCompanyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state ManageCompanyResourceModel

	// Get planned data from Terraform
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

	// Build delete command
	command := fmt.Sprintf("msp-remove '%s' -f", state.Id.ValueString())

	// Only switch to MSP if it's an MSP account
	if r.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, r.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Delete Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to delete managed company")
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Managed Company Failed",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
