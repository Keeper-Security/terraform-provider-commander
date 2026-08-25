// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordhealthinsurance "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/health_insurance"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *HealthInsuranceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError(utils.ERR_MSG_INVALID_IMPORT_ID, "Import ID cannot be empty. Use the health insurance record UID.")
		return
	}
	m := HealthInsuranceResourceModel{
		HealthInsuranceModel: commonrecordhealthinsurance.HealthInsuranceModel{
			BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
				Id:             types.StringValue(id),
				Title:          types.StringNull(),
				Notes:          types.StringNull(),
				FolderLocation: types.StringNull(),
			},
			AccountNumber:  types.StringNull(),
			Name:           nil,
			Login:          types.StringNull(),
			Password:       types.StringNull(),
			WebsiteAddress: types.StringNull(),
			Custom:         nil,
		},
		ShareModel: new_share.ShareModel{
			Share: types.MapNull(new_share.ShareEntryAttrType),
		},
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
