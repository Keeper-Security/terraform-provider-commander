// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"fmt"
	"strings"

	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_directory"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamDirectoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data commonpamdirectory.PamDirectoryResourceModel

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

	command := buildAddPamDirectoryRecordCommand(data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailAddPamDirectoryRecordFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamDirectoryRecordFailed, err.Error())
		return
	}

	createdRecordUID, ok := apiResp.Data.(map[string]interface{})["record_uid"].(string)
	if !ok {
		resp.Diagnostics.AddError(ErrSummaryAddPamDirectoryRecordFailed, fmt.Sprintf("Failed to extract record UID from response. API response: %s", apiResp.Data))
		return
	}
	data.Id = types.StringValue(createdRecordUID)

	if data.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, createdRecordUID, data.PamSettings); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamTunnelSettingsFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildAddPamDirectoryRecordCommand(data commonpamdirectory.PamDirectoryResourceModel) string {
	parts := []string{utils.CmdRecordAdd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamDirectory))
	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))

	appendHostnameOrIPField(&parts, data.HostnameOrIP)
	appendOptionalCheckboxField(&parts, FlagUseSSL, data.UseSSL)
	appendOptionalTextField(&parts, FlagDomainName, data.DomainName)
	appendAlternativeIPsField(&parts, data.AlternativeIPs)
	appendOptionalTextField(&parts, FlagDirectoryId, data.DirectoryId)
	appendOptionalDirectoryTypeField(&parts, data.DirectoryType)
	appendOptionalTextField(&parts, FlagUserMatch, data.UserMatch)
	appendOptionalTextField(&parts, FlagProviderGroup, data.ProviderGroup)
	appendOptionalTextField(&parts, FlagProviderRegion, data.ProviderRegion)

	if !data.Folder.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagFolder, data.Folder.ValueString()))
	}

	if !data.Notes.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}
