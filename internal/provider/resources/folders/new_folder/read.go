// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"
	"errors"

	commonnewfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/new_folder"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Read refreshes the Terraform state for a Keeper Drive folder via nsf-get.
//
// Behavior:
//   - empty id in state -> diagnostic error (this should never happen for a
//     resource that completed Create; usually means hand-edited state).
//   - nsf-get reports the folder no longer exists -> RemoveResource so the next
//     plan re-creates it.
//   - nsf-get succeeds -> refresh Name (and Id, defensively) from the API.
//   - share block is reconciled from the API's user_permissions array only if
//     the user is actually managing shares (state.Share is non-null). This is
//     the Optional-only semantics: a user who never declared `share` in config
//     should not see drift appear out of nowhere on refresh. Owner entries
//     returned by the API are dropped by new_share.MapResponseToModel.
func (r *NewFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NewFolderResourceModel

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
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, "new folder id is empty")
		return
	}

	apiData, err := commonnewfolder.FetchNewFolderByNameOrId(ctx, r.ApiManager, id)
	if err != nil {
		if errors.Is(err, commonnewfolder.ErrNestedSharedFolderNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	if err := commonnewfolder.MapResponseToModel(apiData, &state); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	if !state.Share.IsNull() && !state.Share.IsUnknown() {
		if err := new_share.MapResponseToModel(apiData.UserPermissions, &state.ShareModel); err != nil {
			resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
