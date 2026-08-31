// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecorddriverlicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/driver_license"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *DriverLicenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DriverLicenseResourceModel
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
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "Driver's License record id is empty")
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	// Fetch the vault record.
	apiResp, err := commonrecordsutils.FetchVaultRecord(ctx, r.ApiManager, id)
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}
	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Map the vault record to the state model.
	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}
	if rec.Type != "" && rec.Type != commonrecordsutils.RecordTypeDriverLicense {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("vault record type is %q, expected %q", rec.Type, commonrecordsutils.RecordTypeDriverLicense))
		return
	}

	resp.Diagnostics.Append(commonrecorddriverlicense.MapVaultRecordGetResponseToDriverLicenseModel(&rec, state.FolderLocation, &state.DriverLicenseModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := classic_share.MapResponseToModel(rec.UserPermissions, &state.ShareModel); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
