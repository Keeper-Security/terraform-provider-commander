// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser_test

import (
	"context"
	"testing"

	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var rotationSettingsAttrTypes = map[string]attr.Type{
	"rotation_profile": types.StringType,
	"configuration":    types.StringType,
	"iam_aad_config":   types.StringType,
	"saas_config":      types.StringType,
	"resource":         types.StringType,
	"on_demand":        types.BoolType,
	"schedule_config":  types.BoolType,
	"schedule_cron":    types.StringType,
	"schedule_json":    types.StringType,
}

func rotationSettingsAttrs(overrides map[string]attr.Value) map[string]attr.Value {
	attrs := map[string]attr.Value{
		"rotation_profile": types.StringNull(),
		"configuration":    types.StringNull(),
		"iam_aad_config":   types.StringNull(),
		"saas_config":      types.StringNull(),
		"resource":         types.StringNull(),
		"on_demand":        types.BoolNull(),
		"schedule_config":  types.BoolNull(),
		"schedule_cron":    types.StringNull(),
		"schedule_json":    types.StringNull(),
	}
	for k, v := range overrides {
		attrs[k] = v
	}
	return attrs
}

func generalRotationSettings(overrides map[string]attr.Value) map[string]attr.Value {
	base := map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileGeneral),
		"resource":         types.StringValue("pam-resource-uid"),
	}
	for k, v := range overrides {
		base[k] = v
	}
	return rotationSettingsAttrs(base)
}

func validateRotationSettings(t *testing.T, attrs map[string]attr.Value, validators ...validator.Object) validator.ObjectResponse {
	t.Helper()
	if len(validators) == 0 {
		validators = []validator.Object{
			commonpamuser.RotationProfileRequirementsValidator(),
			commonpamuser.RotationScheduleCombinationValidator(),
		}
	}

	obj, diags := types.ObjectValue(rotationSettingsAttrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	req := validator.ObjectRequest{
		Path:        path.Root("rotation_settings"),
		ConfigValue: obj,
	}
	for _, v := range validators {
		v.ValidateObject(context.Background(), req, &resp)
	}
	return resp
}

func TestRotationProfileRequirementsValidator_MissingProfile(t *testing.T) {
	t.Parallel()

	resp := validateRotationSettings(t, rotationSettingsAttrs(map[string]attr.Value{
		"on_demand": types.BoolValue(true),
	}), commonpamuser.RotationProfileRequirementsValidator())
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when rotation_profile is missing")
	}
}

func TestRotationProfileRequirementsValidator_GeneralRequiresResource(t *testing.T) {
	t.Parallel()

	resp := validateRotationSettings(t, rotationSettingsAttrs(map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileGeneral),
	}), commonpamuser.RotationProfileRequirementsValidator())
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when resource is missing for general profile")
	}
}

func TestRotationProfileRequirementsValidator_GeneralForbidsSaaSConfigAndIamAadConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		overrides map[string]attr.Value
	}{
		{
			name: "saas_config",
			overrides: map[string]attr.Value{
				"rotation_profile": types.StringValue(commonpamuser.RotProfileGeneral),
				"resource":         types.StringValue("pam-resource-uid"),
				"saas_config":      types.StringValue("saas-config-uid"),
			},
		},
		{
			name: "iam_aad_config",
			overrides: map[string]attr.Value{
				"rotation_profile": types.StringValue(commonpamuser.RotProfileGeneral),
				"resource":         types.StringValue("pam-resource-uid"),
				"iam_aad_config":   types.StringValue("iam-config-uid"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := validateRotationSettings(t, rotationSettingsAttrs(tc.overrides), commonpamuser.RotationProfileRequirementsValidator())
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected error when %s is set for general profile", tc.name)
			}
		})
	}
}

func TestRotationProfileRequirementsValidator_IAMUserRequiresIamAadConfig(t *testing.T) {
	t.Parallel()

	resp := validateRotationSettings(t, rotationSettingsAttrs(map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileIAMUser),
	}), commonpamuser.RotationProfileRequirementsValidator())
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when iam_aad_config is missing for iam_user profile")
	}
}

