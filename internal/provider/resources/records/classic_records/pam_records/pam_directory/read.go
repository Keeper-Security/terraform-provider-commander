// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamDirectoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PamDirectoryResourceModel

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
		resp.Diagnostics.AddError(ErrSummaryPamDirectoryReadFailed, "PAM directory record id is empty")
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
		resp.Diagnostics.AddError(ErrSummaryPamDirectoryReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamDirectory {
		resp.Diagnostics.AddError(
			ErrSummaryPamDirectoryReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamDirectory),
		)
		return
	}

	resp.Diagnostics.Append(commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(&rec, &state.PamDirectoryResourceModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := classic_share.MapResponseToModel(rec.UserPermissions, &state.ShareModel); err != nil {
		resp.Diagnostics.AddError(ErrSummaryPamDirectoryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
