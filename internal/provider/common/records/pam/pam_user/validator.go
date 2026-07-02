// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	rotationPasswordComplexityParts     = 5
	rotationPasswordComplexityMaxPart   = 99
	rotationPasswordComplexityMinLength = 20
)

// RotationProfileRequirementsValidator ensures rotation_profile is set when
// rotation_settings is present and enforces profile-specific required fields.
type rotationProfileRequirementsValidator struct{}

func RotationProfileRequirementsValidator() rotationProfileRequirementsValidator {
	return rotationProfileRequirementsValidator{}
}

func (rotationProfileRequirementsValidator) Description(_ context.Context) string {
	return "rotation_profile is required; profile-specific fields must be set or omitted as required for general, iam_user, scripts_only, and saas"
}

func (rotationProfileRequirementsValidator) MarkdownDescription(_ context.Context) string {
	return "`rotation_profile` is **required** when `rotation_settings` is set. Profile-specific fields: `general` requires `configuration` and `resource` (not `saas_config`); `iam_user` and `scripts_only` require `configuration` (not `resource` or `saas_config`); `saas` requires `configuration` and `saas_config` (not `resource`)."
}

func (rotationProfileRequirementsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	ValidateRotationProfileRequirements(req.Path, req.ConfigValue.Attributes(), resp)
}

// ValidateRotationProfileRequirements checks rotation_profile and dependent fields.
func ValidateRotationProfileRequirements(basePath path.Path, attrs map[string]attr.Value, resp *validator.ObjectResponse) {
	hasConfiguration := stringAttrNonEmpty(attrs, "configuration")
	hasSaaSConfig := stringAttrNonEmpty(attrs, "saas_config")
	hasResource := stringAttrNonEmpty(attrs, "resource")

	profile, ok := stringAttrValue(attrs, "rotation_profile")
	if !ok {
		resp.Diagnostics.AddAttributeError(
			basePath.AtName("rotation_profile"),
			"Missing rotation_profile",
			"rotation_profile is required when rotation_settings is set.",
		)
		return
	}

	switch profile {
	case RotProfileGeneral:
		// Validation for fields that are forbidden for the profile
		if hasSaaSConfig {
			addForbiddenRotationFieldError(basePath, resp, "saas_config", profile)
		}

		// Validation for fields that are required for the profile
		if !hasConfiguration {
			addRequiredRotationFieldError(basePath, resp, "configuration", profile)
		}
		if !hasResource {
			addRequiredRotationFieldError(basePath, resp, "resource", profile)
		}
	case RotProfileIAMUser:
		// Validation for fields that are forbidden for the profile
		if hasResource {
			addForbiddenRotationFieldError(basePath, resp, "resource", profile)
		}
		if hasSaaSConfig {
			addForbiddenRotationFieldError(basePath, resp, "saas_config", profile)
		}

		// Validation for fields that are required for the profile
		if !hasConfiguration {
			addRequiredRotationFieldError(basePath, resp, "configuration", profile)
		}

	case RotProfileScriptsOnly:
		// Validation for fields that are forbidden for the profile
		if hasResource {
			addForbiddenRotationFieldError(basePath, resp, "resource", profile)
		}
		if hasSaaSConfig {
			addForbiddenRotationFieldError(basePath, resp, "saas_config", profile)
		}

		// Validation for fields that are required for the profile
		if !hasConfiguration {
			addRequiredRotationFieldError(basePath, resp, "configuration", profile)
		}

	case RotProfileSaaS:
		// Validation for fields that are forbidden for the profile
		if hasResource {
			addForbiddenRotationFieldError(basePath, resp, "resource", profile)
		}

		// Validation for fields that are required for the profile
		if !hasConfiguration {
			addRequiredRotationFieldError(basePath, resp, "configuration", profile)
		}
		if !hasSaaSConfig {
			addRequiredRotationFieldError(basePath, resp, "saas_config", profile)
		}
	}
}

// addForbiddenRotationFieldError adds an error to the response when a field is set for a profile that is not allowed.
func addForbiddenRotationFieldError(basePath path.Path, resp *validator.ObjectResponse, field, profile string) {
	resp.Diagnostics.AddAttributeError(
		basePath.AtName(field),
		fmt.Sprintf("Invalid %s", field),
		fmt.Sprintf("%s must not be set when rotation_profile is %q.", field, profile),
	)
}

// addRequiredRotationFieldError adds an error to the response when a field is required for a profile but is not set.
func addRequiredRotationFieldError(basePath path.Path, resp *validator.ObjectResponse, field, profile string) {
	resp.Diagnostics.AddAttributeError(
		basePath.AtName(field),
		fmt.Sprintf("Missing %s", field),
		fmt.Sprintf("%s is required when rotation_profile is %q.", field, profile),
	)
}

