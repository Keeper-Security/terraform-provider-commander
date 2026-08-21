// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense_test

import (
	"encoding/json"
	"strings"
	"testing"

	commonrecordsoftwarelicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/software_license"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddCommand_IncludesSoftwareLicenseFields(t *testing.T) {
	t.Parallel()

	data := commonrecordsoftwarelicense.SoftwareLicenseModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("my metron software record"),
			Notes: types.StringValue("this is test record"),
		},
		SoftwareLicenseKey: types.StringValue("132456789807867564331423565"),
		ExpirationDate:     types.StringValue("2026-05-20"),
		DateActive:         types.StringValue("2026-07-09"),
	}

	cmd := commonrecordsoftwarelicense.BuildAddCommand(utils.CmdRecordAdd, data)

	for _, want := range []string{
		"record-add",
		"--record-type softwareLicense",
		"f.licenseNumber=",
		"f.expirationDate=",
		"f.date.dateActive=",
		"132456789807867564331423565",
		"2026-05-20",
		"2026-07-09",
		"this is test record",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("command %q missing %q", cmd, want)
		}
	}
}

func TestMapVaultRecordGetResponseToSoftwareLicenseModel_ClassicEpochMillis(t *testing.T) {
	t.Parallel()

	rec := &utils.VaultRecordGetResponse{
		RecordUID: "u5M3yoJUn_gI0rvfBx4zXQ",
		Title:     "my metron software record",
		Type:      "softwareLicense",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "licenseNumber", Value: json.RawMessage(`["132456789807867564331423565"]`)},
			{Type: "expirationDate", Value: json.RawMessage(`[1779258619809]`)},
			{Type: "date", Label: "dateActive", Value: json.RawMessage(`[1783598400000]`)},
		},
	}

	var state commonrecordsoftwarelicense.SoftwareLicenseModel
	commonrecordsoftwarelicense.MapVaultRecordGetResponseToSoftwareLicenseModel(rec, types.StringNull(), &state)

	if state.SoftwareLicenseKey.ValueString() != "132456789807867564331423565" {
		t.Fatalf("software_license_key = %q", state.SoftwareLicenseKey.ValueString())
	}
	if state.ExpirationDate.ValueString() != "2026-05-20" {
		t.Fatalf("expiration_date = %q, want 2026-05-20", state.ExpirationDate.ValueString())
	}
	if state.DateActive.ValueString() != "2026-07-09" {
		t.Fatalf("date_active = %q, want 2026-07-09", state.DateActive.ValueString())
	}
}

func TestMapVaultRecordGetResponseToSoftwareLicenseModel_NSFDateStrings(t *testing.T) {
	t.Parallel()

	// NSF get responses return expirationDate / dateActive as YYYY-MM-DD strings.
	rec := &utils.VaultRecordGetResponse{
		RecordUID: "nsf-software-license-uid",
		Title:     "NSF software license",
		Type:      "softwareLicense",
		Fields: []utils.VaultRecordFieldResponse{
			{Type: "licenseNumber", Value: json.RawMessage(`["NSF-KEY-123"]`)},
			{Type: "expirationDate", Value: json.RawMessage(`["2027-12-31"]`)},
			{Type: "date", Label: "dateActive", Value: json.RawMessage(`["2026-01-15"]`)},
		},
	}

	var state commonrecordsoftwarelicense.SoftwareLicenseModel
	commonrecordsoftwarelicense.MapVaultRecordGetResponseToSoftwareLicenseModel(rec, types.StringNull(), &state)

	if state.SoftwareLicenseKey.ValueString() != "NSF-KEY-123" {
		t.Fatalf("software_license_key = %q", state.SoftwareLicenseKey.ValueString())
	}
	if state.ExpirationDate.ValueString() != "2027-12-31" {
		t.Fatalf("expiration_date = %q, want 2027-12-31", state.ExpirationDate.ValueString())
	}
	if state.DateActive.ValueString() != "2026-01-15" {
		t.Fatalf("date_active = %q, want 2026-01-15", state.DateActive.ValueString())
	}
}
