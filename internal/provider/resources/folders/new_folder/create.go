// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"
	"fmt"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *NewFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NewFolderResourceModel

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

	// Create new folder (Nested Share Folder)
	command := fmt.Sprintf(`%s "%s"`, CmdNsfMkdir, data.Name.ValueString())

	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpCreate)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	data.Id = types.StringValue(string(apiResp.Message))

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdShareFolder, data.Id.ValueString(), data.Share, types.MapNull(new_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
