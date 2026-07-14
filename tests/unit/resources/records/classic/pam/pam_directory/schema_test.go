// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"testing"

	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestPamDirectoryResource_Schema(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)

	expectedAttrs := []string{
		"id", "title", "hostname_or_ip",
		"use_ssl", "domain_name", "alternative_ips",
		"directory_id", "directory_type", "user_match",
		"provider_group", "provider_region",
		"notes", "folder_location",
	}
	for _, attr := range expectedAttrs {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected %s attribute", attr)
		}
	}
}

func TestPamDirectoryResource_SchemaHostnameOrIPBlock(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
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

func TestPamDirectoryResource_SchemaPamSettingsBlock(t *testing.T) {
	r := pamdirectory.NewPamDirectoryResource().(*pamdirectory.PamDirectoryResource)
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
}
