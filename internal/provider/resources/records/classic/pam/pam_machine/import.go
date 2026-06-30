// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the PAM machine record UID when defined.",
		)
		return
	}

	state := PamMachineResourceModel{
		PamMachineResourceModel: commonpammachine.PamMachineResourceModel{
			BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
				Id:             types.StringValue(importID),
				Title:          types.StringNull(),
				Notes:          types.StringNull(),
				FolderLocation: types.StringNull(),
			},
			HostnameOrIP:    nil,
			OperatingSystem: types.StringNull(),
			InstanceName:    types.StringNull(),
			InstanceId:      types.StringNull(),
			ProviderGroup:   types.StringNull(),
			ProviderRegion:  types.StringNull(),
			PamSettings:     nil,
		},
		ShareModel: classic_share.ShareModel{
			Share: types.MapNull(classic_share.ShareEntryAttrType),
		},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
