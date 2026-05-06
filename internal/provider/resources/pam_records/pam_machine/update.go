// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"
	"fmt"
	"strings"

	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_machine"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan commonpammachine.PamMachineResourceModel
	var state commonpammachine.PamMachineResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	resp.Diagnostics.Append(commonpamrecords.ValidatePamSettingsFieldsNotRemoved(plan.PamSettings, state.PamSettings)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id
	recordUID := strings.TrimSpace(plan.Id.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryPamMachineRecordUpdateFailed, "PAM machine record id is empty")
		return
	}

	// Phase 0: Move record to destination folder if folder is changed.
	if err := commonpamrecords.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.Folder.ValueString(), state.Folder.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	// Phase 1: record fields.
	if recordUpdateHasMutations(plan, state) {
		cmd := buildUpdatePamMachineRecordCommand(recordUID, plan, state)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailPamMachineRecordUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamMachineRecordUpdateFailed, err.Error())
			return
		}
	}

	// Phase 2 – apply PAM settings when pam_settings fields are defined.
	if plan.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, recordUID, plan.PamSettings, state.PamSettings); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamSettingsFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func buildUpdatePamMachineRecordCommand(recordUID string, plan, state commonpammachine.PamMachineResourceModel) string {
	parts := []string{
		utils.CmdRecordUpdate,
		fmt.Sprintf("%s '%s'", utils.FlagRecord, recordUID),
	}

	if !plan.Title.Equal(state.Title) {
		parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagTitle, plan.Title.ValueString()))
	}

	if !commonpamrecords.HostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) {
		commonpamrecords.AppendHostnameOrIPField(&parts, plan.HostnameOrIP)
	}

	commonpamrecords.AppendChangedTextField(&parts, FlagOperatingSystem, plan.OperatingSystem, state.OperatingSystem)
	commonpamrecords.AppendChangedTextField(&parts, FlagInstanceName, plan.InstanceName, state.InstanceName)
	commonpamrecords.AppendChangedTextField(&parts, FlagInstanceId, plan.InstanceId, state.InstanceId)
	commonpamrecords.AppendChangedTextField(&parts, FlagProviderGroup, plan.ProviderGroup, state.ProviderGroup)
	commonpamrecords.AppendChangedTextField(&parts, FlagProviderRegion, plan.ProviderRegion, state.ProviderRegion)

	if !plan.Notes.Equal(state.Notes) && !plan.Notes.IsUnknown() {
		if plan.Notes.IsNull() {
			parts = append(parts, fmt.Sprintf("%s ''", utils.FlagNotes))
		} else {
			parts = append(parts, fmt.Sprintf("%s '%s'", utils.FlagNotes, plan.Notes.ValueString()))
		}
	}

	return strings.Join(parts, " ")
}
