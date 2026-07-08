// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a serverCredentials record.
func BuildAddCommand(cmd string, data ServerModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagLogin, data.Login)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)
	appendOptionalHostAdd(&extra, data.Hostname, data.Port)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(
		cmd,
		commonrecordsutils.RecordTypeServerCredentials,
		data.Title.ValueString(),
		data.FolderLocation,
		extra,
		custom,
		data.Notes,
	)
}

// UpdateHasMutations reports whether plan differs from state on updatable server fields.
func UpdateHasMutations(plan, state ServerModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Password.Equal(state.Password) ||
		!plan.Hostname.Equal(state.Hostname) ||
		!plan.Port.Equal(state.Port) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed server fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state ServerModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagLogin, plan.Login, state.Login)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)
	appendChangedHostUpdate(&extra, plan.Hostname, plan.Port, state.Hostname, state.Port)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(
		cmd,
		recordUID,
		plan.Title,
		state.Title,
		extra,
		customPlan,
		customState,
		plan.Notes,
		state.Notes,
	)
}

// MapVaultRecordGetResponseToServerModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToServerModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *ServerModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Login = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeLogin)
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)

	if host := commonrecordsutils.HostFromFields(rec.Fields, ""); host != nil {
		m.Hostname = host.HostName
		m.Port = host.Port
	} else {
		m.Hostname = types.StringNull()
		m.Port = types.StringNull()
	}

	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}

func appendOptionalHostAdd(parts *[]string, hostname, port types.String) {
	if j, ok := hostJSON(hostname, port); ok {
		commonrecordsutils.AppendOptionalJSONAdd(parts, FlagHost, j)
	}
}

func appendChangedHostUpdate(parts *[]string, planHostname, planPort, stateHostname, statePort types.String) {
	planJSON, planOK := hostJSON(planHostname, planPort)
	stateJSON, stateOK := hostJSON(stateHostname, statePort)
	changed := planJSON != stateJSON || planOK != stateOK
	commonrecordsutils.AppendChangedJSONField(parts, FlagHost, planJSON, stateJSON, changed)
}

func hostJSON(hostname, port types.String) (string, bool) {
	host := &commonrecordsutils.HostValue{
		HostName: hostname,
		Port:     port,
	}
	if hostValueEmpty(host) {
		return "", false
	}
	j, err := host.ToJSON()
	if err != nil || strings.TrimSpace(j) == "" {
		return "", false
	}
	return j, true
}

func hostValueEmpty(h *commonrecordsutils.HostValue) bool {
	if h == nil {
		return true
	}
	return stringUnset(h.HostName) && stringUnset(h.Port)
}

func stringUnset(s types.String) bool {
	return s.IsNull() || s.IsUnknown() || strings.TrimSpace(s.ValueString()) == ""
}
