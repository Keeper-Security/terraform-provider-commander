// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"context"
	"testing"

	commander "github.com/Keeper-Security/terraform-provider-commander/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func TestCommanderProvider_Metadata(t *testing.T) {
	p := commander.New("test")()
	req := provider.MetadataRequest{}
	var resp provider.MetadataResponse
	p.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "commander" {
		t.Errorf("expected TypeName commander, got %s", resp.TypeName)
	}
	if resp.Version != "test" {
		t.Errorf("expected Version test, got %s", resp.Version)
	}
}

func TestCommanderProvider_Schema(t *testing.T) {
	p := commander.New("test")()
	req := provider.SchemaRequest{}
	var resp provider.SchemaResponse
	p.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["service_mode_url"] == nil {
		t.Error("expected service_mode_url attribute in schema")
	}
	if resp.Schema.Attributes["service_mode_api_key"] == nil {
		t.Error("expected service_mode_api_key attribute in schema")
	}
	urlAttr, ok := resp.Schema.Attributes["service_mode_url"].(schema.StringAttribute)
	if !ok || !urlAttr.Required {
		t.Error("service_mode_url should be required StringAttribute")
	}
}

func TestCommanderProvider_Resources(t *testing.T) {
	p := commander.New("test")()
	resources := p.Resources(context.Background())
	if len(resources) != 5 {
		t.Errorf("expected 5 resources, got %d", len(resources))
	}
	for i, factory := range resources {
		r := factory()
		if r == nil {
			t.Errorf("resource factory %d returned nil", i)
		}
	}
}

func TestCommanderProvider_DataSources(t *testing.T) {
	p := commander.New("test")()
	dataSources := p.DataSources(context.Background())
	if len(dataSources) != 5 {
		t.Errorf("expected 5 data sources, got %d", len(dataSources))
	}
	for i, factory := range dataSources {
		ds := factory()
		if ds == nil {
			t.Errorf("data source factory %d returned nil", i)
		}
	}
}

func TestCommanderProvider_EphemeralResources(t *testing.T) {
	p := commander.New("test")()
	cp, ok := p.(*commander.CommanderProvider)
	if !ok {
		t.Fatalf("expected *CommanderProvider, got %T", p)
	}
	ephemeral := cp.EphemeralResources(context.Background())
	if len(ephemeral) != 0 {
		t.Errorf("expected 0 ephemeral resources, got %d", len(ephemeral))
	}
}

func TestCommanderProvider_Functions(t *testing.T) {
	p := commander.New("test")()
	cp, ok := p.(*commander.CommanderProvider)
	if !ok {
		t.Fatalf("expected *CommanderProvider, got %T", p)
	}
	functions := cp.Functions(context.Background())
	if len(functions) != 0 {
		t.Errorf("expected 0 functions, got %d", len(functions))
	}
}

func TestCommanderProvider_Actions(t *testing.T) {
	p := commander.New("test")()
	cp, ok := p.(*commander.CommanderProvider)
	if !ok {
		t.Fatalf("expected *CommanderProvider, got %T", p)
	}
	actions := cp.Actions(context.Background())
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestNew(t *testing.T) {
	factory := commander.New("1.0.0")
	if factory == nil {
		t.Fatal("New returned nil")
	}
	p := factory()
	if p == nil {
		t.Fatal("factory returned nil provider")
	}
	var resp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)
	if resp.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", resp.Version)
	}
}
