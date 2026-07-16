// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"testing"

	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestPamConfigurationResource_Schema_Attributes(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	expectedAttrs := []string{
		"id", "environment", "title", "gateway", "application_folder",
		"schedule", "port_mapping",
		"connections", "tunneling", "rotation",
		"remote_browser_isolation", "connections_recording", "typescript_recording",
		"ai_threat_detection", "ai_terminate_session_on_detection",
	}
	for _, attr := range expectedAttrs {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected attribute %q in schema", attr)
		}
	}
}

func TestPamConfigurationResource_Schema_EnvironmentBlocks(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	expectedBlocks := []string{"local_network", "aws", "azure", "domain", "gcp"}
	for _, block := range expectedBlocks {
		if resp.Schema.Blocks[block] == nil {
			t.Errorf("expected block %q in schema", block)
		}
	}
}

func TestPamConfigurationResource_Schema_LocalNetworkBlock(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	block, ok := resp.Schema.Blocks["local_network"].(rschema.SingleNestedBlock)
	if !ok {
		t.Fatal("local_network is not SingleNestedBlock")
	}
	for _, attr := range []string{"network_id", "network_cidr"} {
		if block.Attributes[attr] == nil {
			t.Errorf("expected local_network attribute %q", attr)
		}
	}
}

func TestPamConfigurationResource_Schema_AwsBlock(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	block, ok := resp.Schema.Blocks["aws"].(rschema.SingleNestedBlock)
	if !ok {
		t.Fatal("aws is not SingleNestedBlock")
	}
	for _, attr := range []string{"aws_id", "access_key_id", "access_secret_key", "region_names"} {
		if block.Attributes[attr] == nil {
			t.Errorf("expected aws attribute %q", attr)
		}
	}
}

func TestPamConfigurationResource_Schema_AzureBlock(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	block, ok := resp.Schema.Blocks["azure"].(rschema.SingleNestedBlock)
	if !ok {
		t.Fatal("azure is not SingleNestedBlock")
	}
	for _, attr := range []string{"azure_id", "client_id", "client_secret", "subscription_id", "tenant_id", "resource_groups"} {
		if block.Attributes[attr] == nil {
			t.Errorf("expected azure attribute %q", attr)
		}
	}
}

func TestPamConfigurationResource_Schema_DomainBlock(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	block, ok := resp.Schema.Blocks["domain"].(rschema.SingleNestedBlock)
	if !ok {
		t.Fatal("domain is not SingleNestedBlock")
	}
	for _, attr := range []string{"domain_id", "domain_hostname", "domain_port", "domain_use_ssl", "domain_scan_dc_cidr", "domain_network_cidr", "domain_admin"} {
		if block.Attributes[attr] == nil {
			t.Errorf("expected domain attribute %q", attr)
		}
	}
}

func TestPamConfigurationResource_Schema_GcpBlock(t *testing.T) {
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	block, ok := resp.Schema.Blocks["gcp"].(rschema.SingleNestedBlock)
	if !ok {
		t.Fatal("gcp is not SingleNestedBlock")
	}
	for _, attr := range []string{"gcp_id", "service_account_key", "google_admin_email", "gcp_region"} {
		if block.Attributes[attr] == nil {
			t.Errorf("expected gcp attribute %q", attr)
		}
	}
}
