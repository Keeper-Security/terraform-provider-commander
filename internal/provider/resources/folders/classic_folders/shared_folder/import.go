// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"
	"strings"

	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import by classic shared folder UID.
// Import ID: the classic shared folder UID (e.g. "BTbjhOmqw9iYal3OQJ9UAQ").
// After import, Terraform runs Read to refresh state from the API.
func (r *ClassicSharedFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the classic shared folder name or UID.",
		)
		return
	}

	state := SharedFolderResourceModel{
		CommonFolderModel: folderutils.CommonFolderModel{
			Id:             types.StringValue(importID),
			Name:           types.StringNull(),
			FolderLocation: types.StringNull(),
		},
		UserPermissions:   nil,
		RecordPermissions: nil,
		Records:           types.MapNull(RecordEntryMapElemType),
		Users:             types.MapNull(UserEntryMapElemType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
