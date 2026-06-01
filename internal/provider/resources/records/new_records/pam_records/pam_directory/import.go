// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamDirectoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the new PAM directory record UID when defined.",
		)
		return
	}

	state := PamDirectoryResourceModel{
		PamDirectoryResourceModel: commonpamdirectory.PamDirectoryResourceModel{
			CommonPamRecordsResourceModel: commonpamrecords.CommonPamRecordsResourceModel{
				Id:     types.StringValue(importID),
				Title:  types.StringNull(),
				Notes:  types.StringNull(),
				Folder: types.StringNull(),
			},
			HostnameOrIP:   nil,
			UseSSL:         types.BoolNull(),
			DomainName:     types.StringNull(),
			AlternativeIPs: types.SetNull(types.StringType),
			DirectoryId:    types.StringNull(),
			DirectoryType:  types.StringNull(),
			UserMatch:      types.StringNull(),
			ProviderGroup:  types.StringNull(),
			ProviderRegion: types.StringNull(),
			PamSettings:    nil,
		},
		ShareModel: new_share.ShareModel{
			Share: types.MapNull(new_share.ShareEntryAttrType),
		},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
