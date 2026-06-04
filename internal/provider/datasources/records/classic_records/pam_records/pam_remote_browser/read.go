// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamRemoteBrowserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamRemoteBrowserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	recordUID := strings.TrimSpace(data.RemoteBrowser.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(errSummaryReadPamRemoteBrowserDataSource, "record_uid is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	// Fetch the vault record with DAG.
	apiResp, err := commonpamrecords.FetchVaultRecord(ctx, d.ApiManager, recordUID)
	if err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamRemoteBrowserDataSource, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(errSummaryReadPamRemoteBrowserDataSource, fmt.Sprintf("record %q not found or empty response", recordUID))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamRemoteBrowserDataSource, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamRemoteBrowser {
		resp.Diagnostics.AddError(
			errSummaryReadPamRemoteBrowserDataSource,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamRemoteBrowser),
		)
		return
	}

	var mapped commonpamremotebrowser.PamRemoteBrowserResourceModel
	resp.Diagnostics.Append(commonpamremotebrowser.MapVaultRecordGetResponseToPamRemoteBrowserModel(ctx, &rec, &mapped)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.RemoteBrowser = types.StringValue(recordUID)
	data.Id = mapped.Id
	data.Title = mapped.Title
	data.Url = mapped.Url
	data.Notes = mapped.Notes
	data.Folder = mapped.Folder
	data.PamRemoteBrowserSettings = mapped.PamRemoteBrowserSettings

	if err := classic_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamRemoteBrowserDataSource, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
