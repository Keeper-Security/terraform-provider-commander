// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
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
	commonrecordsutils.AppendOptionalHostFlatAdd(&extra, FlagHost, data.Hostname, data.Port)
	commonrecordsutils.AppendOptionalKeyPairFlatAdd(&extra, FlagKeyPair, data.PublicKey, data.PrivateKey)

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
	commonrecordsutils.AppendChangedHostFlatUpdate(&extra, FlagHost, plan.Hostname, plan.Port, state.Hostname, state.Port)
	commonrecordsutils.AppendChangedKeyPairFlatUpdate(&extra, FlagKeyPair, plan.PublicKey, plan.PrivateKey, state.PublicKey, state.PrivateKey)

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
	m.Hostname, m.Port = commonrecordsutils.FlatHostFromFields(rec.Fields)
	m.PublicKey, m.PrivateKey = commonrecordsutils.FlatKeyPairFromFields(rec.Fields)
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
