// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"testing"

	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic_records/pam_records/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestPamRemoteBrowserDataSource_Schema(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), req, &resp)

	for _, attr := range []string{"remote_browser", "id", "title", "url", "notes", "folder_location", "pam_remote_browser_settings"} {
		if resp.Schema.Attributes[attr] == nil {
			t.Errorf("expected %s attribute", attr)
		}
	}
}

func TestPamRemoteBrowserDataSource_SchemaSettingsBlock(t *testing.T) {
	d := pamremotebrowser.NewPamRemoteBrowserDataSource().(*pamremotebrowser.PamRemoteBrowserDataSource)
	var resp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	settingsAttr, ok := resp.Schema.Attributes["pam_remote_browser_settings"]
	if !ok {
		t.Fatal("pam_remote_browser_settings attribute missing")
	}

	nested, ok := settingsAttr.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected SingleNestedAttribute, got %T", settingsAttr)
	}

	expectedAttrs := []string{
		"configuration", "remote_browser_isolation", "connections_recording",
		"key_events", "allow_url_navigation", "ignore_server_cert",
		"allowed_urls", "allowed_resource_urls", "auto_fill_targets",
		"auto_fill_credentials", "allow_copy", "allow_paste",
		"disable_audio", "audio_channels", "audio_bit_depth", "audio_sample_rate",
	}
	for _, attr := range expectedAttrs {
		if nested.Attributes[attr] == nil {
			t.Errorf("expected nested attribute %s", attr)
		}
	}
}
