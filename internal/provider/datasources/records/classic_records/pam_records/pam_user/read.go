// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	pamUserResource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_user"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamUserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	recordUID := strings.TrimSpace(data.RecordUID.ValueString())
	if recordUID == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "record_uid is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	// Phase 1: fetch the vault record via `get '<uid>' --format json`.
	command := fmt.Sprintf("%s '%s' %s", utils.CmdGetRecord, recordUID, utils.FlagFormatJSON)
	apiResp, err := d.ApiManager.ExecuteCommand(ctx, command, ErrDetailReadFailed)
	if err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, fmt.Sprintf("record %q not found or empty response", recordUID))
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamUser {
		resp.Diagnostics.AddError(
			ErrSummaryReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamUser),
		)
		return
	}

	mapVaultRecordToDataSource(&rec, &data)

	// Phase 2: fetch rotation info via `pam rotation info -r '<uid>'`.
	rotCmd := fmt.Sprintf("%s -r '%s'", pamUserResource.CmdPamRotationInfo, recordUID)
	rotResp, err2 := d.ApiManager.ExecuteCommand(ctx, rotCmd, ErrDetailRotationInfoFailed)
	if err2 == nil && rotResp != nil {
		messages := parseFlexibleMessageToLines(rotResp.Message.String())
		if hasRotationData(messages) {
			rs := parseRotationInfoToDataSource(messages)
			data.RotationSettings = rs
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapVaultRecordToDataSource(rec *utils.VaultRecordGetResponse, data *PamUserDataSourceModel) {
	if strings.TrimSpace(rec.RecordUID) != "" {
		data.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}

	data.Title = stringOrNull(rec.Title)
	data.Notes = stringOrNull(rec.Notes)
	data.Folder = commonpamrecords.ExtractFolderValue(rec.Folder, data.Folder)

	for i := range rec.Fields {
		f := &rec.Fields[i]
		switch f.Type {
		case "login":
			data.Login = firstStringValue(f.Value)
		case "password":
			data.Password = firstStringValue(f.Value)
		case "text":
			switch f.Label {
			case "distinguishedName":
				data.DistinguishedName = firstStringValue(f.Value)
			case "connectDatabase":
				data.ConnectDatabase = firstStringValue(f.Value)
			}
		case "secret":
			if f.Label == "privatePEMKey" {
				data.PrivatePEMKey = firstStringValue(f.Value)
			}
		case "checkbox":
			if f.Label == "managed" {
				data.Managed = firstBoolValue(f.Value)
			}
		}
	}
}

func stringOrNull(s string) types.String {
	if strings.TrimSpace(s) == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func firstStringValue(raw json.RawMessage) types.String {
	var vals []string
	if err := json.Unmarshal(raw, &vals); err != nil {
		return types.StringNull()
	}
	if len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
		return types.StringValue(vals[0])
	}
	return types.StringNull()
}

func firstBoolValue(raw json.RawMessage) types.Bool {
	var vals []bool
	if err := json.Unmarshal(raw, &vals); err != nil {
		return types.BoolNull()
	}
	if len(vals) > 0 {
		return types.BoolValue(vals[0])
	}
	return types.BoolNull()
}

func hasRotationData(messages []string) bool {
	for _, line := range messages {
		if strings.HasPrefix(strings.TrimSpace(line), "PAM Config UID:") {
			return true
		}
	}
	return false
}

func parseRotationInfoToDataSource(messages []string) *PamUserDataSourceRotationSettings {
	rs := &PamUserDataSourceRotationSettings{}

	for _, line := range messages {
		line = strings.TrimSpace(line)
		if k, v, ok := strings.Cut(line, ": "); ok {
			switch k {
			case "PAM Config UID":
				rs.Configuration = stringOrNull(v)
			case "Admin Resource Uid":
				rs.AdminUser = stringOrNull(v)
			case "Is Rotation Disabled":
				rs.Enabled = types.BoolValue(strings.EqualFold(strings.TrimSpace(v), "False"))
			case "Schedule Type":
				if strings.Contains(strings.ToLower(v), "manual") {
					rs.OnDemand = types.BoolValue(true)
				}
			case "Schedule":
				parseScheduleValue(strings.TrimSpace(v), rs)
			case "Password Complexity Data":
				rs.Complexity = stringOrNull(v)
			}
		}
	}

	return rs
}

func parseScheduleValue(raw string, rs *PamUserDataSourceRotationSettings) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}

	var sched map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &sched); err != nil {
		return
	}

	schedType, _ := sched["type"].(string)
	switch strings.ToUpper(schedType) {
	case "CRON":
		if cron, ok := sched["cron"].(string); ok {
			rs.ScheduleCron = types.StringValue(cron)
			rs.OnDemand = types.BoolNull()
			rs.ScheduleJSON = types.StringNull()
		}
	case "DAILY", "WEEKLY", "MONTHLY":
		rs.ScheduleJSON = types.StringValue(raw)
		rs.OnDemand = types.BoolNull()
		rs.ScheduleCron = types.StringNull()
	}
}

func parseFlexibleMessageToLines(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err == nil {
		return arr
	}

	return strings.Split(raw, "\n")
}
