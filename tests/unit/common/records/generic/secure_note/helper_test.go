// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote_test

import (
	"encoding/json"
	"strings"
	"testing"

	commonrecordsecurenote "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/secure_note"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddCommand_IncludesSecureNoteFields(t *testing.T) {
	t.Parallel()

	data := commonrecordsecurenote.SecureNoteModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("secure note record"),
			Notes: types.StringValue("optional management notes"),
		},
		SecuredNote: types.StringValue("hey this is test"),
		Date:        types.StringValue("2026-05-25"),
	}

	cmd := commonrecordsecurenote.BuildAddCommand(utils.CmdRecordAdd, data)

	for _, want := range []string{
		"record-add",
		"--record-type encryptedNotes",
		"f.note=",
		"f.date=",
		"hey this is test",
		"2026-05-25",
		"optional management notes",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("command %q missing %q", cmd, want)
		}
	}
}

func TestMapVaultRecordGetResponseToSecureNoteModel(t *testing.T) {
	t.Parallel()

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "SgIzZdtNw2m_5jsfT-OleQ",
		Title:     "secure note record",
		Type:      "encryptedNotes",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "note", Value: json.RawMessage(`["hey this is test"]`)},
			{Type: "date", Value: json.RawMessage(`[1779690823894]`)},
		},
	}

	var state commonrecordsecurenote.SecureNoteModel
	commonrecordsecurenote.MapVaultRecordGetResponseToSecureNoteModel(rec, types.StringNull(), &state)

	if state.SecuredNote.ValueString() != "hey this is test" {
		t.Fatalf("secured_note = %q", state.SecuredNote.ValueString())
	}
	if state.Date.ValueString() != "2026-05-25" {
		t.Fatalf("date = %q, want 2026-05-25", state.Date.ValueString())
	}
}
