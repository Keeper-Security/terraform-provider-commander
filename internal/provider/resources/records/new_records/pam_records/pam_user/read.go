// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PamUserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	id := strings.TrimSpace(state.Id.ValueString())
	if id == "" {
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryReadFailed, "new PAM User record id is empty")
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonpamrecords.FetchVaultRecord(ctx, r.ApiManager, id)
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamUser {
		resp.Diagnostics.AddError(
			commonpamuser.ErrSummaryReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamUser),
		)
		return
	}

	commonpamuser.MapVaultRecordToState(&rec, &state.PamUserSharedModel)

	rotCmd := fmt.Sprintf("%s %s '%s' %s", commonpamuser.CmdPamRotationInfo, commonpamuser.FlagRecordShort, id, utils.FlagFormatJSON)
	rotResp, err2 := r.ApiManager.ExecuteCommand(ctx, rotCmd, commonpamuser.ErrDetailRotationInfoFailed)
	if err2 == nil && rotResp.Data != nil {
		var rotInfo commonpamuser.PamRotationInfoResponse
		if err := utils.UnmarshalApiResponse(rotResp.Data, &rotInfo); err != nil {
			resp.Diagnostics.AddError(commonpamuser.ErrDetailRotationInfoFailed, err.Error())
			return
		}
		commonpamuser.MapRotationSettingsToState(&rotInfo, &rec, state.RotationSettings, &state.PamUserSharedModel)
	}

	if err := new_share.MapResponseToModel(rec.UserPermissions, &state.ShareModel); err != nil {
		resp.Diagnostics.AddError(commonpamuser.ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
