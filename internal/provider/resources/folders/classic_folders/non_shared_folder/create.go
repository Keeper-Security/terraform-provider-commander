// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"
	"fmt"
	"strings"

	commonnonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/classic_folders/non_shared_folder"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *NonSharedFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NonSharedFolderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	command := buildCreateNonSharedFolderCommand(&data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, folderutils.ErrOpCreate)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	folderUID, err := folderutils.ExtractFolderUIDFromCreateResponse(apiResp.Data)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}
	data.Id = types.StringValue(folderUID)

	if err := folderutils.LinkRecords(ctx, r.ApiManager, CmdLn, folderUID, data.Records); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}

	// `records` is Optional+Computed: resolve it to the real, known value now
	// (re-fetch rather than trust the plan) so an omitted `records` doesn't
	// leave the attribute unknown after apply.
	apiData, err := commonnonsharedfolder.FetchFolderByNameOrId(ctx, r.ApiManager, folderUID)
	if err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}
	if err := commonnonsharedfolder.MapResponseToModel(ctx, apiData, &data); err != nil {
		resp.Diagnostics.AddError(folderutils.ErrSummaryCreateFailed, err.Error())
		return
	}
	if data.Records.IsUnknown() {
		data.Records = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildCreateNonSharedFolderCommand(data *NonSharedFolderResourceModel) string {

	folderPath := folderutils.BuildFolderPath(data.Name.ValueString(), data.FolderLocation.ValueString())
	parts := []string{CmdMkdir, FlagUserFolder, fmt.Sprintf(`"%s"`, folderutils.EscapeDoubleQuotesForCLI(folderPath))}

	return strings.Join(parts, " ")
}
