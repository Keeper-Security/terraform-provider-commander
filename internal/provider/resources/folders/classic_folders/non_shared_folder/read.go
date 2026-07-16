// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"
	"errors"

	commonnonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/non_shared_folder"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *NonSharedFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NonSharedFolderResourceModel

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

	id := state.Id.ValueString()
	if id == "" {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, "folder id is empty")
		return
	}

	apiData, err := commonnonsharedfolder.FetchFolderByNameOrId(ctx, r.ApiManager, id)
	if err != nil {
		if errors.Is(err, commonnonsharedfolder.ErrNonSharedFolderNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	if err := commonnonsharedfolder.MapResponseToModel(ctx, apiData, &state); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
