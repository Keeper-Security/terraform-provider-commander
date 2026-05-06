// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FolderResponse is the data payload from get FOLDER_UID --format json.
type FolderResponse struct {
	FolderUID string                 `json:"folder_uid"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Path      string                 `json:"path,omitempty"`
	Records   []FolderRecordResponse `json:"records,omitempty"`
}

// FolderRecordResponse represents a record entry in the folder get response.
type FolderRecordResponse struct {
	RecordUID string `json:"record_uid"`
}

func (r *FolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FolderResourceModel

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
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "folder id is empty")
		return
	}

	command := fmt.Sprintf("%s '%s' %s %s", CmdGet, id, FlagFormat, FormatJSON)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpGetFolder)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var folderData FolderResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &folderData); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if folderData.FolderUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(folderData.FolderUID)
	state.Name = types.StringValue(folderData.Name)

	if folderData.Path != "" {
		parent, _ := SplitFolderPath(folderData.Path)
		if parent != "" {
			state.FolderLocation = types.StringValue(parent)
		} else {
			state.FolderLocation = types.StringNull()
		}
	}

	if folderData.Records != nil {
		uids := make([]string, 0, len(folderData.Records))
		for _, r := range folderData.Records {
			if r.RecordUID != "" {
				uids = append(uids, r.RecordUID)
			}
		}
		if len(uids) > 0 {
			recordSet, diags := types.SetValueFrom(ctx, types.StringType, uids)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				state.Records = recordSet
			}
		} else {
			state.Records = types.SetNull(types.StringType)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
