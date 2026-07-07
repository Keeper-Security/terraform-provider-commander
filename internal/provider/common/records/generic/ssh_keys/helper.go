// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for an sshKeys record.
func BuildAddCommand(cmd string, data SshKeysModel) string {
	var extra []string

	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagLogin, data.Login)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassphrase, data.Passphrase)
	appendOptionalHostAdd(&extra, data.Hostname, data.Port)
	appendOptionalKeyPairAdd(&extra, data.PublicKey, data.PrivateKey)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(
		cmd,
		commonrecordsutils.RecordTypeSshKeys,
		data.Title.ValueString(),
		data.FolderLocation,
		extra,
		custom,
		data.Notes,
	)
}

// UpdateHasMutations reports whether plan differs from state on updatable sshKeys fields.
func UpdateHasMutations(plan, state SshKeysModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Passphrase.Equal(state.Passphrase) ||
		!plan.Hostname.Equal(state.Hostname) ||
		!plan.Port.Equal(state.Port) ||
		!plan.PublicKey.Equal(state.PublicKey) ||
		!plan.PrivateKey.Equal(state.PrivateKey) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed sshKeys fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state SshKeysModel) string {
	var extra []string

	commonrecordsutils.AppendChangedStringField(&extra, FlagLogin, plan.Login, state.Login)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassphrase, plan.Passphrase, state.Passphrase)
	appendChangedHostUpdate(&extra, plan.Hostname, plan.Port, state.Hostname, state.Port)
	appendChangedKeyPairUpdate(&extra, plan.PublicKey, plan.PrivateKey, state.PublicKey, state.PrivateKey)

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

// MapVaultRecordGetResponseToSshKeysModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToSshKeysModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *SshKeysModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Login = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeLogin)
	m.Passphrase = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypePassword, PassphraseFieldLabel)

	if host := commonrecordsutils.HostFromFields(rec.Fields, ""); host != nil {
		m.Hostname = host.HostName
		m.Port = host.Port
	} else {
		m.Hostname = types.StringNull()
		m.Port = types.StringNull()
	}

	if kp := commonrecordsutils.KeyPairFromFields(rec.Fields, ""); kp != nil {
		m.PublicKey = kp.PublicKey
		m.PrivateKey = kp.PrivateKey
	} else {
		m.PublicKey = types.StringNull()
		m.PrivateKey = types.StringNull()
	}

	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}

func appendOptionalHostAdd(parts *[]string, hostname, port types.String) {
	if j, ok := hostJSON(hostname, port); ok {
		commonrecordsutils.AppendOptionalJSONAdd(parts, FlagHost, j)
	}
}

func appendOptionalKeyPairAdd(parts *[]string, publicKey, privateKey types.String) {
	if j, ok := keyPairJSON(publicKey, privateKey); ok {
		commonrecordsutils.AppendOptionalJSONAdd(parts, FlagKeyPair, j)
	}
}

func appendChangedHostUpdate(parts *[]string, planHostname, planPort, stateHostname, statePort types.String) {
	planJSON, planOK := hostJSON(planHostname, planPort)
	stateJSON, stateOK := hostJSON(stateHostname, statePort)
	changed := planJSON != stateJSON || planOK != stateOK
	commonrecordsutils.AppendChangedJSONField(parts, FlagHost, planJSON, stateJSON, changed)
}

func appendChangedKeyPairUpdate(parts *[]string, planPublic, planPrivate, statePublic, statePrivate types.String) {
	planJSON, planOK := keyPairJSON(planPublic, planPrivate)
	stateJSON, stateOK := keyPairJSON(statePublic, statePrivate)
	changed := planJSON != stateJSON || planOK != stateOK
	commonrecordsutils.AppendChangedJSONField(parts, FlagKeyPair, planJSON, stateJSON, changed)
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

func keyPairJSON(publicKey, privateKey types.String) (string, bool) {
	kp := &commonrecordsutils.KeyPairValue{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	if keyPairEmpty(kp) {
		return "", false
	}
	j, err := kp.ToJSON()
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

func keyPairEmpty(k *commonrecordsutils.KeyPairValue) bool {
	if k == nil {
		return true
	}
	return stringUnset(k.PublicKey) && stringUnset(k.PrivateKey)
}

func stringUnset(s types.String) bool {
	return s.IsNull() || s.IsUnknown() || strings.TrimSpace(s.ValueString()) == ""
}
