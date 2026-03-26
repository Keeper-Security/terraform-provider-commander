// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ManagedCompanyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state ManagedCompanyResourceModel

	// Get planned data from Terraform
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
		command := fmt.Sprintf("msp-remove '%s' -f", state.Id.ValueString())
		_, err := r.ApiManager.ExecuteCommand(ctx, command, "Unable to delete managed company")
		if err != nil {
			return err
		}
		resp.State.RemoveResource(ctx)
		return nil
	}, "Delete Managed Company Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
}
