// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser_test

import (
	"context"
	"testing"

	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
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
		{"schedule_json": types.StringValue(`{"type":"WEEKLY","weekday":"SATURDAY"}`)},
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
			"schedule_json": types.StringValue(`{"type":"DAILY"}`),
		},
		{
			"on_demand":       types.BoolValue(true),
			"schedule_config": types.BoolValue(true),
		},
		{
			"schedule_config": types.BoolValue(true),
			"schedule_json":   types.StringValue(`{"type":"DAILY"}`),
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
