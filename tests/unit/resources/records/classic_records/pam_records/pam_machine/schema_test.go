// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine_test

import (
	"context"
	"testing"

	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestPamMachineResource_Schema(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)

	expectedAttrs := []string{
		"id", "title", "hostname_or_ip",
		"operating_system", "instance_name", "instance_id",
		"provider_group", "provider_region",
		"notes", "folder_location",
	}
	for _, attr := range expectedAttrs {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected %s attribute", attr)
		}
	}
}

func TestPamMachineResource_SchemaHostnameOrIPBlock(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	hostAttr, ok := resp.Schema.Attributes["hostname_or_ip"]
	if !ok {
		t.Fatal("hostname_or_ip attribute missing")
	}

	nested, ok := hostAttr.(rschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected SingleNestedAttribute, got %T", hostAttr)
	}

	for _, attr := range []string{"hostname", "administrative_port"} {
		if nested.Attributes[attr] == nil {
			t.Errorf("expected nested attribute %s", attr)
		}
	}
}

func TestPamMachineResource_SchemaPamSettingsBlock(t *testing.T) {
	r := pammachine.NewPamMachineResource().(*pammachine.PamMachineResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	pamBlock, ok := resp.Schema.Blocks["pam_settings"]
	if !ok {
		t.Fatal("pam_settings block missing")
	}

	snb, ok := pamBlock.(rschema.SingleNestedBlock)
	if !ok {
		t.Fatalf("expected SingleNestedBlock, got %T", pamBlock)
	}

	for _, attr := range []string{"allow_supply_host", "configuration", "administrative_credentials"} {
		if snb.Attributes[attr] == nil {
			t.Errorf("expected pam_settings attribute %s", attr)
		}
	}

	for _, block := range []string{"connection", "tunnel"} {
		if snb.Blocks[block] == nil {
			t.Errorf("expected pam_settings block %s", block)
		}
	}

	connBlock, ok := snb.Blocks["connection"].(rschema.SingleNestedBlock)
	if !ok {
		t.Fatal("expected connection SingleNestedBlock")
	}
	for _, proto := range []string{"kubernetes", "rdp", "ssh", "telnet", "vnc"} {
		if connBlock.Blocks[proto] == nil {
			t.Errorf("expected connection protocol block %s", proto)
		}
	}
	for _, proto := range []string{"mysql", "postgresql", "sql_server"} {
		if connBlock.Blocks[proto] != nil {
			t.Errorf("unexpected database connection protocol block %s on pam_machine", proto)
		}
	}
}
