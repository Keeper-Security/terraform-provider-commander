// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"testing"

	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_records/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestPamMachineDataSource_Schema(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)

	expectedAttrs := []string{
		"pam_machine", "id", "title", "hostname_or_ip",
		"operating_system", "instance_name", "instance_id",
		"provider_group", "provider_region",
		"notes", "folder", "pam_settings",
	}
	for _, attr := range expectedAttrs {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected %s attribute", attr)
		}
	}
}

func TestPamMachineDataSource_SchemaHostnameOrIPNested(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	hostAttr, ok := resp.Schema.Attributes["hostname_or_ip"]
	if !ok {
		t.Fatal("hostname_or_ip attribute missing")
	}

	nested, ok := hostAttr.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected SingleNestedAttribute, got %T", hostAttr)
	}

	for _, attr := range []string{"hostname", "administrative_port"} {
		if nested.Attributes[attr] == nil {
			t.Errorf("expected nested attribute %s", attr)
		}
	}
}

func TestPamMachineDataSource_SchemaPamSettingsNested(t *testing.T) {
	d := pammachine.NewPamMachineDataSource().(*pammachine.PamMachineDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	pamAttr, ok := resp.Schema.Attributes["pam_settings"]
	if !ok {
		t.Fatal("pam_settings attribute missing")
	}

	nested, ok := pamAttr.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected SingleNestedAttribute, got %T", pamAttr)
	}

	for _, attr := range []string{"allow_supply_host", "configuration", "administrative_credentials"} {
		if nested.Attributes[attr] == nil {
			t.Errorf("expected pam_settings attribute %s", attr)
		}
	}

	connAttr, ok := nested.Attributes["connection"]
	if !ok {
		t.Fatal("connection attribute missing in pam_settings")
	}
	connNested, ok := connAttr.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected SingleNestedAttribute for connection, got %T", connAttr)
	}
	for _, proto := range []string{"kubernetes", "mysql", "postgresql", "rdp", "sql_server", "ssh", "telnet", "vnc"} {
		if connNested.Attributes[proto] == nil {
			t.Errorf("expected connection attribute %s", proto)
		}
	}

	tunnelAttr, ok := nested.Attributes["tunnel"]
	if !ok {
		t.Fatal("tunnel attribute missing in pam_settings")
	}
	tunnelNested, ok := tunnelAttr.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected SingleNestedAttribute for tunnel, got %T", tunnelAttr)
	}
	for _, attr := range []string{"enable", "remote_target_port", "re_use_port", "use_specified_local_port", "local_port"} {
		if tunnelNested.Attributes[attr] == nil {
			t.Errorf("expected tunnel attribute %s", attr)
		}
	}
}
