// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"
	"strings"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *AddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError(utils.ERR_MSG_INVALID_IMPORT_ID, "Import ID cannot be empty. Use the Address record UID.")
		return
	}
	m := AddressResourceModel{
		BaseVaultRecordModel: records.BaseVaultRecordModelImportValues(id),
		Street1:              types.StringNull(),
		Street2:              types.StringNull(),
		City:                 types.StringNull(),
		State:                types.StringNull(),
		Zip:                  types.StringNull(),
		Country:              types.StringNull(),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
