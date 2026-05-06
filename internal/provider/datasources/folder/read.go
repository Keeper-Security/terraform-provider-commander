// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"

	folderres "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folder"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *FolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FolderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(folderres.ErrSummarySyncDownFailed, err.Error())
		return
	}

	key := data.Folder.ValueString()
	command := fmt.Sprintf("%s '%s' %s %s", folderres.CmdGet, key, folderres.FlagFormat, folderres.FormatJSON)
	apiResp, err := d.ApiManager.ExecuteCommand(ctx, command, folderres.ErrOpGetFolder)
	if err != nil {
		resp.Diagnostics.AddError(folderres.ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(folderres.ErrSummaryReadFailed, fmt.Sprintf("folder %q not found", key))
		return
	}

	var folderData folderres.FolderResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &folderData); err != nil {
		resp.Diagnostics.AddError(folderres.ErrSummaryReadFailed, err.Error())
		return
	}

	if folderData.FolderUID == "" {
		resp.Diagnostics.AddError(folderres.ErrSummaryReadFailed, fmt.Sprintf("folder %q not found", key))
		return
	}

	data.Id = types.StringValue(folderData.FolderUID)
	data.Name = types.StringValue(folderData.Name)
	data.Type = types.StringValue(folderData.Type)

	if folderData.Path != "" {
		parent, _ := folderres.SplitFolderPath(folderData.Path)
		if parent != "" {
			data.FolderLocation = types.StringValue(parent)
		} else {
			data.FolderLocation = types.StringValue("")
		}
	} else {
		data.FolderLocation = types.StringValue("")
	}

	if len(folderData.Records) > 0 {
		uids := make([]string, 0, len(folderData.Records))
		for _, r := range folderData.Records {
			if r.RecordUID != "" {
				uids = append(uids, r.RecordUID)
			}
		}
		recordSet, diags := types.SetValueFrom(ctx, types.StringType, uids)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Records = recordSet
	} else {
		data.Records = types.SetValueMust(types.StringType, []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
