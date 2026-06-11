// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	rotationPasswordComplexityParts     = 5
	rotationPasswordComplexityMaxPart   = 99
	rotationPasswordComplexityMinLength = 1
)

// RotationScheduleCombinationValidator ensures rotation schedule attributes match
// Commander precedence: on_demand excludes schedule_config, schedule_cron, and
// schedule_json; otherwise at most one of schedule_config, schedule_cron, or
// schedule_json may be set.
type rotationScheduleCombinationValidator struct{}

func RotationScheduleCombinationValidator() rotationScheduleCombinationValidator {
	return rotationScheduleCombinationValidator{}
}

func (rotationScheduleCombinationValidator) Description(_ context.Context) string {
	return "on_demand cannot be combined with schedule_config, schedule_cron, or schedule_json; otherwise at most one of those three may be set."
}

func (rotationScheduleCombinationValidator) MarkdownDescription(_ context.Context) string {
	return "`on_demand` cannot be combined with `schedule_config`, `schedule_cron`, or `schedule_json`. When `on_demand` is not used, set **at most one** of `schedule_config`, `schedule_cron`, or `schedule_json` (same rules as Commander `pam rotation edit`)."
}

func (rotationScheduleCombinationValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	attrs := req.ConfigValue.Attributes()

	onDemand := boolAttrExplicitTrue(attrs, "on_demand")
	scheduleFromConfig := boolAttrExplicitTrue(attrs, "schedule_config")
	cronSet := stringAttrNonEmpty(attrs, "schedule_cron")
	jsonSet := stringAttrNonEmpty(attrs, "schedule_json")

	if onDemand {
		if scheduleFromConfig || cronSet || jsonSet {
			resp.Diagnostics.AddAttributeError(req.Path,
				"Invalid rotation_settings combination",
				"When on_demand is true, omit schedule_config, schedule_cron, and schedule_json. Commander only sends --on-demand in that case; other schedule flags are ignored.")
		}
		return
	}

	n := 0
	if scheduleFromConfig {
		n++
	}
	if cronSet {
		n++
	}
	if jsonSet {
		n++
	}
	if n > 1 {
		resp.Diagnostics.AddAttributeError(req.Path,
			"Invalid rotation_settings schedule",
			"Set at most one of schedule_config (true), schedule_cron, or schedule_json. Multiple values map to conflicting Commander flags (--schedule-config, --schedulecron, --schedulejson).")
	}
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

// RotationPasswordComplexityValidator validates rotation_settings.complexity:
// exactly five comma-separated non-negative integers: length, upper, lower, digits, symbols.
// The password length (first value) must be 1–99 to match Keeper UI limits.
type rotationPasswordComplexityValidator struct{}

func RotationPasswordComplexityValidator() rotationPasswordComplexityValidator {
	return rotationPasswordComplexityValidator{}
}

func (rotationPasswordComplexityValidator) Description(_ context.Context) string {
	return "must be five comma-separated integers length,upper,lower,digits,symbols with length 1–99 and each count 0–99 (Keeper UI limits)"
}

func (rotationPasswordComplexityValidator) MarkdownDescription(_ context.Context) string {
	return "Five comma-separated integers: `length,upper,lower,digits,symbols`. Password **length** (first value) must be **1–99**; each count must be **0–99**, matching Keeper UI limits."
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
