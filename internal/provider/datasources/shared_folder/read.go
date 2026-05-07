// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"

	commonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/shared_folder"
	sfres "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/shared_folder"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *SharedFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
		resp.Diagnostics.AddError(sfres.ErrSummarySyncDownFailed, err.Error())
		return
	}

	key := data.SharedFolder.ValueString()
	apiData, err := commonsharedfolder.FetchSharedFolderByNameOrId(ctx, d.ApiManager, key)
	if err != nil {
		resp.Diagnostics.AddError(sfres.ErrSummaryReadFailed, err.Error())
		return
	}

	priorUsers := types.MapNull(commonsharedfolder.UserEntryMapElemType)
	priorRecords := types.MapNull(commonsharedfolder.RecordEntryMapElemType)
	if err := commonsharedfolder.MapResponseToModel(apiData, &data.Model, priorUsers, priorRecords); err != nil {
		resp.Diagnostics.AddError(sfres.ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
