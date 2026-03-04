// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"
	"fmt"

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

		_, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpCreateScim)
		if err != nil {
			return err
		}

		// TODO: Remove this once we get the data after create
		data.Id = types.StringValue("123")
		data.ScimURL = types.StringValue(fmt.Sprintf("https://keepersecurity.com/api/rest/scim/v2/%s", data.Node.ValueString()))
		data.Status = types.StringValue("inactive")
		return nil
	}, ErrSummaryCreateFailed, &resp.Diagnostics); err != nil {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// // Currently we are not getting data after create
	// if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
	// 	scimInfo, err := utils.FetchEnterpriseScimById(ctx, r.ApiManager, data.Id.ValueString())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if scimInfo != nil {
	// 		mapScimReadResponseToModel(scimInfo, &data)
	// 	}
	// 	return nil
	// }, "Read Enterprise SCIM After Create Failed", &resp.Diagnostics); err != nil {
	// 	return
	// }
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
