// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"
	"errors"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

		state.Id = types.StringValue(strconv.Itoa(scimInfo.ScimID))
		state.ScimURL = types.StringValue(scimInfo.ScimURL)
		// use same state.node value as it is immutable
		state.Node = types.StringValue(state.Node.ValueString())
		state.Status = types.StringValue(scimInfo.Status)
		// Store empty prefix as null so state matches config (prefix = null) and plan shows no change.
		if scimInfo.Prefix == "" {
			state.Prefix = types.StringNull()
		} else {
			state.Prefix = types.StringValue(scimInfo.Prefix)
		}
		state.UniqueGroups = types.BoolValue(scimInfo.UniqueGroups)
		// Use token from API so Read matches Update (API may rotate token on update); avoids "inconsistent result after apply"
		state.ProvisioningToken = types.StringValue(scimInfo.ProvisioningToken)

		return nil
	}, ErrSummaryReadFailed, &resp.Diagnostics); err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
