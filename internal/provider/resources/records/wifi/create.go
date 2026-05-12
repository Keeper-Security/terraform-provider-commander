// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *WifiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WifiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}
	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	cmd := buildRecordAddCommand(data)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailCreateFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, err.Error())
		return
	}

	createdUID, ok := apiResp.Data.(map[string]interface{})["record_uid"].(string)
	if !ok {
		resp.Diagnostics.AddError(ErrSummaryCreateFailed, fmt.Sprintf("Failed to extract record_uid from response: %v", apiResp.Data))
		return
	}
	data.Id = types.StringValue(createdUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
