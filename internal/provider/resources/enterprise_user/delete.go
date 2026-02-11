// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnterpriseUserResourceModel

	// Get state from Terraform
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

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, state.ManagedCompany, func() error {
		command := fmt.Sprintf("enterprise-user --delete '%s' -f", state.Id.ValueString())
		_, err := r.ApiManager.ExecuteCommand(ctx, command, "Unable to delete enterprise user")
		return err
	}, "Delete Enterprise User Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
}
