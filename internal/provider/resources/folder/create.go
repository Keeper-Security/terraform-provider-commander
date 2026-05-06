// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FolderResourceModel

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
		resp.Diagnostics.AddError(ErrSummarySyncDownFailed, err.Error())
		return
	}

	command, err := buildCreateFolderCommand(&data)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryInvalidConfig, err.Error())
		return
	}

	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrOpCreateFolder)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	folderUID, err := extractFolderUIDFromCreateResponse(apiResp.Data)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}
	data.Id = types.StringValue(folderUID)

	if err := LinkRecords(ctx, r.ApiManager, folderUID, data.Records); err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildCreateFolderCommand(data *FolderResourceModel) (string, error) {
	name := data.Name.ValueString()
	if name == "" {
		return "", fmt.Errorf("name is required")
	}

	folderPath := BuildFolderPath(name, data.FolderLocation.ValueString())
	parts := []string{CmdMkdir, FlagUserFolder, fmt.Sprintf(`"%s"`, EscapeDoubleQuotesForCLI(folderPath))}

	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		if color != "" {
			parts = append(parts, FlagColor, color)
		}
	}

	return strings.Join(parts, " "), nil
}

func extractFolderUIDFromCreateResponse(data any) (string, error) {
	m, _ := data.(map[string]interface{})
	v, ok := m[KeyFolderUID]
	if !ok || v == nil {
		return "", fmt.Errorf("API response missing %s", KeyFolderUID)
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return "", fmt.Errorf("API response %s is empty", KeyFolderUID)
		}
		return s, nil
	}
	return fmt.Sprintf("%v", v), nil
}
