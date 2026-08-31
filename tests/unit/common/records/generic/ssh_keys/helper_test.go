// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys_test

import (
	"encoding/json"
	"strings"
	"testing"

	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddCommand_IncludesSshKeysFields(t *testing.T) {
	t.Parallel()

	data := commonrecordsshkeys.SshKeysModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("ssh record"),
		},
		Login:      types.StringValue("user@example.com"),
		Passphrase: types.StringValue("secret-pass"),
		Hostname:   types.StringValue("12.0.0.1"),
		Port:       types.StringValue("22"),
		PublicKey:  types.StringValue("pub"),
		PrivateKey: types.StringValue("priv"),
	}

	cmd := commonrecordsshkeys.BuildAddCommand(utils.CmdRecordAdd, data)

	for _, want := range []string{
		"record-add",
		"--record-type sshKeys",
		"login=",
		"f.password.passphrase=",
		"f.host=",
		"f.keyPair=",
		"user@example.com",
		"12.0.0.1",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("command %q missing %q", cmd, want)
		}
	}
}

func TestMapVaultRecordGetResponseToSshKeysModel(t *testing.T) {
	t.Parallel()

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "uid-ssh-1",
		Title:     "ssh record",
		Type:      "sshKeys",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "login", Value: json.RawMessage(`["user@example.com"]`)},
			{Type: "password", Label: "passphrase", Value: json.RawMessage(`["secret"]`)},
			{Type: "host", Value: json.RawMessage(`[{"hostName":"12.0.0.1","port":"22"}]`)},
			{Type: "keyPair", Value: json.RawMessage(`[{"publicKey":"pub","privateKey":"priv"}]`)},
		},
	}

	var state commonrecordsshkeys.SshKeysModel
	commonrecordsshkeys.MapVaultRecordGetResponseToSshKeysModel(rec, types.StringNull(), &state)

	if state.Login.ValueString() != "user@example.com" {
		t.Fatalf("login = %q", state.Login.ValueString())
	}
	if state.Passphrase.ValueString() != "secret" {
		t.Fatalf("passphrase = %q", state.Passphrase.ValueString())
	}
	if state.Hostname.ValueString() != "12.0.0.1" || state.Port.ValueString() != "22" {
		t.Fatalf("host = %q:%q", state.Hostname.ValueString(), state.Port.ValueString())
	}
	if state.PublicKey.ValueString() != "pub" || state.PrivateKey.ValueString() != "priv" {
		t.Fatalf("keys = %q:%q", state.PublicKey.ValueString(), state.PrivateKey.ValueString())
	}
}