// RotationScheduleCombinationValidator ensures only one rotation schedule option
// is set within rotation_settings: on_demand, schedule_config, schedule_cron, or
// schedule_json (same mutual exclusivity as Commander pam rotation edit).
type rotationScheduleCombinationValidator struct{}

func RotationScheduleCombinationValidator() rotationScheduleCombinationValidator {
	return rotationScheduleCombinationValidator{}
}

func (rotationScheduleCombinationValidator) Description(_ context.Context) string {
	return "only one of on_demand, schedule_config, schedule_cron, or schedule_json may be set"
}

func (rotationScheduleCombinationValidator) MarkdownDescription(_ context.Context) string {
	return "Set **only one** of `on_demand`, `schedule_config`, `schedule_cron`, or `schedule_json`."
}

func (rotationScheduleCombinationValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if active := activeRotationScheduleOptions(req.ConfigValue.Attributes()); len(active) > 1 {
		resp.Diagnostics.AddAttributeError(req.Path,
			"Invalid rotation_settings schedule",
			fmt.Sprintf(
				"Set only one schedule option in rotation_settings; %s are mutually exclusive.",
				strings.Join(active, ", "),
			),
		)
	}
}

func activeRotationScheduleOptions(attrs map[string]attr.Value) []string {
	var active []string
	if boolAttrExplicitTrue(attrs, "on_demand") {
		active = append(active, "on_demand")
	}
	if boolAttrExplicitTrue(attrs, "schedule_config") {
		active = append(active, "schedule_config")
	}
	if stringAttrNonEmpty(attrs, "schedule_cron") {
		active = append(active, "schedule_cron")
	}
	if stringAttrNonEmpty(attrs, "schedule_json") {
		active = append(active, "schedule_json")
	}
	return active
}

func boolAttrExplicitTrue(attrs map[string]attr.Value, key string) bool {
	v, ok := attrs[key].(types.Bool)
	if !ok || v.IsNull() || v.IsUnknown() {
		return false
	}
	return v.ValueBool()
}

func stringAttrNonEmpty(attrs map[string]attr.Value, key string) bool {
	v, ok := attrs[key].(types.String)
	if !ok || v.IsNull() || v.IsUnknown() {
		return false
	}
	return strings.TrimSpace(v.ValueString()) != ""
}

func stringAttrValue(attrs map[string]attr.Value, key string) (string, bool) {
	v, ok := attrs[key].(types.String)
	if !ok || v.IsNull() || v.IsUnknown() {
		return "", false
	}
	s := strings.TrimSpace(v.ValueString())
	if s == "" {
		return "", false
	}
	return s, true
}

// RotationPasswordComplexityValidator validates rotation_settings.complexity:
// exactly five comma-separated non-negative integers: length, upper, lower, digits, symbols.
// The password length (first value) must be 20–99 to match Keeper UI limits.
type rotationPasswordComplexityValidator struct{}

func RotationPasswordComplexityValidator() rotationPasswordComplexityValidator {
	return rotationPasswordComplexityValidator{}
}

func (rotationPasswordComplexityValidator) Description(_ context.Context) string {
	return "must be five comma-separated integers length,upper,lower,digits,symbols with length 20–99 and each count 0–99 (Keeper UI limits)"
}

func (rotationPasswordComplexityValidator) MarkdownDescription(_ context.Context) string {
	return "Five comma-separated integers: `length,upper,lower,digits,symbols`. Password **length** (first value) must be **20–99**; each count must be **0–99**, matching Keeper UI limits."
}

func (rotationPasswordComplexityValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := ValidateRotationPasswordComplexity(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid rotation complexity", err.Error())
	}
}

// ValidateRotationPasswordComplexity parses complexity for tests and shared checks.
func ValidateRotationPasswordComplexity(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	if len(parts) != rotationPasswordComplexityParts {
		return fmt.Errorf("expected exactly %d comma-separated integers (length,upper,lower,digits,symbols), got %d part(s)", rotationPasswordComplexityParts, len(parts))
	}
	nums := make([]int, rotationPasswordComplexityParts)
	for i, p := range parts {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("part %d %q is not a valid integer", i+1, strings.TrimSpace(parts[i]))
		}
		nums[i] = n
	}
	if nums[0] < rotationPasswordComplexityMinLength || nums[0] > rotationPasswordComplexityMaxPart {
		return fmt.Errorf("password length (first value) must be between %d and %d (Keeper UI maximum), got %d", rotationPasswordComplexityMinLength, rotationPasswordComplexityMaxPart, nums[0])
	}
	for i := 1; i < rotationPasswordComplexityParts; i++ {
		if nums[i] < 0 || nums[i] > rotationPasswordComplexityMaxPart {
			return fmt.Errorf("value for part %d must be between 0 and %d, got %d", i+1, rotationPasswordComplexityMaxPart, nums[i])
		}
	}
	return nil
}
