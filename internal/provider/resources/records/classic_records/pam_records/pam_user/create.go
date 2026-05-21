// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"fmt"

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

	// Phase 1: create the record.
	command := buildRecordAddPamUserCommand(data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailCreateFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	createdRecordUID, ok := apiResp.Data.(map[string]interface{})["record_uid"].(string)
	if !ok {
		resp.Diagnostics.AddError(
			ErrSummaryCreateFailed,
			fmt.Sprintf("Failed to extract record UID from response. API response: %v", apiResp.Data),
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

	// Phase 2: apply rotation settings (`pam rotation edit`).
	if data.RotationSettings != nil {
		editCmd := buildPamRotationEditCommand(createdRecordUID, data.RotationSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, ErrDetailRotationEditFailed); err != nil {
			delCmd := fmt.Sprintf("%s '%s' %s", utils.CmdRecordDelete, createdRecordUID, utils.FlagForce)
			if _, delErr := r.ApiManager.ExecuteCommand(ctx, delCmd, utils.ErrDetailRecordDeleteFailed); delErr != nil {
				resp.Diagnostics.AddWarning(
					"PAM User rotation failed and rollback delete failed",
					fmt.Sprintf("Rotation configuration failed: %v. Removing the newly created record also failed: %v. The pamUser may still exist (uid %s).", err, delErr, createdRecordUID),
				)
			}
			resp.Diagnostics.AddError(ErrSummaryRotationEditFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