func TestRotationProfileRequirementsValidator_IAMUserForbidsConfigurationResourceAndSaaSConfig(t *testing.T) {
	t.Parallel()

	base := map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileIAMUser),
		"iam_aad_config":   types.StringValue("iam-config-uid"),
	}

	cases := []struct {
		name      string
		field     string
		overrides map[string]attr.Value
	}{
		{
			name:  "configuration",
			field: "configuration",
			overrides: map[string]attr.Value{
				"configuration": types.StringValue("pam-config-uid"),
			},
		},
		{
			name:  "resource",
			field: "resource",
			overrides: map[string]attr.Value{
				"resource": types.StringValue("pam-resource-uid"),
			},
		},
		{
			name:  "saas_config",
			field: "saas_config",
			overrides: map[string]attr.Value{
				"saas_config": types.StringValue("saas-config-uid"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			overrides := map[string]attr.Value{}
			for k, v := range base {
				overrides[k] = v
			}
			for k, v := range tc.overrides {
				overrides[k] = v
			}
			resp := validateRotationSettings(t, rotationSettingsAttrs(overrides), commonpamuser.RotationProfileRequirementsValidator())
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected error when %s is set for iam_user profile", tc.field)
			}
		})
	}
}

func TestRotationProfileRequirementsValidator_ScriptsOnlyRequiresConfiguration(t *testing.T) {
	t.Parallel()

	resp := validateRotationSettings(t, rotationSettingsAttrs(map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileScriptsOnly),
	}), commonpamuser.RotationProfileRequirementsValidator())
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when configuration is missing for scripts_only profile")
	}
}

func TestRotationProfileRequirementsValidator_ScriptsOnlyForbidsIamAadConfigResourceAndSaaSConfig(t *testing.T) {
	t.Parallel()

	base := map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileScriptsOnly),
		"configuration":    types.StringValue("pam-config-uid"),
	}

	cases := []struct {
		name      string
		overrides map[string]attr.Value
	}{
		{
			name: "iam_aad_config",
			overrides: map[string]attr.Value{
				"iam_aad_config": types.StringValue("iam-config-uid"),
			},
		},
		{
			name: "resource",
			overrides: map[string]attr.Value{
				"resource": types.StringValue("pam-resource-uid"),
			},
		},
		{
			name: "saas_config",
			overrides: map[string]attr.Value{
				"saas_config": types.StringValue("saas-config-uid"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			overrides := map[string]attr.Value{}
			for k, v := range base {
				overrides[k] = v
			}
			for k, v := range tc.overrides {
				overrides[k] = v
			}
			resp := validateRotationSettings(t, rotationSettingsAttrs(overrides), commonpamuser.RotationProfileRequirementsValidator())
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected error when %s is set for scripts_only profile", tc.name)
			}
		})
	}
}

func TestRotationProfileRequirementsValidator_SaaSRequiresConfigurationAndSaaSConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		overrides map[string]attr.Value
	}{
		{
			name: "missing configuration",
			overrides: map[string]attr.Value{
				"rotation_profile": types.StringValue(commonpamuser.RotProfileSaaS),
				"saas_config":      types.StringValue("saas-config-uid"),
			},
		},
		{
			name: "missing saas_config",
			overrides: map[string]attr.Value{
				"rotation_profile": types.StringValue(commonpamuser.RotProfileSaaS),
				"configuration":    types.StringValue("pam-config-uid"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := validateRotationSettings(t, rotationSettingsAttrs(tc.overrides), commonpamuser.RotationProfileRequirementsValidator())
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected error for %#v", tc.overrides)
			}
		})
	}
}

