// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamRemoteBrowserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state commonpamremotebrowser.PamRemoteBrowserResourceModel

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
		resp.Diagnostics.AddError(ErrSummaryPamRemoteBrowserReadFailed, "PAM remote browser record id is empty")
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonpamrecords.FetchVaultRecord(ctx, r.ApiManager, id)
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(utils.ErrSummaryFetchVaultRecordFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryPamRemoteBrowserReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamRemoteBrowser {
		resp.Diagnostics.AddError(
			ErrSummaryPamRemoteBrowserReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamRemoteBrowser),
		)
		return
	}

	resp.Diagnostics.Append(commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(ctx, &rec, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
