// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// newAwsStateRaw builds a full state tftypes.Value for environment=aws with the given id,
// mirroring newAwsPlanRaw but with an explicit id (required for state).
func newAwsStateRaw(t *testing.T, sch rschema.Schema, id string, configureUss, ussEnabled bool) tftypes.Value {
	t.Helper()
	awsType := nestedObjectType(t, sch, "aws")

	awsOverrides := map[string]tftypes.Value{
		"aws_id": stringVal("aws-1"),
	}
	if configureUss {
		ussType := childObjectType(t, awsType, "universal_secrets_sync")
		awsOverrides["universal_secrets_sync"] = nullObject(ussType, map[string]tftypes.Value{
			"enabled": boolVal(ussEnabled),
		})
	}

	return newValidateConfigRaw(t, sch, map[string]tftypes.Value{
		"id":                 stringVal(id),
		"environment":        stringVal("aws"),
		"title":              stringVal("Test"),
		"gateway":            stringVal("gw-uid"),
		"application_folder": stringVal("folder-uid"),
		"aws":                nullObject(awsType, awsOverrides),
	})
}

func TestUpdate_SendsUssEditWhenConfiguredInPlan(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: newAwsStateRaw(t, sch, "", true, true)},
		State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", true, false)},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", true, false)}}

	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}

	var ussCmd string
	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config edit") {
			ussCmd = c
		}
		if strings.Contains(c, "universal-sync-config remove") {
			t.Fatalf("expected no remove command when plan still has universal_secrets_sync configured, got: %s", c)
		}
	}
	if ussCmd == "" {
		t.Fatalf("expected a universal-sync-config edit command, got: %v", captured)
	}
	if !strings.Contains(ussCmd, "--network 'network-uid-123'") {
		t.Errorf("expected --network flag with the state's network uid, got: %s", ussCmd)
	}
	if !strings.Contains(ussCmd, "--enabled true") {
		t.Errorf("expected --enabled true, got: %s", ussCmd)
	}
}

func TestUpdate_CreatesUssConfigWhenConfiguredForFirstTime(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: newAwsStateRaw(t, sch, "", true, true)},
		State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", false, false)},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", false, false)}}

	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}

	var newCmd string
	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config edit") || strings.Contains(c, "universal-sync-config remove") {
			t.Fatalf("expected no edit/remove command when universal_secrets_sync is being configured for the first time, got: %s", c)
		}
		if strings.Contains(c, "universal-sync-config add") {
			newCmd = c
		}
	}
	if newCmd == "" {
		t.Fatalf("expected a universal-sync-config add command, got: %v", captured)
	}
	if !strings.Contains(newCmd, "--network 'network-uid-123'") {
		t.Errorf("expected --network flag with the state's network uid, got: %s", newCmd)
	}
	if !strings.Contains(newCmd, "--enabled true") {
		t.Errorf("expected --enabled true, got: %s", newCmd)
	}
}

func TestUpdate_RemovesUssConfigWhenRemovedFromPlan(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: newAwsStateRaw(t, sch, "", false, false)},
		State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", true, true)},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", true, true)}}

	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}

	var removeCmd string
	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config edit") || strings.Contains(c, "universal-sync-config add") {
			t.Fatalf("expected no edit/new command when universal_secrets_sync was removed from the plan, got: %s", c)
		}
		if strings.Contains(c, "universal-sync-config remove") {
			removeCmd = c
		}
	}
	if removeCmd == "" {
		t.Fatalf("expected a universal-sync-config remove command, got: %v", captured)
	}
	if !strings.Contains(removeCmd, "--network 'network-uid-123'") {
		t.Errorf("expected --network flag with the state's network uid, got: %s", removeCmd)
	}
	// pam universal-sync-config remove doesn't support --format json, unlike list.
	if strings.Contains(removeCmd, "--format json") {
		t.Errorf("expected no --format json (remove doesn't support it), got: %s", removeCmd)
	}
}

func TestUpdate_SkipsUssWhenNeverConfigured(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: newAwsStateRaw(t, sch, "", false, false)},
		State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", false, false)},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", false, false)}}

	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}

	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config") {
			t.Errorf("expected no universal-sync-config command when it was never configured, got: %s", c)
		}
	}
}

func TestUpdate_SkipsUssEditWhenUnchanged(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: newAwsStateRaw(t, sch, "", true, true)},
		State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", true, true)},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: newAwsStateRaw(t, sch, "network-uid-123", true, true)}}

	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}

	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config") {
			t.Errorf("expected no universal-sync-config command when the block is unchanged, got: %s", c)
		}
	}
}
