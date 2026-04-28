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
)

func (r *PamDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan commonpamdatabase.PamDatabaseResourceModel
	var state commonpamdatabase.PamDatabaseResourceModel

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

	plan.Id = state.Id
	recordUID := strings.TrimSpace(plan.Id.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryPamDatabaseRecordUpdateFailed, "PAM database record id is empty")
		return
	}

	if err := commonpamrecords.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.Folder.ValueString(), state.Folder.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	if recordUpdateHasMutations(plan, state) {
		cmd := buildUpdatePamDatabaseRecordCommand(recordUID, plan, state)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailPamDatabaseRecordUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryPamDatabaseRecordUpdateFailed, err.Error())
			return
		}
	}

	if plan.PamSettings != nil {
		if err := commonpamrecords.ApplyPamSettings(ctx, r.ApiManager, recordUID, plan.PamSettings); err != nil {
			resp.Diagnostics.AddError(utils.ErrSummaryApplyPamSettingsFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func buildUpdatePamDatabaseRecordCommand(recordUID string, plan, state commonpamdatabase.PamDatabaseResourceModel) string {
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

	commonpamrecords.AppendChangedCheckboxField(&parts, FlagUseSSL, plan.UseSSL, state.UseSSL)
	commonpamrecords.AppendChangedTextField(&parts, FlagDatabaseId, plan.DatabaseId, state.DatabaseId)

	if !plan.DatabaseType.Equal(state.DatabaseType) {
		appendOptionalDatabaseTypeField(&parts, plan.DatabaseType)
	}

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
