// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamDirectoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PamDirectoryResourceModel

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

	command := commonpamdirectory.BuildAddCommand(utils.CmdNsfRecordAdd, data.PamDirectoryResourceModel)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailAddPamDirectoryRecordFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamDirectoryRecordFailed, err.Error())
		return
	}

	createdRecordUID := string(apiResp.Message)
	if createdRecordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryAddPamDirectoryRecordFailed, fmt.Sprintf("Failed to extract record UID from response. API response: %v", apiResp))
		return
	}
	data.Id = types.StringValue(createdRecordUID)

	if data.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, createdRecordUID, data.PamSettings, nil); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamTunnelSettingsFailed, err.Error())
			return
		}
	}

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdShareRecord, createdRecordUID, data.Share, types.MapNull(new_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamDirectoryRecordFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
