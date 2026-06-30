// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"
	"fmt"
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	apiResp, err := commonrecordsutils.FetchVaultRecord(ctx, d.ApiManager, lookup)
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

	resp.Diagnostics.Append(mapVaultRecordToDataSource(ctx, &rec, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapVaultRecordToDataSource hydrates the data source model from a `get <uid> --format json` payload.
// `SSID` is matched by the `SSID` label since `text` is a generic field type. The remaining
// built-in fields are matched by type alone — Keeper returns label="" or label=<type>
// inconsistently, and either form is supported by the AnyLabel helpers.
func mapVaultRecordToDataSource(ctx context.Context, rec *utils.VaultRecordGetResponse, data *WifiDataSourceModel) diag.Diagnostics {
	if uid := strings.TrimSpace(rec.RecordUID); uid != "" {
		data.Id = types.StringValue(uid)
	}
	data.Title = utils.StringOrNull(rec.Title)
	data.Notes = utils.StringOrNull(rec.Notes)
	data.FolderLocation = utils.ExtractFolderValue(rec.FolderLocation, data.FolderLocation)

	data.SSID = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, "SSID")
	data.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	data.Encryption = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeWifiEncryption)
	data.IsSSIDHidden = commonrecordsutils.FirstBoolFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeIsSSIDHidden)

	data.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
