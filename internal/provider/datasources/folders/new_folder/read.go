// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"

	commonnewfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/new_folder"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Read loads the Keeper Drive folder identified by `new_folder` (UID or name)
func (d *NewFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NewFolderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	key := data.NewFolder.ValueString()
	if key == "" {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, "new_folder lookup value is empty; provide the folder UID or name.")
		return
	}

	apiData, err := commonnewfolder.FetchNewFolderByNameOrId(ctx, d.ApiManager, key)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	if err := commonnewfolder.MapResponseToModel(apiData, &data.Model); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	if err := new_share.MapResponseToModel(apiData.UserPermissions, &data.Model.ShareModel); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
