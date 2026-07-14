// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *PamUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamUserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	recordUID := strings.TrimSpace(data.PamUser.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "pam_user is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	// Phase 1: fetch the vault record
	apiResp, err := commonrecordsutils.FetchVaultRecord(ctx, d.ApiManager, recordUID)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("record %q not found or empty response", recordUID))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamUser {
		resp.Diagnostics.AddError(
			ErrSummaryReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamUser),
		)
		return
	}

	commonpamuser.MapVaultRecordToState(&rec, &data.PamUserSharedModel)

	// Phase 2: fetch rotation info
	rotCmd := fmt.Sprintf("%s %s '%s' %s", commonpamuser.CmdPamRotationInfo, commonpamuser.FlagRecordShort, recordUID, utils.FlagFormatJSON)
	if rotResp, err := d.ApiManager.ExecuteCommand(ctx, rotCmd, commonpamuser.ErrDetailRotationInfoFailed); err == nil && rotResp != nil {
		var rotInfo commonpamuser.PamRotationInfoResponse
		if err := utils.UnmarshalApiResponse(rotResp.Data, &rotInfo); err != nil {
			resp.Diagnostics.AddError(commonpamuser.ErrDetailRotationInfoFailed, err.Error())
			return
		}
		commonpamuser.MapRotationSettingsToState(&rotInfo, &rec, data.RotationSettings, &data.PamUserSharedModel)
	}

	if err := classic_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