func TestRotationProfileRequirementsValidator_SaaSForbidsResourceAndIamAadConfig(t *testing.T) {
	t.Parallel()

	base := map[string]attr.Value{
		"rotation_profile": types.StringValue(commonpamuser.RotProfileSaaS),
		"configuration":    types.StringValue("pam-config-uid"),
		"saas_config":      types.StringValue("saas-config-uid"),
	}

	cases := []struct {
		name      string
		overrides map[string]attr.Value
	}{
		{
			name: "resource",
			overrides: map[string]attr.Value{
				"resource": types.StringValue("pam-resource-uid"),
			},
		},
		{
			name: "iam_aad_config",
			overrides: map[string]attr.Value{
				"iam_aad_config": types.StringValue("iam-config-uid"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			overrides := map[string]attr.Value{}
			for k, v := range base {
				overrides[k] = v
			}
			for k, v := range tc.overrides {
				overrides[k] = v
			}
			resp := validateRotationSettings(t, rotationSettingsAttrs(overrides), commonpamuser.RotationProfileRequirementsValidator())
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected error when %s is set for saas profile", tc.name)
			}
		})
	}
}

func TestRotationProfileRequirementsValidator_ValidProfiles(t *testing.T) {
	t.Parallel()

	cases := []map[string]attr.Value{
		{
			"rotation_profile": types.StringValue(commonpamuser.RotProfileGeneral),
			"resource":         types.StringValue("pam-resource-uid"),
			"configuration":    types.StringValue("pam-config-uid"),
		},
		{
			"rotation_profile": types.StringValue(commonpamuser.RotProfileIAMUser),
			"iam_aad_config":   types.StringValue("iam-config-uid"),
		},
		{
			"rotation_profile": types.StringValue(commonpamuser.RotProfileScriptsOnly),
			"configuration":    types.StringValue("pam-config-uid"),
		},
		{
			"rotation_profile": types.StringValue(commonpamuser.RotProfileSaaS),
			"configuration":    types.StringValue("pam-config-uid"),
			"saas_config":      types.StringValue("saas-config-uid"),
		},
	}

	for _, overrides := range cases {
		resp := validateRotationSettings(t, rotationSettingsAttrs(overrides), commonpamuser.RotationProfileRequirementsValidator())
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected valid for %#v, got: %v", overrides, resp.Diagnostics)
		}
	}
}

func TestRotationScheduleCombinationValidator_SingleOptionValid(t *testing.T) {
	t.Parallel()

	cases := []map[string]attr.Value{
		{"on_demand": types.BoolValue(true)},
		{"schedule_config": types.BoolValue(true)},
		{"schedule_cron": types.StringValue("0 28 17 ? * *")},
		{"schedule_json": types.StringValue(`{"type":"WEEKLY","weekday":"SATURDAY","time":"17:00:00","tz":"Etc/UTC"}`)},
	}

	for _, overrides := range cases {
		resp := validateRotationSettings(t, generalRotationSettings(overrides), commonpamuser.RotationScheduleCombinationValidator())
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected valid for %#v, got: %v", overrides, resp.Diagnostics)
		}
	}
}

func TestRotationScheduleCombinationValidator_MultipleOptionsInvalid(t *testing.T) {
	t.Parallel()

	cases := []map[string]attr.Value{
		{
			"on_demand":     types.BoolValue(true),
			"schedule_cron": types.StringValue("0 28 17 ? * *"),
		},
		{
			"schedule_cron": types.StringValue("0 28 17 ? * *"),
			"schedule_json": types.StringValue(`{"type":"DAILY","time":"17:00:00","tz":"Etc/UTC"}`),
		},
		{
			"on_demand":       types.BoolValue(true),
			"schedule_config": types.BoolValue(true),
		},
		{
			"schedule_config": types.BoolValue(true),
			"schedule_json":   types.StringValue(`{"type":"DAILY","utcTime":"17:56","tz":"Etc/UTC"}`),
		},
	}

	for _, overrides := range cases {
		resp := validateRotationSettings(t, generalRotationSettings(overrides), commonpamuser.RotationScheduleCombinationValidator())
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected error for %#v", overrides)
		}
	}
}

func TestRotationScheduleCombinationValidator_FalseOnDemandWithCronValid(t *testing.T) {
	t.Parallel()

	resp := validateRotationSettings(t, generalRotationSettings(map[string]attr.Value{
		"on_demand":     types.BoolValue(false),
		"schedule_cron": types.StringValue("0 28 17 ? * *"),
	}), commonpamuser.RotationScheduleCombinationValidator())
	if resp.Diagnostics.HasError() {
		t.Fatalf("on_demand=false with schedule_cron should be valid: %v", resp.Diagnostics)
	}
}

