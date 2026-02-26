// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ShareFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ShareFolderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			ErrSummaryProviderConfig,
			err.Error(),
		)
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(ErrSummarySyncDownFailed, err.Error())
		return
	}

	command := fmt.Sprintf("%s '%s' %s %s", CmdRmdir, state.Id.ValueString(), FlagForce, FlagQuiet)
	if _, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpDeleteSF); err != nil {
		resp.Diagnostics.AddError(ErrSummaryDeleteFailed, err.Error())
		return
	}
}
