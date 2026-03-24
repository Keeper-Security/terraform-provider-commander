// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EpmPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EpmPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
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

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
		command := buildCreateCommand(&data)

		_, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpCreateEpmPolicy)
		if err != nil {
			return err
		}

		data.Id = types.StringValue("123")
		return nil
	}, ErrSummaryCreateFailed, &resp.Diagnostics); err != nil {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
