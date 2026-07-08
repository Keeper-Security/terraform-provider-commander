// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server_test

import (
	"encoding/json"
	"strings"
	"testing"

	commonrecordserver "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/server"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddCommand_IncludesServerFields(t *testing.T) {
	t.Parallel()

	data := commonrecordserver.ServerModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("Test Server Record Type"),
			Notes: types.StringValue("test notes"),
		},
		Login:    types.StringValue("root"),
		Password: types.StringValue("secret"),
		Hostname: types.StringValue("localhost"),
		Port:     types.StringValue("22"),
	}

	cmd := commonrecordserver.BuildAddCommand(utils.CmdRecordAdd, data)

	for _, want := range []string{
		"record-add",
		"--record-type serverCredentials",
		"login=",
		"password=",
		"f.host=",
		"root",
		"localhost",
		"test notes",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("command %q missing %q", cmd, want)
		}
	}
}

func TestMapVaultRecordGetResponseToServerModel(t *testing.T) {
	t.Parallel()

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "OjgmY-sRsNmu3tVFY7F4_A",
		Title:     "Test Server Record Type",
		Type:      "serverCredentials",
		Notes:     "test notes",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "host", Value: json.RawMessage(`[{"hostName":"localhost","port":"22"}]`)},
			{Type: "login", Value: json.RawMessage(`["root"]`)},
			{Type: "password", Value: json.RawMessage(`["k2j<ZC*W/YO8H5'}cE$Z"]`)},
		},
	}

	var state commonrecordserver.ServerModel
	commonrecordserver.MapVaultRecordGetResponseToServerModel(rec, types.StringNull(), &state)

	if state.Login.ValueString() != "root" {
		t.Fatalf("login = %q", state.Login.ValueString())
	}
	if state.Password.ValueString() != "k2j<ZC*W/YO8H5'}cE$Z" {
		t.Fatalf("password = %q", state.Password.ValueString())
	}
	if state.Hostname.ValueString() != "localhost" || state.Port.ValueString() != "22" {
		t.Fatalf("host = %q:%q", state.Hostname.ValueString(), state.Port.ValueString())
	}
	if state.Notes.ValueString() != "test notes" {
		t.Fatalf("notes = %q", state.Notes.ValueString())
	}
}
