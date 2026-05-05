// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"
	"fmt"
	"strings"

	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_database"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data commonpamdatabase.PamDatabaseResourceModel

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

	command := buildAddPamDatabaseRecordCommand(data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailAddPamDatabaseRecordFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamDatabaseRecordFailed, err.Error())
		return
	}

	createdRecordUID, ok := apiResp.Data.(map[string]interface{})["record_uid"].(string)
	if !ok {
		resp.Diagnostics.AddError(ErrSummaryAddPamDatabaseRecordFailed, fmt.Sprintf("Failed to extract record UID from response. API response: %s", apiResp.Data))
		return
	}
	data.Id = types.StringValue(createdRecordUID)

	if data.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, createdRecordUID, data.PamSettings, nil); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamSettingsFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildAddPamDatabaseRecordCommand(data commonpamdatabase.PamDatabaseResourceModel) string {
	parts := []string{utils.CmdRecordAdd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamDatabase))
	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))

	commonpamrecords.AppendHostnameOrIPField(&parts, data.HostnameOrIP)
	commonpamrecords.AppendOptionalCheckboxField(&parts, FlagUseSSL, data.UseSSL)
	commonpamrecords.AppendOptionalTextField(&parts, FlagDatabaseId, data.DatabaseId)
	appendOptionalDatabaseTypeField(&parts, data.DatabaseType)
	commonpamrecords.AppendOptionalTextField(&parts, FlagProviderGroup, data.ProviderGroup)
	commonpamrecords.AppendOptionalTextField(&parts, FlagProviderRegion, data.ProviderRegion)

	if !data.Folder.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagFolder, data.Folder.ValueString()))
	}

	if !data.Notes.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}
