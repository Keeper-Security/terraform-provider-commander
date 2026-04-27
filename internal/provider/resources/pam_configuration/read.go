// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"
	"errors"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state commonpamconfiguration.PamConfigurationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	id := strings.TrimSpace(state.Id.ValueString())
	if id == "" {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, "pam configuration id is empty")
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(ErrSummarySyncDownFailed, err.Error())
		return
	}

	command := commonpamconfiguration.FetchPamConfigByUIDCommand(id)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, commonpamconfiguration.ErrOpFetchPamConfig)
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var apiData utils.PamConfigListResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &apiData); err != nil {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, err.Error())
		return
	}

	if err := commonpamconfiguration.MapPamConfigAPIResponseToModel(&state, &apiData); err != nil {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
