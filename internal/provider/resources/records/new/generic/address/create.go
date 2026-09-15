// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordaddress "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/address"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *AddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AddressResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}
	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}
	cmd := commonrecordaddress.BuildAddCommand(utils.CmdNsfRecordAdd, data.AddressModel)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailCreateFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	// nsf-record-add returns the new record UID in Message (same as other NSF resources).
	createdUID := string(apiResp.Message)
	if createdUID == "" {
		resp.Diagnostics.AddError(
			ErrSummaryCreateFailed,
			fmt.Sprintf("Failed to extract record UID from response. API response: %v", apiResp),
		)
		return
	}
	data.Id = types.StringValue(createdUID)

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdNsfShareRecord, createdUID, data.Share, types.MapNull(new_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
