// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords_test

import (
	"context"
	"strings"
	"testing"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func runProtocolValidator(t *testing.T, allowed []string, val types.String) validator.StringResponse {
	t.Helper()
	v := commonpamrecords.ConnectionProtocolValidator(allowed)
	req := validator.StringRequest{ConfigValue: val, Path: path.Root("pam_settings").AtName("connection").AtName("protocol")}
	var resp validator.StringResponse
	v.ValidateString(context.Background(), req, &resp)
	return resp
}

func TestConnectionProtocolValidator_AllowsMachineDirectoryProtocols(t *testing.T) {
	for _, p := range commonpamrecords.MachineDirectoryProtocols {
		resp := runProtocolValidator(t, commonpamrecords.MachineDirectoryProtocols, types.StringValue(p))
		if resp.Diagnostics.HasError() {
			t.Errorf("protocol %q should be allowed for machine/directory, got %v", p, resp.Diagnostics)
		}
	}
}

func TestConnectionProtocolValidator_RejectsDatabaseProtocolsForMachineDirectory(t *testing.T) {
	for _, p := range commonpamrecords.DatabaseProtocols {
		resp := runProtocolValidator(t, commonpamrecords.MachineDirectoryProtocols, types.StringValue(p))
		if !resp.Diagnostics.HasError() {
			t.Errorf("protocol %q should be rejected on machine/directory records", p)
		}
	}
}

func TestConnectionProtocolValidator_AllowsDatabaseProtocols(t *testing.T) {
	for _, p := range commonpamrecords.DatabaseProtocols {
		resp := runProtocolValidator(t, commonpamrecords.DatabaseProtocols, types.StringValue(p))
		if resp.Diagnostics.HasError() {
			t.Errorf("protocol %q should be allowed for database, got %v", p, resp.Diagnostics)
		}
	}
}

func TestConnectionProtocolValidator_RejectsMachineDirectoryProtocolsForDatabase(t *testing.T) {
	for _, p := range commonpamrecords.MachineDirectoryProtocols {
		resp := runProtocolValidator(t, commonpamrecords.DatabaseProtocols, types.StringValue(p))
		if !resp.Diagnostics.HasError() {
			t.Errorf("protocol %q should be rejected on database records", p)
		}
	}
}

func TestConnectionProtocolValidator_RejectsUnknownProtocol(t *testing.T) {
	resp := runProtocolValidator(t, commonpamrecords.MachineDirectoryProtocols, types.StringValue("ftp"))
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostic for unsupported protocol")
	}
}

func TestConnectionProtocolValidator_NullOrUnknownNoDiag(t *testing.T) {
	for _, val := range []types.String{types.StringNull(), types.StringUnknown()} {
		resp := runProtocolValidator(t, commonpamrecords.MachineDirectoryProtocols, val)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for null/unknown protocol, got %v", resp.Diagnostics)
		}
	}
}

func TestConnectionProtocolValidator_DescriptionIncludesAllowedList(t *testing.T) {
	v := commonpamrecords.ConnectionProtocolValidator(commonpamrecords.DatabaseProtocols)
	got := v.Description(context.Background())
	for _, p := range commonpamrecords.DatabaseProtocols {
		if !strings.Contains(got, p) {
			t.Errorf("Description %q must mention allowed protocol %q", got, p)
		}
	}
	if strings.Contains(got, commonpamrecords.ConnectionProtocolSsh) {
		t.Errorf("Database Description must not mention machine/directory-only protocol %q: %s", commonpamrecords.ConnectionProtocolSsh, got)
	}
	if v.MarkdownDescription(context.Background()) != got {
		t.Error("MarkdownDescription should equal Description")
	}
}
