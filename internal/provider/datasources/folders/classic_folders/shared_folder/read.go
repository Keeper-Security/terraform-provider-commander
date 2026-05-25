// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"

	commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/shared_folder"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *ClassicSharedFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SharedFolderDataSourceModel

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

	key := data.SharedFolder.ValueString()
	apiData, err := commonsharedfolder.FetchSharedFolderByNameOrId(ctx, d.ApiManager, key)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	priorUsers := types.MapNull(commonsharedfolder.UserEntryMapElemType)
	priorRecords := types.MapNull(commonsharedfolder.RecordEntryMapElemType)
	if err := commonsharedfolder.MapResponseToModel(apiData, &data.Model, priorUsers, priorRecords); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
