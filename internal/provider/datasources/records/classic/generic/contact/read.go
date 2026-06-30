// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *ContactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContactDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	lookup := strings.TrimSpace(data.Contact.ValueString())
	if lookup == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "Contact record identifier is empty")
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
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("Contact record '%s' not found", lookup))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}
	if rec.Type != "" && rec.Type != utils.RecordTypeContact {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypeContact))
		return
	}

	data.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	data.Title = utils.StringOrNull(rec.Title)
	data.Notes = utils.StringOrNull(rec.Notes)
	data.FolderLocation = utils.ExtractFolderValue(rec.FolderLocation, types.StringNull())
	data.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	data.Company = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, "company")
	data.Email = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeEmail, "")
	data.Phone = commonrecordsutils.PhonesFromField(rec.Fields, "")
	data.AddressRef = commonrecordsutils.FirstRefUID(rec.Fields, commonrecordsutils.FieldTypeAddressRef, "")
	data.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)

	if err := classic_share.MapResponseToModel(rec.UserPermissions, &data.ShareModel); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadContactDataSource, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
