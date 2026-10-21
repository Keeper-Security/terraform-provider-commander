// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"testing"

	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestModifyPlan_UssFirstTimeEnabledFalse(t *testing.T) {
	sch := getResSchema(t)
	nullRaw := tftypes.NewValue(sch.Type().TerraformType(context.Background()), nil)

	awsType := nestedObjectType(t, sch, "aws")
	ussType := childObjectType(t, awsType, "universal_secrets_sync")
	enabledOmittedPlan := newValidateConfigRaw(t, sch, map[string]tftypes.Value{
		"id":                 stringVal(""),
		"environment":        stringVal("aws"),
		"title":              stringVal("Test"),
		"gateway":            stringVal("gw-uid"),
		"application_folder": stringVal("folder-uid"),
		"aws": nullObject(awsType, map[string]tftypes.Value{
			"aws_id":                 stringVal("aws-1"),
			"universal_secrets_sync": nullObject(ussType, map[string]tftypes.Value{"dry_run": boolVal(true)}),
		}),
	})

	tests := []struct {
		name      string
		plan      tftypes.Value
		state     tftypes.Value
		wantError bool
	}{
		{"create with enabled false", newAwsPlanRaw(t, sch, true, false), nullRaw, true},
		{"create with enabled omitted", enabledOmittedPlan, nullRaw, true},
		{"create with enabled true", newAwsPlanRaw(t, sch, true, true), nullRaw, false},
		{"create without uss", newAwsPlanRaw(t, sch, false, false), nullRaw, false},
		{"update adds uss with enabled false", newAwsStateRaw(t, sch, "net-1", true, false), newAwsStateRaw(t, sch, "net-1", false, false), true},
		{"update adds uss with enabled true", newAwsStateRaw(t, sch, "net-1", true, true), newAwsStateRaw(t, sch, "net-1", false, false), false},
		{"update disables existing uss", newAwsStateRaw(t, sch, "net-1", true, false), newAwsStateRaw(t, sch, "net-1", true, true), false},
		{"update keeps existing disabled uss", newAwsStateRaw(t, sch, "net-1", true, false), newAwsStateRaw(t, sch, "net-1", true, false), false},
		{"destroy", nullRaw, newAwsStateRaw(t, sch, "net-1", true, false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
			req := resource.ModifyPlanRequest{
				Plan:  tfsdk.Plan{Schema: sch, Raw: tt.plan},
				State: tfsdk.State{Schema: sch, Raw: tt.state},
			}
			resp := resource.ModifyPlanResponse{Plan: tfsdk.Plan{Schema: sch, Raw: tt.plan}}

			r.ModifyPlan(context.Background(), req, &resp)

			if got := resp.Diagnostics.HasError(); got != tt.wantError {
				t.Errorf("HasError = %v, want %v; diagnostics: %v", got, tt.wantError, resp.Diagnostics)
			}
		})
	}
}
