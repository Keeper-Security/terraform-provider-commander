// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"
	"strings"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *WifiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError(utils.ERR_MSG_INVALID_IMPORT_ID, "Import ID cannot be empty. Use the WiFi record UID.")
		return
	}
	m := WifiResourceModel{
		BaseVaultRecordModel: records.BaseVaultRecordModel{
			Id:     types.StringValue(id),
			Title:  types.StringNull(),
			Notes:  types.StringNull(),
			Folder: types.StringNull(),
			Custom: nil,
		},
		SSID:         types.StringNull(),
		Password:     types.StringNull(),
		Encryption:   types.StringNull(),
		IsSSIDHidden: types.BoolNull(),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
