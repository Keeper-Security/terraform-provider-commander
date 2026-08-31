// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database_test

import (
	"encoding/json"
	"strings"
	"testing"

	commonrecorddatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/database"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddCommand_IncludesDatabaseFields(t *testing.T) {
	t.Parallel()

	data := commonrecorddatabase.DatabaseModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("test DB"),
			Notes: types.StringValue("test notes"),
		},
		Login:    types.StringValue("root"),
		Hostname: types.StringValue("127.0.0.1"),
		Port:     types.StringValue("8000"),
		Type:     types.StringValue("SQL"),
		Password: types.StringValue("secret"),
	}

	cmd := commonrecorddatabase.BuildAddCommand(utils.CmdRecordAdd, data)

	for _, want := range []string{
		"record-add",
		"--record-type databaseCredentials",
		"login=",
		"password=",
		"f.host=",
		"f.text.type=",
		"root",
		"127.0.0.1",
		"SQL",
		"test notes",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("command %q missing %q", cmd, want)
		}
	}
}

func TestMapVaultRecordGetResponseToDatabaseModel(t *testing.T) {
	t.Parallel()

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "sy7czlOU1MoI2Mcrbi9LNw",
		Title:     "test DB",
		Type:      "databaseCredentials",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "text", Label: "type", Value: json.RawMessage(`["SQL"]`)},
			{Type: "host", Value: json.RawMessage(`[{"hostName":"127.0.0.1","port":"8000"}]`)},
			{Type: "login", Value: json.RawMessage(`["root"]`)},
			{Type: "password", Value: json.RawMessage(`["XvsFZ'@UDC&o6>DuS\">%"]`)},
		},
	}

	var state commonrecorddatabase.DatabaseModel
	commonrecorddatabase.MapVaultRecordGetResponseToDatabaseModel(rec, types.StringNull(), &state)

	if state.Login.ValueString() != "root" {
		t.Fatalf("login = %q", state.Login.ValueString())
	}
	if state.Password.ValueString() != "XvsFZ'@UDC&o6>DuS\">%" {
		t.Fatalf("password = %q", state.Password.ValueString())
	}
	if state.Hostname.ValueString() != "127.0.0.1" || state.Port.ValueString() != "8000" {
		t.Fatalf("host = %q:%q", state.Hostname.ValueString(), state.Port.ValueString())
	}
	if state.Type.ValueString() != "SQL" {
		t.Fatalf("type = %q", state.Type.ValueString())
	}
}
