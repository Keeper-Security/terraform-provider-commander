// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"
	"fmt"
	"strings"

	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamRemoteBrowserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data commonpamremotebrowser.PamRemoteBrowserResourceModel

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

	// Phase 1: create the record
	command := buildAddPamRemoteBrowserRecordCommand(data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailAddPamRemoteBrowserRecordFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryAddPamRemoteBrowserRecordFailed, err.Error())
		return
	}

	createdRecordUID, ok := apiResp.Data.(map[string]interface{})["record_uid"].(string)
	if !ok {
		resp.Diagnostics.AddError(ErrSummaryAddPamRemoteBrowserRecordFailed, fmt.Sprintf("Failed to extract record UID from response. API response: %s", apiResp.Data))
		return
	}
	data.Id = types.StringValue(createdRecordUID)

	// Phase 2: apply PAM remote browser settings (`pam rbi edit`).
	if data.PamRemoteBrowserSettings != nil {
		editCmd := BuildPamRbiEditCommand(createdRecordUID, data.PamRemoteBrowserSettings)
		if _, err := r.ApiManager.ExecuteCommand(ctx, editCmd, ErrDetailPamRbiEditFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamRbiEditFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildAddPamRemoteBrowserRecordCommand(data commonpamremotebrowser.PamRemoteBrowserResourceModel) string {
	parts := []string{utils.CmdRecordAdd}

	// Record type
	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamRemoteBrowser))

	// Title
	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))

	// URL
	parts = append(parts, fmt.Sprintf("'%s=%s'", utils.FlagRbiUrl, data.Url.ValueString()))

	// Folder
	if !data.Folder.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagFolder, data.Folder.ValueString()))
	}

	// Notes
	if !data.Notes.IsNull() {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, data.Notes.ValueString()))
	}

	return strings.Join(parts, " ")
}
