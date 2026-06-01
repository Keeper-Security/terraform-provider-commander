// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamremotebrowser

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamRemoteBrowserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PamRemoteBrowserResourceModel

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

	command := commonpamremotebrowser.BuildAddCommand(utils.CmdNsfRecordAdd, data.PamRemoteBrowserResourceModel)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailAddPamRemoteBrowserRecordFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamRemoteBrowserRecordFailed, err.Error())
		return
	}

	createdRecordUID := string(apiResp.Message)
	if createdRecordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryAddPamRemoteBrowserRecordFailed, fmt.Sprintf("Failed to extract record UID from response. API response: %v", apiResp))
		return
	}
	data.Id = types.StringValue(createdRecordUID)

	if data.PamRemoteBrowserSettings != nil {
		editCmd := commonpamremotebrowser.BuildPamRbiEditCommand(createdRecordUID, data.PamRemoteBrowserSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, ErrDetailPamRbiEditFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamRbiEditFailed, err.Error())
			return
		}
	}

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdShareRecord, createdRecordUID, data.Share, types.MapNull(new_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamRemoteBrowserRecordFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
