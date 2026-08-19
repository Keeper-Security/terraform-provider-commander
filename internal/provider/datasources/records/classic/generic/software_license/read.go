// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsoftwarelicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/software_license"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *SoftwareLicenseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SoftwareLicenseDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	lookup := strings.TrimSpace(data.SoftwareLicense.ValueString())
	if lookup == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "Software license record identifier is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonrecordsutils.FetchVaultRecord(ctx, d.ApiManager, lookup)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}
	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("Software license record '%s' not found", lookup))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}
	if rec.Type != "" && rec.Type != commonrecordsutils.RecordTypeSoftwareLicense {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("vault record type is %q, expected %q", rec.Type, commonrecordsutils.RecordTypeSoftwareLicense))
		return
	}

	resp.Diagnostics.Append(commonrecordsoftwarelicense.MapVaultRecordGetResponseToSoftwareLicenseModel(&rec, data.FolderLocation, &data.SoftwareLicenseModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := classic_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadSoftwareLicenseDataSource, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
