// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"testing"

	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestPamConfigurationDataSource_Schema_Attributes(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	expectedAttrs := []string{
		"pam_configuration", "id", "environment", "title", "gateway", "application_folder",
		"schedule", "port_mapping",
		"connections", "tunneling", "rotation",
		"remote_browser_isolation", "connections_recording", "typescript_recording",
		"ai_threat_detection", "ai_terminate_session_on_detection",
		"local_network", "aws", "azure", "domain", "gcp",
	}
	for _, attr := range expectedAttrs {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected attribute %q in schema", attr)
		}
	}
}

func TestPamConfigurationDataSource_Schema_PamConfigurationRequired(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	attr, ok := resp.Schema.Attributes["pam_configuration"].(dschema.StringAttribute)
	if !ok {
		t.Fatal("pam_configuration is not StringAttribute")
	}
	if !attr.Required {
		t.Error("pam_configuration should be Required")
	}
}

func TestPamConfigurationDataSource_Schema_ComputedAttributes(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	computedAttrs := []string{
		"id", "environment", "title", "gateway", "application_folder",
		"connections", "tunneling", "rotation",
	}
	for _, name := range computedAttrs {
		attr, ok := resp.Schema.Attributes[name].(dschema.StringAttribute)
		if !ok {
			boolAttr, ok2 := resp.Schema.Attributes[name].(dschema.BoolAttribute)
			if !ok2 {
				t.Errorf("attribute %q is not StringAttribute or BoolAttribute", name)
				continue
			}
			if !boolAttr.Computed {
				t.Errorf("expected %q to be Computed", name)
			}
			continue
		}
		if !attr.Computed {
			t.Errorf("expected %q to be Computed", name)
		}
	}
}

func TestPamConfigurationDataSource_Schema_LocalNetworkNested(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	nested, ok := resp.Schema.Attributes["local_network"].(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("local_network is not SingleNestedAttribute")
	}
	for _, attr := range []string{"network_id", "network_cidr"} {
		if nested.Attributes[attr] == nil {
			t.Errorf("expected local_network attribute %q", attr)
		}
	}
}

func TestPamConfigurationDataSource_Schema_GcpNested(t *testing.T) {
	d := pamconfiguration.NewPamConfigurationDataSource().(*pamconfiguration.PamConfigurationDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	nested, ok := resp.Schema.Attributes["gcp"].(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("gcp is not SingleNestedAttribute")
	}
	for _, attr := range []string{"gcp_id", "service_account_key", "google_admin_email", "gcp_region"} {
		if nested.Attributes[attr] == nil {
			t.Errorf("expected gcp attribute %q", attr)
		}
	}
}
