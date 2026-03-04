// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"
	"errors"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *EnterpriseScimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseScimResourceModel

	// Get current state (old values)
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
		scimInfo, err := utils.FetchEnterpriseScimById(ctx, r.ApiManager, state.Id.ValueString())
		if err != nil {
			return err
		}
		if scimInfo == nil {
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}
		mapScimReadResponseToModel(scimInfo, &state)
		return nil
	}, ErrSummaryReadFailed, &resp.Diagnostics); err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
