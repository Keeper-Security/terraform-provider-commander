// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"
	"fmt"
	"strings"

	commonnewfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/new_folder"
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
	command := buildCreateNewFolderCommand(&data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpCreate)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	data.Id = types.StringValue(string(apiResp.Message))

	// Sync the share permissions - share the folder to the users in the share block.
	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdNsfShareFolder, data.Id.ValueString(), data.Share, types.MapNull(new_share.ShareEntryAttrType)); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	// Link the records to the folder.
	if err := folderutils.LinkRecords(ctx, r.ApiManager, CmdNsfLn, data.Id.ValueString(), data.Records); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	// `records` is Optional+Computed: resolve it to the real, known value now
	// (re-fetch rather than trust the plan) so an omitted `records` doesn't
	// leave the attribute unknown after apply.
	apiData, err := commonnewfolder.FetchNsfFolderByNameOrId(ctx, r.ApiManager, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}
	if err := commonnewfolder.MapResponseToModel(ctx, apiData, &data); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}
	if data.Records.IsUnknown() {
		data.Records = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// It creates the nsf directory command and the folder path used as NAME.
func buildCreateNewFolderCommand(data *NewFolderResourceModel) (command string) {
	name := data.Name.ValueString()

	// Build the folder path: if folder_location is set, "folder_location/name", otherwise "name".
	folderPath := folderutils.BuildFolderPath(name, data.FolderLocation.ValueString())
	parts := []string{CmdNsfMkdir, fmt.Sprintf(`"%s"`, folderutils.EscapeDoubleQuotesForCLI(folderPath))}

	command = strings.Join(parts, " ")
	return command
}