func TestValidateRotationScheduleJSON_ValidExamples(t *testing.T) {
	t.Parallel()

	cases := []string{
		`{"type":"DAILY","intervalCount":1,"time":"17:00:00","tz":"Asia/Calcutta"}`,
		`{"type":"WEEKLY","intervalCount":1,"time":"17:00:00","tz":"Asia/Calcutta","weekday":"WEDNESDAY"}`,
		`{"type":"WEEKLY","utcTime":"00:00","weekday":"SATURDAY","intervalCount":1,"tz":"Etc/UTC"}`,
		`{"type":"MONTHLY_BY_WEEKDAY","intervalCount":1,"time":"09:30:00","tz":"America/New_York","weekday":"TUESDAY","occurrence":"SECOND"}`,
		`{"type":"YEARLY","intervalCount":1,"time":"00:00:00","tz":"Etc/UTC","month":"MAY","monthDay":20}`,
	}

	for _, jsonVal := range cases {
		if err := commonpamuser.ValidateRotationScheduleJSON(jsonVal); err != nil {
			t.Fatalf("expected valid schedule_json %s, got: %v", jsonVal, err)
		}
	}
}

func TestValidateRotationScheduleJSON_InvalidExamples(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		json string
	}{
		{name: "empty", json: ""},
		{name: "not json", json: "not-json"},
		{name: "missing type", json: `{"time":"17:00:00"}`},
		{name: "unknown type", json: `{"type":"HOURLY","time":"17:00:00"}`},
		{name: "daily with weekday", json: `{"type":"DAILY","time":"17:00:00","weekday":"MONDAY"}`},
		{name: "weekly missing weekday", json: `{"type":"WEEKLY","time":"17:00:00"}`},
		{name: "weekly invalid weekday", json: `{"type":"WEEKLY","time":"17:00:00","weekday":"FUNDAY"}`},
		{name: "both time fields", json: `{"type":"DAILY","time":"17:00:00","utcTime":"17:00"}`},
		{name: "missing time fields", json: `{"type":"DAILY","tz":"Etc/UTC"}`},
		{name: "monthly by weekday missing occurrence", json: `{"type":"MONTHLY_BY_WEEKDAY","time":"09:30:00","weekday":"TUESDAY"}`},
		{name: "monthly by day type not supported", json: `{"type":"MONTHLY_BY_DAY","time":"02:00:00","tz":"Etc/UTC","monthDay":15}`},
		{name: "yearly missing month", json: `{"type":"YEARLY","time":"00:00:00","monthDay":20}`},
		{name: "cron type not supported", json: `{"type":"CRON","cron":"0 0 17 * * ?","tz":"Etc/UTC"}`},
		{name: "invalid intervalCount", json: `{"type":"DAILY","time":"17:00:00","intervalCount":0}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := commonpamuser.ValidateRotationScheduleJSON(tc.json); err == nil {
				t.Fatalf("expected error for %q", tc.json)
			}
		})
	}
}

func TestRotationScheduleJSONValidator_ValidString(t *testing.T) {
	t.Parallel()

	var resp validator.StringResponse
	req := validator.StringRequest{
		Path:        path.Root("schedule_json"),
		ConfigValue: types.StringValue(`{"type":"DAILY","time":"17:00:00","tz":"Etc/UTC"}`),
	}
	commonpamuser.RotationScheduleJSONValidator().ValidateString(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected valid schedule_json, got: %v", resp.Diagnostics)
	}
}

func TestRotationScheduleJSONValidator_InvalidString(t *testing.T) {
	t.Parallel()

	var resp validator.StringResponse
	req := validator.StringRequest{
		Path:        path.Root("schedule_json"),
		ConfigValue: types.StringValue(`{"type":"WEEKLY","time":"17:00:00"}`),
	}
	commonpamuser.RotationScheduleJSONValidator().ValidateString(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for weekly schedule without weekday")
	}
}
