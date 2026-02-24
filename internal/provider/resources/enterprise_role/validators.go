// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// generatedPasswordComplexityAllowedKeys defines allowed keys and their JSON types for GENERATED_PASSWORD_COMPLEXITY objects.
// Type: "string", "number", "bool", "array_string" (array of strings).
var generatedPasswordComplexityAllowedKeys = map[string]string{
	"domains":               "array_string",
	"length":                "number",
	"lower-use":             "bool",
	"lower-min":             "number",
	"upper-use":             "bool",
	"upper-min":             "number",
	"digit-use":             "bool",
	"digit-min":             "number",
	"special-use":           "bool",
	"special-min":           "number",
	"special":               "string",
	"passphrase-allow":      "bool",
	"passphrase-length":     "number",
	"passphrase-capitalize": "bool",
	"passphrase-number":     "bool",
	"passphrase-separator":  "string",
	"apply-privacy-screen":  "bool",
}

// validateGeneratedPasswordComplexity validates that the JSON string is an array of objects
// with only allowed keys and strict value types. Adds diagnostics to resp on failure.
func validateGeneratedPasswordComplexity(jsonStr string, attrPath path.Path, resp *validator.MapResponse) {
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		resp.Diagnostics.AddAttributeError(
			attrPath,
			"Invalid GENERATED_PASSWORD_COMPLEXITY Value",
			fmt.Sprintf("Value must be a valid JSON array of objects: %v", err),
		)
		return
	}
	for i, obj := range arr {
		for k, v := range obj {
			expectedType, allowed := generatedPasswordComplexityAllowedKeys[k]
			if !allowed {
				resp.Diagnostics.AddAttributeError(
					attrPath,
					"Invalid GENERATED_PASSWORD_COMPLEXITY Key",
					fmt.Sprintf("Object at index %d: key %q is not allowed. Allowed keys: domains, length, lower-use, lower-min, upper-use, upper-min, digit-use, digit-min, special-use, special-min, special, passphrase-allow, passphrase-length, passphrase-capitalize, passphrase-number, passphrase-separator, apply-privacy-screen.", i, k),
				)
				continue
			}
			if !jsonValueMatchesType(v, expectedType) {
				resp.Diagnostics.AddAttributeError(
					attrPath,
					"Invalid GENERATED_PASSWORD_COMPLEXITY Value Type",
					fmt.Sprintf("Object at index %d, key %q: expected type %s, got %T with value %v.", i, k, expectedType, v, v),
				)
			}
		}
	}
}

