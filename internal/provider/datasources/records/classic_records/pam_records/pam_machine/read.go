// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"
	"fmt"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamMachineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamMachineDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	recordUID := strings.TrimSpace(data.PamMachine.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(errSummaryReadPamMachineDataSource, "pam_machine is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	apiResp, err := commonpamrecords.FetchVaultRecord(ctx, d.ApiManager, recordUID)
	if err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamMachineDataSource, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(errSummaryReadPamMachineDataSource, fmt.Sprintf("record %q not found or empty response", recordUID))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(errSummaryReadPamMachineDataSource, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamMachine {
		resp.Diagnostics.AddError(
			errSummaryReadPamMachineDataSource,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamMachine),
		)
		return
	}

	var state commonpammachine.PamMachineResourceModel
	resp.Diagnostics.Append(commonpammachine.MapVaultRecordGetResponseToPamMachineModel(&rec, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.PamMachine = types.StringValue(recordUID)
	data.Id = state.Id
	data.Title = state.Title
	data.Notes = state.Notes
	data.Folder = state.Folder
	data.HostnameOrIP = state.HostnameOrIP
	data.OperatingSystem = state.OperatingSystem
	data.InstanceName = state.InstanceName
	data.InstanceId = state.InstanceId
	data.ProviderGroup = state.ProviderGroup
	data.ProviderRegion = state.ProviderRegion
	data.PamSettings = state.PamSettings

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
