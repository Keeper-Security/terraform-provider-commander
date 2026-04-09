// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseScimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseScimResourceModel

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

		createScimResponse, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpCreateScim)
		if err != nil {
			return err
		}

		var createdScimInfo utils.EnterpriseScimResponse
		if err := utils.UnmarshalApiResponse(createScimResponse.Data, &createdScimInfo); err != nil {
			return err
		}

		data.Id = types.StringValue(strconv.Itoa(createdScimInfo.ScimID))
		data.ScimURL = types.StringValue(createdScimInfo.ScimURL)
		data.UniqueGroups = types.BoolValue(createdScimInfo.UniqueGroups)
		data.ProvisioningToken = types.StringValue(createdScimInfo.ProvisioningToken)
		data.Status = types.StringValue("inactive")

		return nil
	}, ErrSummaryCreateFailed, &resp.Diagnostics); err != nil {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
