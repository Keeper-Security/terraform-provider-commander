// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseScimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnterpriseScimResourceModel

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

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, state.ManagedCompany, func() error {
		command := fmt.Sprintf("%s %s %s", CmdScimDelete, state.Id.ValueString(), FlagForce)
		_, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpDeleteScim)

		if err != nil {
			return err
		}
		return nil
	}, ErrSummaryDeleteFailed, &resp.Diagnostics); err != nil {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
}
