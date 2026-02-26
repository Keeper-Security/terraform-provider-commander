// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ShareFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ShareFolderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(ErrSummaryProviderConfig, err.Error())
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

	command := fmt.Sprintf("%s '%s' %s %s", CmdGet, id, FlagFormat, FormatJSON)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpGetSF)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	model, err := MapGetResponseToModel(ctx, apiResp.Data)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
