// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PamMachineResourceModel

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

	command := commonpammachine.BuildAddCommand(utils.CmdRecordAdd, data.PamMachineResourceModel)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailAddPamMachineRecordFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamMachineRecordFailed, err.Error())
		return
	}

	createdRecordUID, ok := apiResp.Data.(map[string]interface{})["record_uid"].(string)
	if !ok {
		resp.Diagnostics.AddError(ErrSummaryAddPamMachineRecordFailed, fmt.Sprintf("Failed to extract record UID from response. API response: %s", apiResp.Data))
		return
	}
	data.Id = types.StringValue(createdRecordUID)

	if data.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, createdRecordUID, data.PamSettings, nil); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamTunnelSettingsFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, createdRecordUID, data.Share, types.MapNull(classic_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamMachineRecordFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
