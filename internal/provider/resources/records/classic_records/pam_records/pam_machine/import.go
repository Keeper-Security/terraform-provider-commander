// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_machine"
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

	state := commonpammachine.PamMachineResourceModel{
		CommonPamRecordsResourceModel: commonpamrecords.CommonPamRecordsResourceModel{
			Id:     types.StringValue(importID),
			Title:  types.StringNull(),
			Notes:  types.StringNull(),
			Folder: types.StringNull(),
		},
		HostnameOrIP:    nil,
		OperatingSystem: types.StringNull(),
		InstanceName:    types.StringNull(),
		InstanceId:      types.StringNull(),
		ProviderGroup:   types.StringNull(),
		ProviderRegion:  types.StringNull(),
		PamSettings:     nil,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
