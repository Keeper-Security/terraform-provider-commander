// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DatabaseResourceModel
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
		resp.Diagnostics.AddError(utils.ErrSummaryRecordDeleteFailed, "Database record id is empty")
		return
	}
	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}
	cmd := fmt.Sprintf("%s '%s' %s", utils.CmdNsfRecordDelete, id, utils.FlagForce)
	if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, utils.ErrDetailRecordDeleteFailed); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryRecordDeleteFailed, err.Error())
		return
	}
}
