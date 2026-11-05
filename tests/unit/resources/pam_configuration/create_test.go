// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// newAwsPlanRaw builds a full config/plan tftypes.Value for environment=aws, with aws_id
// set (required), and - if configureUss is true - an aws.universal_secrets_sync block with
// enabled set to ussEnabled.
func newAwsPlanRaw(t *testing.T, sch rschema.Schema, configureUss, ussEnabled bool) tftypes.Value {
	t.Helper()
	awsType := nestedObjectType(t, sch, "aws")

	awsOverrides := map[string]tftypes.Value{
		"aws_id": stringVal("aws-1"),
	}
	if configureUss {
		ussType := ussObjectType(t, awsType)
		awsOverrides["universal_secrets_sync"] = nullObject(ussType, map[string]tftypes.Value{
			"enabled": boolVal(ussEnabled),
		})
	}

	return newValidateConfigRaw(t, sch, map[string]tftypes.Value{
		"id":                 stringVal(""),
		"environment":        stringVal("aws"),
		"title":              stringVal("Test"),
		"gateway":            stringVal("gw-uid"),
		"application_folder": stringVal("folder-uid"),
		"aws":                nullObject(awsType, awsOverrides),
	})
}

func TestCreate_SkipsUssConfigWhenNotConfigured(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		if strings.Contains(cmd, "pam config new") {
			return "network-uid-123", nil
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: newAwsPlanRaw(t, sch, false, false)}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: newResStateRaw(t, sch, "")}}

	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed (should not panic/error when universal_secrets_sync is absent): %v", resp.Diagnostics)
	}

	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config") {
			t.Errorf("expected no universal-sync-config command when block isn't configured, got: %s", c)
		}
	}
}

func TestCreate_SendsUssConfigWhenConfigured(t *testing.T) {
	mock := &helpers.CommandServer{}
	var captured []string
	server := helpers.StartCommandServer(mock, func(cmd string, _ int) (string, interface{}) {
		captured = append(captured, cmd)
		if strings.Contains(cmd, "pam config new") {
			return "network-uid-123", nil
		}
		return "ok", nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: newAwsPlanRaw(t, sch, true, true)}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: newResStateRaw(t, sch, "")}}

	r.Create(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics)
	}

	var ussCmd string
	for _, c := range captured {
		if strings.Contains(c, "universal-sync-config add") {
			ussCmd = c
		}
	}
	if ussCmd == "" {
		t.Fatalf("expected a universal-sync-config add command, got: %v", captured)
	}
	if !strings.Contains(ussCmd, "--network 'network-uid-123'") {
		t.Errorf("expected --network flag with the created network uid, got: %s", ussCmd)
	}
	if !strings.Contains(ussCmd, "--enabled true") {
		t.Errorf("expected --enabled true, got: %s", ussCmd)
	}
}

func TestCreate_UssFailureStillSavesPamConfigState(t *testing.T) {
	mock := &helpers.CommandServer{}
	server := helpers.StartCommandServerWithResultHook(mock, func(cmd string, _ int) (string, interface{}) {
		if strings.Contains(cmd, "pam config new") {
			return "network-uid-123", nil
		}
		return "ok", nil
	}, func(cmd string, _ int) (int, []byte) {
		if strings.Contains(cmd, "universal-sync-config add") {
			return http.StatusInternalServerError, []byte(`{"error":"uss add failed"}`)
		}
		return 0, nil
	})
	defer server.Close()

	r := newConfiguredRes(t, server)
	sch := getResSchema(t)

	req := resource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: newAwsPlanRaw(t, sch, true, true)}}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: newResStateRaw(t, sch, "")}}

	r.Create(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when universal-sync-config add fails")
	}

	var data commonpamconfiguration.PamConfigurationResourceModel
	if diags := resp.State.Get(context.Background(), &data); diags.HasError() {
		t.Fatalf("reading state: %v", diags)
	}
	if data.Id.ValueString() != "network-uid-123" {
		t.Errorf("expected state id network-uid-123 so the created PAM configuration isn't orphaned, got %q", data.Id.ValueString())
	}
	if data.Aws == nil {
		t.Fatal("expected aws block in state")
	}
	if data.Aws.UniversalSecretsSync != nil {
		t.Error("expected universal_secrets_sync to be omitted from state since it wasn't created")
	}
}