func jsonValueMatchesType(v interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := v.(string)
		return ok
	case "number":
		switch v.(type) {
		case float64, int, int64:
			return true
		}
		return false
	case "bool":
		_, ok := v.(bool)
		return ok
	case "array_string":
		arr, ok := v.([]interface{})
		if !ok {
			return false
		}
		for _, el := range arr {
			if _, ok := el.(string); !ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// enforcementPoliciesGPCPlanModifier canonicalizes GENERATED_PASSWORD_COMPLEXITY in the plan
// so that semantically equal JSON does not show as a diff (whitespace/key order).
type enforcementPoliciesGPCPlanModifier struct{}

func (enforcementPoliciesGPCPlanModifier) Description(ctx context.Context) string {
	return "Canonicalizes GENERATED_PASSWORD_COMPLEXITY JSON so equivalent values do not produce a diff."
}

func (enforcementPoliciesGPCPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Canonicalizes GENERATED_PASSWORD_COMPLEXITY JSON so equivalent values do not produce a diff."
}

func (enforcementPoliciesGPCPlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	elems := req.PlanValue.Elements()
	modified := make(map[string]attr.Value)
	for k, v := range elems {
		strVal, ok := v.(types.String)
		if !ok {
			modified[k] = v
			continue
		}
		if k == GeneratedPasswordComplexity && !strVal.IsNull() && !strVal.IsUnknown() {
			modified[k] = types.StringValue(utils.CanonicalizeGeneratedPasswordComplexityJSON(strVal.ValueString()))
		} else {
			modified[k] = v
		}
	}
	newMap, diags := types.MapValue(types.StringType, modified)
	resp.Diagnostics.Append(diags...)
	if !diags.HasError() {
		resp.PlanValue = newMap
	}
}

// stringInSlice returns true if s is in list.
func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}

// ----- ENFORCEMENT POLICIES MAP KEY VALIDATOR --------------------------------
// Validates that all map keys (policy keys) in enforcement_policies are valid and values match key-specific rules.
type enforcementPoliciesMapKeyValidator struct{}

func (v enforcementPoliciesMapKeyValidator) Description(ctx context.Context) string {
	return "All enforcement policy keys (map keys) must be valid enforcement policy keys; values must be strings and may have key-specific allowed values."
}

func (v enforcementPoliciesMapKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "All enforcement policy keys (map keys) must be valid enforcement policy keys; values must be strings and may have key-specific allowed values."
}

func (v enforcementPoliciesMapKeyValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for key := range elements {
		pathForKey := req.Path.AtMapKey(key)

		if key == "" {
			resp.Diagnostics.AddAttributeError(
				pathForKey,
				"Invalid Enforcement Policy Key",
				"Enforcement policy key (map key) cannot be empty.",
			)
			continue
		}

		// Validate value type: must be string (MapAttribute value is types.String)
		rawVal := elements[key]
		strVal, ok := rawVal.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				pathForKey,
				"Invalid Enforcement Policy Value Type",
				fmt.Sprintf("Expected string value for enforcement policy '%s', got: %T", key, rawVal),
			)
			continue
		}

		// Skip value validation when unknown (e.g. during plan)
		if strVal.IsUnknown() {
			continue
		}
		// For GENERATED_PASSWORD_COMPLEXITY we do not check "value is string" / null / empty; we only validate JSON array structure (correct keys and value types per object)
		if strVal.IsNull() && key != GeneratedPasswordComplexity {
			resp.Diagnostics.AddAttributeError(
				pathForKey,
				"Invalid Enforcement Policy Value",
				fmt.Sprintf("Enforcement policy value for key '%s' cannot be null.", key),
			)
			continue
		}

		value := strVal.ValueString()

		// Key must be a valid enforcement policy key
		isValidKey := false
		for _, validKey := range ValidEnforcementPolicyKeys {
			if key == validKey {
				isValidKey = true
				break
			}
		}
		if !isValidKey {
			resp.Diagnostics.AddAttributeError(
				pathForKey,
				"Invalid Enforcement Policy Key",
				fmt.Sprintf("Enforcement policy key '%s' is not valid. Must be one of: %s", key, strings.Join(ValidEnforcementPolicyKeys, ", ")),
			)
			continue
		}

		// Key-specific value validation.
		// For GENERATED_PASSWORD_COMPLEXITY we do not check that the value is a simple string (null/empty);
		// we only validate that the value is a JSON array of objects with allowed keys and correct value types per object.
		switch {
		case key == GeneratedPasswordComplexity:
			validateGeneratedPasswordComplexity(value, pathForKey, resp)
		case TwoFactorDurationPolicyKeys[key]:
			if !stringInSlice(value, TwoFactorDurationAllowedValues) {
				resp.Diagnostics.AddAttributeError(
					pathForKey,
					"Invalid Enforcement Policy Value",
					fmt.Sprintf("Enforcement policy '%s' value must be one of: %s", key, strings.Join(TwoFactorDurationAllowedValues, ", ")),
				)
			}
		case KeeperFillPolicyKeys[key]:
			if !stringInSlice(value, KeeperFillAllowedValues) {
				resp.Diagnostics.AddAttributeError(
					pathForKey,
					"Invalid Enforcement Policy Value",
					fmt.Sprintf("Enforcement policy '%s' value must be one of: %s", key, strings.Join(KeeperFillAllowedValues, ", ")),
				)
			}
		default:
			// All other policies: value must be non-empty string
			if value == "" {
				resp.Diagnostics.AddAttributeError(
					pathForKey,
					"Invalid Enforcement Policy Value",
					fmt.Sprintf("Enforcement policy value for key '%s' cannot be an empty string.", key),
				)
			}
		}
	}
}

// ----- PRIVILEGES SET VALIDATOR --------------------------------.
type privilegesValidator struct{}

func (v privilegesValidator) Description(ctx context.Context) string {
	return "Privileges must be valid privilege values."
}

func (v privilegesValidator) MarkdownDescription(ctx context.Context) string {
	return "Privileges must be valid privilege values. Valid values: " + strings.Join(ValidPrivileges, ", ")
}

func (v privilegesValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	for _, elem := range elements {
		// Get the string value
		strValue, ok := elem.(types.String)
		if !ok {
			// This shouldn't happen if ElementType is correct, but handle it anyway
			resp.Diagnostics.AddError(
				"Invalid Privilege Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		// Skip unknown values (e.g. from data source references not yet evaluated at plan time)
		if strValue.IsUnknown() {
			continue
		}

		value := strValue.ValueString()

		// Check for empty strings
		if value == "" {
			resp.Diagnostics.AddError(
				"Empty Privilege String",
				"Privileges set cannot contain empty strings. Each privilege must be a non-empty string.",
			)
			continue
		}

		// Check if the value is in the valid privileges list
		isValid := false
		for _, validPrivilege := range ValidPrivileges {
			if value == validPrivilege {
				isValid = true
				break
			}
		}

		if !isValid {
			resp.Diagnostics.AddError(
				"Invalid Privilege",
				fmt.Sprintf("Privilege '%s' is not valid. Must be one of: %s", value, strings.Join(ValidPrivileges, ", ")),
			)
		}
	}
}
