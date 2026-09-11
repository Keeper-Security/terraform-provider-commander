// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordwifi "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/wifi"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *WifiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WifiDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	lookup := strings.TrimSpace(data.Wifi.ValueString())
	if lookup == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "wifi is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonrecordsutils.FetchNsfVaultRecord(ctx, d.ApiManager, lookup)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}
	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("record %q not found or empty response", lookup))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != commonrecordsutils.RecordTypeWifiCredentials {
		resp.Diagnostics.AddError(
			ErrSummaryReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, commonrecordsutils.RecordTypeWifiCredentials),
		)
		return
	}

	resp.Diagnostics.Append(commonrecordwifi.MapVaultRecordGetResponseToWifiModel(&rec, data.FolderLocation, &data.WifiModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := new_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
