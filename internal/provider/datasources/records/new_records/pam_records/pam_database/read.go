// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdatabase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_database"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamDatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamDatabaseDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	recordUID := strings.TrimSpace(data.PamDatabase.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(errSummaryReadPamDatabaseDataSource, "pam_database is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonrecordsutils.FetchVaultRecord(ctx, d.ApiManager, recordUID)
	if err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamDatabaseDataSource, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(errSummaryReadPamDatabaseDataSource, fmt.Sprintf("record %q not found or empty response", recordUID))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamDatabaseDataSource, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamDatabase {
		resp.Diagnostics.AddError(
			errSummaryReadPamDatabaseDataSource,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamDatabase),
		)
		return
	}

	var state commonpamdatabase.PamDatabaseResourceModel
	resp.Diagnostics.Append(commonpamdatabase.MapVaultRecordGetResponseToPamDatabaseModel(&rec, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.PamDatabase = types.StringValue(recordUID)
	data.Id = state.Id
	data.Title = state.Title
	data.Notes = state.Notes
	data.FolderLocation = state.FolderLocation
	data.HostnameOrIP = state.HostnameOrIP
	data.UseSSL = state.UseSSL
	data.DatabaseId = state.DatabaseId
	data.DatabaseType = state.DatabaseType
	data.ProviderGroup = state.ProviderGroup
	data.ProviderRegion = state.ProviderRegion
	data.PamSettings = state.PamSettings

	if err := new_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamDatabaseDataSource, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
