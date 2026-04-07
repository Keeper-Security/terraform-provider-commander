// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"errors"

	commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/shared_folder"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *SharedFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SharedFolderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(ErrSummarySyncDownFailed, err.Error())
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "shared folder id is empty")
		return
	}

	// State always stores the canonical UID; Commander get accepts UID or path, same as the data source.
	apiData, err := commonsharedfolder.FetchSharedFolderByNameOrId(ctx, r.ApiManager, id)
	if err != nil {
		if errors.Is(err, commonsharedfolder.ErrSharedFolderNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if err := commonsharedfolder.MapResponseToModel(apiData, &state, state.Users, state.Records); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
