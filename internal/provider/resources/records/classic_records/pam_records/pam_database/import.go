// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the PAM database record UID when defined.",
		)
		return
	}

	state := PamDatabaseResourceModel{
		PamDatabaseResourceModel: commonpamdatabase.PamDatabaseResourceModel{
			CommonPamRecordsResourceModel: commonpamrecords.CommonPamRecordsResourceModel{
				Id:     types.StringValue(importID),
				Title:  types.StringNull(),
				Notes:  types.StringNull(),
				Folder: types.StringNull(),
			},
			HostnameOrIP:   nil,
			UseSSL:         types.BoolNull(),
			DatabaseId:     types.StringNull(),
			DatabaseType:   types.StringNull(),
			ProviderGroup:  types.StringNull(),
			ProviderRegion: types.StringNull(),
			PamSettings:    nil,
		},
		ShareModel: classic_share.ShareModel{
			Share: types.MapNull(classic_share.ShareEntryAttrType),
		},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
