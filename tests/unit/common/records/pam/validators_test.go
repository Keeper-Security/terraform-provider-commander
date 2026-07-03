// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords_test

import (
	"context"
	"strings"
	"testing"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func TestValidateAllowSupplyHostHostnameConstraints_ConflictWhenBothSet(t *testing.T) {
	diags := commonpamrecords.ValidateAllowSupplyHostHostnameConstraints(
		types.BoolValue(true),
		&commonpamrecords.HostnameOrIPModel{HostName: types.StringValue("db.example.com")},
	)
	if !diags.HasError() {
		t.Fatal("expected error when allow_supply_host is true and hostname_or_ip is set")
	}
}

func TestValidateAllowSupplyHostHostnameConstraints_RequiresHostnameWhenSupplyHostFalse(t *testing.T) {
	diags := commonpamrecords.ValidateAllowSupplyHostHostnameConstraints(types.BoolValue(false), nil)
	if !diags.HasError() {
		t.Fatal("expected error when allow_supply_host is false and hostname is omitted")
	}
}

func TestValidateAllowSupplyHostHostnameConstraints_RequiresHostnameWhenSupplyHostUnset(t *testing.T) {
	diags := commonpamrecords.ValidateAllowSupplyHostHostnameConstraints(types.BoolNull(), nil)
	if !diags.HasError() {
		t.Fatal("expected error when allow_supply_host is unset and hostname is omitted")
	}
}

func TestValidateAllowSupplyHostHostnameConstraints_RequiresHostnameWhenOnlyPortSet(t *testing.T) {
	diags := commonpamrecords.ValidateAllowSupplyHostHostnameConstraints(
		types.BoolValue(false),
		&commonpamrecords.HostnameOrIPModel{AdministrativePort: types.Int32Value(22)},
	)
	if !diags.HasError() {
		t.Fatal("expected error when hostname is missing but administrative_port is set")
	}
}

func TestValidateAllowSupplyHostHostnameConstraints_AllowsSupplyHostWithoutHostname(t *testing.T) {
	diags := commonpamrecords.ValidateAllowSupplyHostHostnameConstraints(types.BoolValue(true), nil)
	if diags.HasError() {
		t.Fatalf("expected no error, got %v", diags)
	}
}

func TestValidateAllowSupplyHostHostnameConstraints_AllowsHostnameWhenSupplyHostFalse(t *testing.T) {
	diags := commonpamrecords.ValidateAllowSupplyHostHostnameConstraints(
		types.BoolValue(false),
		&commonpamrecords.HostnameOrIPModel{HostName: types.StringValue("db.example.com")},
	)
	if diags.HasError() {
		t.Fatalf("expected no error, got %v", diags)
	}
}

func TestIsHostnameOrIPSet(t *testing.T) {
	if commonpamrecords.IsHostnameOrIPSet(nil) {
		t.Fatal("nil model should not be set")
	}
	if commonpamrecords.IsHostnameOrIPSet(&commonpamrecords.HostnameOrIPModel{}) {
		t.Fatal("empty model should not be set")
	}
	if !commonpamrecords.IsHostnameOrIPSet(&commonpamrecords.HostnameOrIPModel{HostName: types.StringValue("host")}) {
		t.Fatal("hostname should count as set")
	}
	if !commonpamrecords.IsHostnameOrIPSet(&commonpamrecords.HostnameOrIPModel{AdministrativePort: types.Int32Value(22)}) {
		t.Fatal("administrative_port alone should count as set")
	}
}

func TestPamSettingsRequiredFieldsValidator_AllowsUnknownConfiguration(t *testing.T) {
	v := commonpamrecords.PamSettingsRequiredFieldsValidator()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"configuration": types.StringType,
	}, map[string]attr.Value{
		"configuration": types.StringUnknown(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("pam_settings"),
		ConfigValue: obj,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unknown configuration should be allowed during plan, got %v", resp.Diagnostics)
	}
}

func TestPamSettingsRequiredFieldsValidator_ErrorsOnNullConfiguration(t *testing.T) {
	v := commonpamrecords.PamSettingsRequiredFieldsValidator()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"configuration": types.StringType,
	}, map[string]attr.Value{
		"configuration": types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("pam_settings"),
		ConfigValue: obj,
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when configuration is null")
	}
}

func TestConnectionFieldsRequireEnabledValidator_AllowsUnknownLaunchCredential(t *testing.T) {
	v := commonpamrecords.ConnectionFieldsRequireEnabledValidator()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"enable":            types.BoolType,
		"connection_port":   types.Int32Type,
		"launch_credential": types.StringType,
	}, map[string]attr.Value{
		"enable":            types.BoolValue(true),
		"connection_port":   types.Int32Value(22),
		"launch_credential": types.StringUnknown(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("pam_settings").AtName("connection"),
		ConfigValue: obj,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unknown launch_credential should be allowed during plan, got %v", resp.Diagnostics)
	}
}

func TestConnectionFieldsRequireEnabledValidator_ErrorsOnNullLaunchCredential(t *testing.T) {
	v := commonpamrecords.ConnectionFieldsRequireEnabledValidator()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"enable":            types.BoolType,
		"connection_port":   types.Int32Type,
		"launch_credential": types.StringType,
	}, map[string]attr.Value{
		"enable":            types.BoolValue(true),
		"connection_port":   types.Int32Value(22),
		"launch_credential": types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("pam_settings").AtName("connection"),
		ConfigValue: obj,
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when launch_credential is null and connection is enabled")
	}
}

func TestSftpUserUidRequiredValidator_AllowsUnknownUserUID(t *testing.T) {
	v := commonpamrecords.SftpUserUidRequiredValidator()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"enable_sftp":   types.BoolType,
		"sftp_user_uid": types.StringType,
	}, map[string]attr.Value{
		"enable_sftp":   types.BoolValue(true),
		"sftp_user_uid": types.StringUnknown(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("pam_settings").AtName("connection").AtName("rdp").AtName("sftp"),
		ConfigValue: obj,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unknown sftp_user_uid should be allowed during plan, got %v", resp.Diagnostics)
	}
}

func TestTunnelLocalPortRequiredValidator_AllowsUnknownLocalPort(t *testing.T) {
	v := commonpamrecords.TunnelLocalPortRequiredValidator()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"use_specified_local_port": types.BoolType,
		"local_port":               types.Int32Type,
	}, map[string]attr.Value{
		"use_specified_local_port": types.BoolValue(true),
		"local_port":               types.Int32Unknown(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("pam_settings").AtName("tunnel"),
		ConfigValue: obj,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unknown local_port should be allowed during plan, got %v", resp.Diagnostics)
	}
}
