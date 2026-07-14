// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PamUserResourceModel

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

	command := commonpamuser.BuildAddCommand(utils.CmdNsfRecordAdd, data.PamUserSharedModel)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, commonpamuser.ErrDetailCreateFailed)
	if err != nil {
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryCreateFailed, err.Error())
		return
	}

	createdRecordUID := string(apiResp.Message)
	if createdRecordUID == "" {
		resp.Diagnostics.AddError(
			commonpamuser.ErrSummaryCreateFailed,
			fmt.Sprintf("Failed to extract record UID from response. API response: %v", apiResp),
		)
		return
	}
	data.Id = types.StringValue(createdRecordUID)
	if data.Login.IsUnknown() {
		data.Login = types.StringNull()
	}
	if data.Password.IsUnknown() {
		data.Password = types.StringNull()
	}

	if data.RotationSettings != nil {
		editCmd := commonpamuser.BuildPamRotationEditCommand(createdRecordUID, data.RotationSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, commonpamuser.ErrDetailRotationEditFailed); err != nil {
			delCmd := fmt.Sprintf("%s '%s' %s", utils.CmdRecordDelete, createdRecordUID, utils.FlagForce)
			if _, delErr := r.ApiManager.ExecuteCommand(ctx, delCmd, utils.ErrDetailRecordDeleteFailed); delErr != nil {
				resp.Diagnostics.AddWarning(
					"New PAM User rotation failed and rollback delete failed",
					fmt.Sprintf("Rotation configuration failed: %v. Removing the newly created record also failed: %v. The pamUser may still exist (uid %s).", err, delErr, createdRecordUID),
				)
			}
			resp.Diagnostics.AddError(commonpamuser.ErrSummaryRotationEditFailed, err.Error())
			return
		}
	}

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdNsfShareRecord, createdRecordUID, data.Share, types.MapNull(new_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamUserRecordFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
