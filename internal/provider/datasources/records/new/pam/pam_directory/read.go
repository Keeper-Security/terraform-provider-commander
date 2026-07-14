// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_directory"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamDirectoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamDirectoryDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	recordUID := strings.TrimSpace(data.PamDirectory.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(errSummaryReadPamDirectoryDataSource, "pam_directory is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonrecordsutils.FetchVaultRecord(ctx, d.ApiManager, recordUID)
	if err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamDirectoryDataSource, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(errSummaryReadPamDirectoryDataSource, fmt.Sprintf("record %q not found or empty response", recordUID))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamDirectoryDataSource, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamDirectory {
		resp.Diagnostics.AddError(
			errSummaryReadPamDirectoryDataSource,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamDirectory),
		)
		return
	}

	var state commonpamdirectory.PamDirectoryResourceModel
	resp.Diagnostics.Append(commonpamdirectory.MapVaultRecordGetResponseToPamDirectoryModel(&rec, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.PamDirectory = types.StringValue(recordUID)
	data.Id = state.Id
	data.Title = state.Title
	data.Notes = state.Notes
	data.FolderLocation = state.FolderLocation
	data.HostnameOrIP = state.HostnameOrIP
	data.UseSSL = state.UseSSL
	data.DomainName = state.DomainName
	data.AlternativeIPs = state.AlternativeIPs
	data.DirectoryId = state.DirectoryId
	data.DirectoryType = state.DirectoryType
	data.UserMatch = state.UserMatch
	data.ProviderGroup = state.ProviderGroup
	data.ProviderRegion = state.ProviderRegion
	data.PamSettings = state.PamSettings

	if err := new_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamDirectoryDataSource, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
