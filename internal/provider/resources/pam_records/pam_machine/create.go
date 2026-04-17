// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"
	"fmt"
	"strings"

	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data commonpammachine.PamMachineResourceModel

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

	command := buildAddPamMachineRecordCommand(data)
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

	// TODO: Phase 2 – apply PAM settings when pam_settings fields are defined.

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildAddPamMachineRecordCommand(data commonpammachine.PamMachineResourceModel) string {
	parts := []string{utils.CmdRecordAdd}

	parts = append(parts, fmt.Sprintf("%s %s", utils.FlagRecordType, utils.RecordTypePamMachine))

	parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, data.Title.ValueString()))

	appendHostnameOrIPField(&parts, data.HostnameOrIP)

	appendOptionalTextField(&parts, FlagOperatingSystem, data.OperatingSystem)
	appendOptionalTextField(&parts, FlagInstanceName, data.InstanceName)
	appendOptionalTextField(&parts, FlagInstanceId, data.InstanceId)
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
